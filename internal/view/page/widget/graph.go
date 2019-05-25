package widget

import (
	"context"
	"sort"
	"time"

	"github.com/slok/grafterm/internal/controller"
	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/log"
	"github.com/slok/grafterm/internal/service/unit"
	"github.com/slok/grafterm/internal/view/render"
	"github.com/slok/grafterm/internal/view/sync"
	"github.com/slok/grafterm/internal/view/template"
)

const (
	graphPointQuantityRetries = 5
)

// graph is a widget that represents values in a two axis graph.
type graph struct {
	controller     controller.Controller
	rendererWidget render.GraphWidget
	widgetCfg      model.Widget
	syncLock       syncingFlag
	logger         log.Logger
}

// NewGraph returns new Graph widget syncer.
func NewGraph(controller controller.Controller, rendererWidget render.GraphWidget, logger log.Logger) sync.Syncer {
	wcfg := rendererWidget.GetWidgetCfg()

	return &graph{
		controller:     controller,
		rendererWidget: rendererWidget,
		widgetCfg:      wcfg,
		logger:         logger,
	}
}

// metricSeries is a helper type that has the metric series and the query
// that has been used to get them.
type metricSeries struct {
	query  model.Query
	series model.MetricSeries
}

func (g *graph) Sync(ctx context.Context, r *sync.Request) error {
	// If already syncing ignore call.
	if g.syncLock.Get() {
		return nil
	}

	// If didn't changed the value means some other sync process
	// already entered before us.
	if !g.syncLock.Set(true) {
		return nil
	}
	defer g.syncLock.Set(false)

	// Get the max capacity of render points (this will be the number of metrics retrieved
	// for the range) of the X axis.
	// If we don't have capacity then return as a dummy sync (no error).
	cap := g.getWindowCapacity()
	if cap <= 0 {
		return nil
	}

	// Gather metrics from multiple queries.
	start := r.TimeRangeStart
	end := r.TimeRangeEnd
	step := end.Sub(start) / time.Duration(cap)
	allSeries := []metricSeries{}
	for _, q := range g.widgetCfg.Graph.Queries {
		//TODO(slok): concurrent queries.
		templatedQ := q
		templatedQ.Expr = r.TemplateData.Render(q.Expr)
		series, err := g.controller.GetRangeMetrics(ctx, templatedQ, start, end, step)
		if err != nil {
			return err
		}

		// Append all received series.
		for _, serie := range series {
			ms := metricSeries{
				query:  q,
				series: serie,
			}
			allSeries = append(allSeries, ms)
		}
	}

	// Merge sort all series.
	metrics := g.sortSeries(allSeries)

	// Transform metric to the ones the render part understands.
	xLabels, indexedTime := g.createIndexedSlices(start, end, step, cap)
	series := g.transformToRenderable(r, metrics, xLabels, indexedTime)

	// Update the render view value.
	g.rendererWidget.Sync(series)
	return nil
}

func (g *graph) sortSeries(allseries []metricSeries) []metricSeries {
	// Sort.
	sort.Slice(allseries, func(i, j int) bool {
		return allseries[i].series.ID < allseries[j].series.ID
	})

	return allseries
}

// createIndexedSlices will create the slices required create a render.Series based on these slices
func (g *graph) createIndexedSlices(start, end time.Time, step time.Duration, capacity int) (xLabels []string, indexedTime []time.Time) {
	xLabels = make([]string, capacity)
	indexedTime = make([]time.Time, capacity)

	// TODO(slok): Calculate the best time format.
	format := unit.TimeRangeTimeStringFormat(end.Sub(start), capacity)
	for i := 0; i < capacity; i++ {
		t := start.Add(time.Duration(i) * step).Local()
		xLabels[i] = t.Format(format)
		indexedTime[i] = t
	}

	return xLabels, indexedTime
}

func (g *graph) transformToRenderable(r *sync.Request, series []metricSeries, xLabels []string, indexedTime []time.Time) []render.Series {
	renderSeries := []render.Series{}

	var colorman widgetColorManager

	// Create the different series to render.
	for _, serie := range series {
		// Create the template data for each series form the sync template
		// data (upper layer template data).
		tplLabels := map[string]interface{}{}
		for k, v := range serie.series.Labels {
			tplLabels[k] = v
		}
		templateData := r.TemplateData.WithData(tplLabels)

		// Get legend and series override based on the legend.
		legend := g.legend(templateData, serie)
		seriesOverride, _ := seriesOverride(g.widgetCfg.Graph.Visualization.SeriesOverride, legend)

		// Init data.
		// This indexes will be used to query the different slices
		// into one single time based XY graph.
		values := make([]*render.Value, len(xLabels))
		timeIndex := 0
		metricIndex := 0
		valueIndex := 0

		// For every value/datapoint we will find where does it belong, to
		// do so we will check one by one each of the metrics if belongs
		// to a current time range, we do this checking if the metric timestamp
		// is after the current timestamp and before the next timestamp.
		for {
			if metricIndex >= len(serie.series.Metrics) ||
				timeIndex >= len(indexedTime) ||
				valueIndex >= len(values) {
				break
			}

			m := serie.series.Metrics[metricIndex]
			ts := indexedTime[timeIndex]

			// If metric is before the timestamp being processed in this
			// iteration then we don't need this metric (too late for it).
			if m.TS.Before(ts) {
				metricIndex++
				continue
			}

			// If we have a next Timestamp then check if the current TS
			// is before the next TS, if not then this metric doesn't
			// belong to this iteration, and belong to a future one.
			if timeIndex < len(indexedTime)-1 {
				nextTS := indexedTime[timeIndex+1]
				// If after means we should ignore this range, so we
				// check the null policy in case we need to fill the
				// empty datapoint space.
				if m.TS.After(nextTS) {
					// The null point mode setting is used to fill the gaps in the values,
					// sometimes the graph has N datapoints and we don't have enough datapoints
					// to create a good renderable graph. This way we can fill this gaps and make
					// the graph renderable.
					switch seriesOverride.NullPointMode {
					case model.NullPointModeAsZero:
						v := render.Value(0)
						values[valueIndex] = &v
					case model.NullPointModeConnected:
						v := render.Value(m.Value)
						values[valueIndex] = &v
					}
					timeIndex++
					valueIndex++
					continue
				}
			}
			// This value belongs here.
			v := render.Value(m.Value)
			values[valueIndex] = &v
			valueIndex++
			metricIndex++
			timeIndex++
		}
		// Create the renderable series.
		serie := render.Series{
			Label:   legend,
			Color:   colorman.GetColorFromSeriesLegend(*g.widgetCfg.Graph, legend),
			XLabels: xLabels,
			Values:  values,
		}

		renderSeries = append(renderSeries, serie)
	}

	return renderSeries
}

func (g *graph) getWindowCapacity() int {
	// Sometimes the widget is not ready to return the capacity of the window, so we try a
	// best effort by trying multiple times with a small sleep so if we are lucky we can get
	// on one of the retries and we don't need to wait for a full sync iteration (e.g 10s),
	// this is not common but happens almost when creating the widgets for the first time.
	cap := 0
	for i := 0; i < graphPointQuantityRetries; i++ {
		cap = g.rendererWidget.GetGraphPointQuantity()
		if cap != 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	return cap
}

// legend will get the correct legend based on the query legend value.
// if this is not set, the legend will be the ID of the metric series,
// if set it will tru rendering the template using the template data.
func (g *graph) legend(templateData template.Data, series metricSeries) string {
	// If no special legend then render with the ID.
	if series.query.Legend == "" {
		return series.series.ID
	}

	// Template the legend.
	return templateData.Render(series.query.Legend)
}

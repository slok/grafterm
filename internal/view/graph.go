package view

import (
	"context"
	"sort"
	"time"

	"github.com/slok/meterm/internal/controller"
	"github.com/slok/meterm/internal/model"
	"github.com/slok/meterm/internal/service/log"
	"github.com/slok/meterm/internal/view/render"
)

const (
	graphPointQuantityRetries = 5
)

var (
	defColors = []string{
		"#7EB26D",
		"#EAB839",
		"#6ED0E0",
		"#EF843C",
		"#E24D42",
	}
)

// graph is a widget that represents values in a two axis graph.
type graph struct {
	controller     controller.Controller
	rendererWidget render.GraphWidget
	appCfg         AppConfig
	widgetCfg      model.Widget
	syncLock       syncingFlag
	logger         log.Logger
}

func newGraph(appCfg AppConfig, controller controller.Controller, rendererWidget render.GraphWidget, logger log.Logger) widget {
	wcfg := rendererWidget.GetWidgetCfg()

	return &graph{
		controller:     controller,
		rendererWidget: rendererWidget,
		appCfg:         appCfg,
		widgetCfg:      wcfg,
		logger:         logger,
	}
}

func (g *graph) sync(ctx context.Context) error {
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
	start := g.appCfg.TimeRangeStart
	end := g.appCfg.TimeRangeEnd
	step := end.Sub(start) / time.Duration(cap)
	allSeries := [][]model.MetricSeries{}
	for _, q := range g.widgetCfg.Graph.Queries {
		//TODO(slok): concurrent queries.
		series, err := g.controller.GetRangeMetrics(ctx, q, start, end, step)
		if err != nil {
			return err
		}
		allSeries = append(allSeries, series)
	}

	// Merge sort all series.
	metrics := g.mergeAndSortSeries(allSeries...)

	// Transform metric to the ones the render part understands.
	xLabels, indexedTime := g.createIndexedSlices(start, step, cap)
	series := g.transformMetrics(metrics, xLabels, indexedTime)

	// Update the render view value.
	g.rendererWidget.Sync(series)
	return nil
}

func (g *graph) mergeAndSortSeries(allseries ...[]model.MetricSeries) []model.MetricSeries {
	res := []model.MetricSeries{}

	// Merge.
	for _, series := range allseries {
		for _, serie := range series {
			res = append(res, serie)
		}
	}

	// Sort.
	sort.Slice(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})

	return res
}

// createIndexedSlices will create the slices required create a render.Series based on these slices
func (g *graph) createIndexedSlices(start time.Time, step time.Duration, capacity int) (xLabels []string, indexedTime []time.Time) {
	xLabels = make([]string, capacity)
	indexedTime = make([]time.Time, capacity)

	// TODO(slok): Calculate the best time format.
	format := "15:04:05"
	for i := 0; i < capacity; i++ {
		t := start.Add(time.Duration(i) * step)
		xLabels[i] = t.Format(format)
		indexedTime[i] = t
	}

	return xLabels, indexedTime
}

func (g *graph) transformMetrics(series []model.MetricSeries, xLabels []string, indexedTime []time.Time) []render.Series {
	renderSeries := []render.Series{}

	// Create the different series to render.
	for i, serie := range series {
		// Init data.
		label := serie.ID
		color := defColors[i%len(defColors)]
		values := make([]*render.Value, len(xLabels))
		timeIndex := 0
		metricIndex := 0
		valueIndex := 0

		// For every value/datapoint.
		for {
			if metricIndex >= len(serie.Metrics) ||
				timeIndex >= len(indexedTime) ||
				valueIndex >= len(values) {
				break
			}

			m := serie.Metrics[metricIndex]
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
				if m.TS.After(nextTS) {
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

		serie := render.Series{
			Label:   label,
			Color:   color,
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

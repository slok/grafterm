package page

import (
	"context"

	"github.com/slok/grafterm/internal/view/sync"
	"github.com/slok/grafterm/internal/view/template"
)

// withWidgetDataMiddleware controls the variables data
// the widget receives, it wraps any widget and will
// mutate the variable data (updating, adding, deleting...)
// the widget receives on every sync.
//
// It has the static data the widget will receive on all
// the syncs, this way the widget doesn't need to store
// the static data.
//
// It also controls the data that the user wants to override
// (for example via cmd flags).
//
// Priority chain.
// 1- OverrideData
// 2- SyncData
// 3- StaticData
func withWidgetDataMiddleware(data template.Data, overrideData template.Data, next sync.Syncer) sync.Syncer {
	return &widgetDataMiddleware{
		staticData:   data,
		overrideData: overrideData,
		next:         next,
	}
}

type widgetDataMiddleware struct {
	staticData   template.Data
	overrideData template.Data
	next         sync.Syncer
}

func (w widgetDataMiddleware) Sync(ctx context.Context, r *sync.Request) error {
	// Add the sync data to the static data and place it again on the cfg.
	data := w.staticData.WithData(r.TemplateData)
	// Override the data asked by the user.
	data = data.WithData(w.overrideData)

	r.TemplateData = data
	return w.next.Sync(ctx, r)
}

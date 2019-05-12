package view

import (
	"context"

	"github.com/slok/grafterm/internal/view/template"
)

// withWidgetDataMiddleware controls the variables data
// the widget receives, it wraps any widget and will
// inject this data variables the widget will receive so the
// widget doesn't need to store this data.
//
// It has the static data the widget will receive on all
// the syncs, this way the widget doesn't need to store
// the static data.
//
// It also controls the data that the user asked to override
// (for example via cmd flags) on every data variables received
// on each sync.
//
// Priority chain.
// 1- overrideData
// 2- SyncData
// 3- StaticData

func withWidgetDataMiddleware(data template.Data, overrideData template.Data, next widget) widget {
	return &widgetDataMiddleware{
		staticData:   data,
		overrideData: overrideData,
		next:         next,
	}
}

type widgetDataMiddleware struct {
	staticData   template.Data
	overrideData template.Data
	next         widget
}

func (w widgetDataMiddleware) sync(ctx context.Context, cfg syncConfig) error {
	// Add the sync data to the static data and place it again on the cfg.
	data := w.staticData.WithData(cfg.templateData)
	// Override the data asked by the user.
	data = data.WithData(w.overrideData)

	cfg.templateData = data
	return w.next.sync(ctx, cfg)
}

package view

import (
	"context"

	"github.com/slok/meterm/internal/view/template"
)

// withWidgetDataMiddleware will wrap a widget and every time is synced it will
// add the static data to the sync data so it has it's corresponding
// data so the widget itself doesn't need to store static data like the
// dashboard, widget static data or similar.
func withWidgetDataMiddleware(data template.Data, next widget) widget {
	return &widgetDataMiddleware{
		staticData: data,
		next:       next,
	}
}

type widgetDataMiddleware struct {
	staticData template.Data
	next       widget
}

func (w widgetDataMiddleware) sync(ctx context.Context, cfg syncConfig) error {
	// Add the sync data to the static data and place it again on the cfg.
	data := w.staticData.WithData(cfg.templateData)
	cfg.templateData = data

	return w.next.sync(ctx, cfg)
}

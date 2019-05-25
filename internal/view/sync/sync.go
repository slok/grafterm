package sync

import (
	"context"
	"time"

	"github.com/slok/grafterm/internal/view/template"
)

// Request is a sync iteration request, similar approach to an HTTP request.
type Request struct {
	TimeRangeStart time.Time
	TimeRangeEnd   time.Time
	TemplateData   template.Data
}

// Syncer is a component that will be synced with the app iteration sync,
// depending on the page it could end rendered in the screen.
type Syncer interface {
	Sync(ctx context.Context, r *Request) error
}

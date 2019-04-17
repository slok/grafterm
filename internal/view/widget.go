package view

import (
	"context"
	"time"

	"github.com/slok/meterm/internal/view/template"
)

type syncConfig struct {
	timeRangeStart time.Time
	timeRangeEnd   time.Time
	templateData   template.Data
}

type widget interface {
	sync(ctx context.Context, cfg syncConfig) error
}

package view

import (
	"context"
	"time"
)

type syncConfig struct {
	timeRangeStart time.Time
	timeRangeEnd   time.Time
}

type widget interface {
	sync(ctx context.Context, cfg syncConfig) error
}

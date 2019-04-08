package view

import (
	"context"
)

type widget interface {
	sync(ctx context.Context) error
}

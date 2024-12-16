package repository

import (
	"context"
)

type PdrisRepository interface {
	UpdateValue(ctx context.Context, value int) error
	GetValue(ctx context.Context) (int, error)
}

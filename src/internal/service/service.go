package service

import (
	"context"
)

type PdrisService interface {
	UpdateValue(ctx context.Context, value int) error
	GetValue(ctx context.Context) (int, error)
}

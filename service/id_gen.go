package service

import (
	"context"

	"github.com/busyfree/leaf-go/models"
)

type IDGen interface {
	Get(ctx context.Context, key string) models.Result
	Init(ctx context.Context) bool
}

package service

import (
	"context"

	"github.com/busyfree/leaf-go/models"
)

type ZeroIDGenImpl struct {
}

func NewZeroIDGenImpl() *ZeroIDGenImpl {
	s := new(ZeroIDGenImpl)
	return s
}

func (s *ZeroIDGenImpl) Init(ctx context.Context) bool {
	return true
}

func (s *ZeroIDGenImpl) Get(ctx context.Context, key string) models.Result {
	return models.Result{Status: models.SUCCESS}
}

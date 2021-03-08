package service

import (
	"context"

	"github.com/busyfree/leaf-go/dao"
)

type LeafAllocService struct {
	Dao *dao.LeafAllocDao
}

func NewLeafAllocService() *LeafAllocService {
	s := new(LeafAllocService)
	s.Dao = dao.NewLeafAllocDao()
	return s
}

func (s *LeafAllocService) GetAllTags(ctx context.Context) (array []string, err error) {
	return s.Dao.GetAllTags(ctx)
}

func (s *LeafAllocService) GetAllLeafAllocs(ctx context.Context) (array []*dao.LeafAllocDao, err error) {
	return s.Dao.GetAllLeafAllocs(ctx)
}

func (s *LeafAllocService) GetLeafAlloc(ctx context.Context, tag string) (err error) {
	return s.Dao.GetLeafAlloc(ctx, tag)
}

func (s *LeafAllocService) UpdateMaxId(ctx context.Context, tag string) (err error) {
	return s.Dao.UpdateMaxId(ctx, tag)
}

func (s *LeafAllocService) UpdateMaxIdByCustomStep(ctx context.Context, step int, tag string) (err error) {
	return s.Dao.UpdateMaxIdByCustomStep(ctx, step, tag)
}

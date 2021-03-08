package dao

import (
	"github.com/busyfree/leaf-go/models"
	"go.uber.org/atomic"
)

type SegmentDao struct {
	models.Segment
	Buffer *SegmentBufferDao
}

func NewSegmentDao(dao *SegmentBufferDao) *SegmentDao {
	s := new(SegmentDao)
	s.Buffer = dao
	s.Value = atomic.NewInt64(0)
	return s
}
func (dao *SegmentDao) GetValue() *atomic.Int64 {
	return dao.Value
}

func (dao *SegmentDao) SetValue(value *atomic.Int64) {
	dao.Value = value
}

func (dao *SegmentDao) GetMax() int64 {
	return dao.Max
}

func (dao *SegmentDao) SetMax(max int64) {
	dao.Max = max
}

func (dao *SegmentDao) GetStep() int {
	return dao.Step
}

func (dao *SegmentDao) SetStep(step int) {
	dao.Step = step
}

func (dao *SegmentDao) GetBuffer() *SegmentBufferDao {
	return dao.Buffer
}

func (dao *SegmentDao) GetIdle() int64 {
	value := dao.GetValue().Load()
	return dao.GetMax() - value
}

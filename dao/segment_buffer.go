package dao

import (
	"sync"

	"go.uber.org/atomic"

	"github.com/busyfree/leaf-go/models"
)

type SegmentBufferDao struct {
	models.SegmentBuffer
	Segments []*SegmentDao
}

func NewSegmentBufferDao() *SegmentBufferDao {
	s := new(SegmentBufferDao)
	s.Segments = make([]*SegmentDao, 0, 0)
	segment1 := NewSegmentDao(s)
	segment2 := NewSegmentDao(s)
	s.Segments = append(s.Segments, segment1, segment2)
	s.CurrentPos = 0
	s.NextReady = false
	s.InitOk = false
	s.ThreadRunning = atomic.NewBool(false)
	s.RWMutex = &sync.RWMutex{}
	return s
}

func (dao *SegmentBufferDao) GetKey() string {
	return dao.Key
}

func (dao *SegmentBufferDao) SetKey(key string) {
	dao.Key = key
}

func (dao *SegmentBufferDao) GetSegments() []*SegmentDao {
	return dao.Segments
}

func (dao *SegmentBufferDao) GetCurrent() *SegmentDao {
	return dao.Segments[dao.CurrentPos]
}

func (dao *SegmentBufferDao) GetCurrentPos() int {
	return dao.CurrentPos
}

func (dao *SegmentBufferDao) NextPos() int {
	return (dao.CurrentPos + 1) % 2
}

func (dao *SegmentBufferDao) SwitchPos() {
	dao.CurrentPos = dao.NextPos()
}

func (dao *SegmentBufferDao) IsInitOk() bool {
	return dao.InitOk
}

func (dao *SegmentBufferDao) SetInitOk(initOk bool) {
	dao.InitOk = initOk
}

func (dao *SegmentBufferDao) IsNextReady() bool {
	return dao.NextReady
}

func (dao *SegmentBufferDao) SetNextReady(nextReady bool) {
	dao.NextReady = nextReady
}

func (dao *SegmentBufferDao) GetThreadRunning() *atomic.Bool {
	return dao.ThreadRunning
}

func (dao *SegmentBufferDao) ReadLock() {
	dao.RWMutex.RLock()
}

func (dao *SegmentBufferDao) ReadUnLock() {
	dao.RWMutex.RUnlock()
}

func (dao *SegmentBufferDao) WriteLock() {
	dao.RWMutex.Lock()
}

func (dao *SegmentBufferDao) WriteULock() {
	dao.RWMutex.Unlock()
}

func (dao *SegmentBufferDao) GetStep() int {
	return dao.Step
}

func (dao *SegmentBufferDao) SetStep(step int) {
	dao.Step = step
}

func (dao *SegmentBufferDao) GetMinStep() int {
	return dao.MinStep
}

func (dao *SegmentBufferDao) SetMinStep(minStep int) {
	dao.MinStep = minStep
}

func (dao *SegmentBufferDao) GetUpdateTimeStamp() int64 {
	return dao.UpdatedAt
}

func (dao *SegmentBufferDao) SetUpdateTimeStamp(ts int64) {
	dao.UpdatedAt = ts
}

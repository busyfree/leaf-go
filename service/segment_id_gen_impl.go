package service

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/busyfree/leaf-go/dao"
	"github.com/busyfree/leaf-go/models"
	"github.com/busyfree/leaf-go/util/log"
	"github.com/busyfree/leaf-go/util/timeutil"
	"go.uber.org/atomic"
)

type SegmentIDGenImpl struct {
	maxStep         int
	segmentDuration int64
	initOk          bool
	cache           *sync.Map
	leafAllocDao    *dao.LeafAllocDao
}

func NewSegmentIDGenImpl() *SegmentIDGenImpl {
	s := new(SegmentIDGenImpl)
	s.maxStep = 1000000
	s.segmentDuration = 15 * 60 * 1000
	s.initOk = false
	s.leafAllocDao = dao.NewLeafAllocDao()
	s.cache = new(sync.Map)
	s.Init()
	return s
}

func (s *SegmentIDGenImpl) Init() bool {
	s.updateCacheFromDb(nil)
	s.initOk = true
	s.updateCacheFromDbAtEveryMinute()
	return s.initOk
}

func (s *SegmentIDGenImpl) updateCacheFromDbAtEveryMinute() {
	ticker := time.NewTicker(time.Duration(60) * time.Second)
	runtime.LockOSThread()
	go func(t *time.Ticker) {
		for {
			select {
			case <-t.C:
				s.updateCacheFromDb(context.Background())
			}
		}
	}(ticker)
}

func (s *SegmentIDGenImpl) Get(ctx context.Context, key string) models.Result {
	if !s.initOk {
		return models.NewResult(int64(models.EXCEPTION_ID_IDCACHE_INIT_FALSE), models.EXCEPTION)
	}
	cacheTags := make(map[string]*dao.SegmentBufferDao, 0)
	s.cache.Range(func(k, v interface{}) bool {
		cacheTags[k.(string)] = v.(*dao.SegmentBufferDao)
		return true
	})
	cacheSegmentBuffer, ok := cacheTags[key]
	if !ok {
		return models.NewResult(int64(models.EXCEPTION_ID_KEY_NOT_EXISTS), models.EXCEPTION)
	}
	if !cacheSegmentBuffer.InitOk {
		err := s.updateSegmentFromDb(ctx, key, cacheSegmentBuffer.GetCurrent())
		if err != nil {
			cacheSegmentBuffer.SetInitOk(false)
			return models.NewResult(int64(models.EXCEPTION_ID_IDCACHE_INIT_FALSE), models.EXCEPTION)
		}
		cacheSegmentBuffer.SetInitOk(true)
		s.cache.Store(key, cacheSegmentBuffer)
	}
	return s.getIdFromSegmentBuffer(cacheSegmentBuffer)
}

func (s *SegmentIDGenImpl) loadNextSegmentFromDb(cacheSegmentBufferDao *dao.SegmentBufferDao) {
	nextSegment := cacheSegmentBufferDao.GetSegments()[cacheSegmentBufferDao.NextPos()]
	err := s.updateSegmentFromDb(context.Background(), cacheSegmentBufferDao.GetKey(), nextSegment)
	if err != nil {
		cacheSegmentBufferDao.GetThreadRunning().Store(false)
	}
	cacheSegmentBufferDao.WriteLock()
	cacheSegmentBufferDao.SetNextReady(true)
	cacheSegmentBufferDao.GetThreadRunning().Store(false)
	cacheSegmentBufferDao.WriteULock()
	return
}

func (s *SegmentIDGenImpl) getIdFromSegmentBuffer(cacheSegmentBufferDao *dao.SegmentBufferDao) models.Result {
	for {
		cacheSegmentBufferDao.RWMutex.RLock()
		segmentDao := cacheSegmentBufferDao.GetCurrent()
		if !cacheSegmentBufferDao.IsNextReady() && (segmentDao.GetIdle() < int64(0.9*float64(segmentDao.GetStep()))) && cacheSegmentBufferDao.GetThreadRunning().CAS(false, true) {
			go s.loadNextSegmentFromDb(cacheSegmentBufferDao)
		}
		cacheSegmentBufferDao.RWMutex.RUnlock()
		value := segmentDao.GetValue().Load()
		_ = segmentDao.GetValue().Inc()
		if value < segmentDao.GetMax() {
			return models.NewResult(value, models.SUCCESS)
		}
		s.waitAndSleep(cacheSegmentBufferDao)
		cacheSegmentBufferDao.WriteLock()
		if cacheSegmentBufferDao.IsNextReady() {
			cacheSegmentBufferDao.SwitchPos()
			cacheSegmentBufferDao.SetNextReady(false)
		} else {
			return models.NewResult(int64(models.EXCEPTION_ID_TWO_SEGMENTS_ARE_NULL), models.EXCEPTION)
		}
		cacheSegmentBufferDao.WriteULock()
		s.cache.Store(cacheSegmentBufferDao.GetKey(), cacheSegmentBufferDao)
	}
}

func (s *SegmentIDGenImpl) waitAndSleep(segmentBufferDao *dao.SegmentBufferDao) {
	roll := 0
	for segmentBufferDao.GetThreadRunning().Load() {
		roll++
		if roll > 10000 {
			time.Sleep(time.Duration(10) * time.Millisecond)
			break
		}
	}
}

func (s *SegmentIDGenImpl) updateSegmentFromDb(ctx context.Context, key string, segment *dao.SegmentDao) (err error) {
	logger := log.Get(ctx)
	segmentBufferDao := segment.GetBuffer()
	var leafAllocDao = dao.NewLeafAllocDao()
	if !segmentBufferDao.InitOk {
		err = leafAllocDao.UpdateMaxIdAndGetLeafAlloc(ctx, key)
		if err != nil {
			logger.Infof("leafAllocDao.UpdateMaxIdAndGetLeafAllocErr:%v", err)
			return
		}
		segmentBufferDao.SetStep(leafAllocDao.Step)
		segmentBufferDao.SetMinStep(leafAllocDao.Step)
	} else if segmentBufferDao.GetUpdateTimeStamp() == 0 {
		err = leafAllocDao.UpdateMaxIdAndGetLeafAlloc(ctx, key)
		if err != nil {
			logger.Infof("leafAllocDao.UpdateMaxIdAndGetLeafAllocErr:%v", err)
			return
		}
		segmentBufferDao.SetUpdateTimeStamp(timeutil.MsTimestampNow())
		segmentBufferDao.SetStep(leafAllocDao.Step)
		segmentBufferDao.SetMinStep(leafAllocDao.Step)
	} else {
		duration := timeutil.MsTimestampNow() - segmentBufferDao.GetUpdateTimeStamp()
		nextStep := segmentBufferDao.GetStep()
		if duration < s.segmentDuration {
			if nextStep*2 > s.maxStep {
			} else {
				nextStep = nextStep * 2
			}
		} else if duration < s.segmentDuration*2 {

		} else {
			if nextStep/2 >= segmentBufferDao.GetMinStep() {
				nextStep = nextStep / 2
			}
		}
		leafAllocNewDao := dao.NewLeafAllocDao()
		leafAllocNewDao.BizTag = key
		leafAllocNewDao.Step = nextStep
		err = leafAllocDao.UpdateMaxIdByCustomStepAndGetLeafAlloc(ctx, leafAllocNewDao)
		if err != nil {
			logger.Infof("leafAllocDao.UpdateMaxIdByCustomStepAndGetLeafAllocErr:%v", err)
			return
		}
		segmentBufferDao.SetUpdateTimeStamp(timeutil.MsTimestampNow())
		segmentBufferDao.SetStep(nextStep)
		segmentBufferDao.SetMinStep(leafAllocDao.Step)
	}

	value := leafAllocDao.MaxId - int64(segmentBufferDao.GetStep())
	segment.GetValue().Store(value)
	segment.SetMax(leafAllocDao.MaxId)
	segment.SetStep(segmentBufferDao.GetStep())
	s.cache.Store(key, segmentBufferDao)
	return
}

func (s *SegmentIDGenImpl) updateCacheFromDb(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	var leafAllocDao = dao.NewLeafAllocDao()
	dbTags, err := leafAllocDao.GetAllTags(ctx)
	if err != nil {
		return
	}
	if len(dbTags) == 0 {
		return
	}
	cacheTags := make(map[string]string, 0)
	insertTags := make([]string, 0, 0)
	removeTags := make([]string, 0, 0)
	s.cache.Range(func(k, v interface{}) bool {
		cacheTags[k.(string)] = k.(string)
		return true
	})
	if len(dbTags) > 0 {
		for _, k := range dbTags {
			if _, ok := cacheTags[k]; !ok {
				insertTags = append(insertTags, k)
			}
		}
	}
	if len(insertTags) > 0 {
		for _, k := range insertTags {
			segmentBuffer := dao.NewSegmentBufferDao()
			segmentBuffer.SetKey(k)
			segment := segmentBuffer.GetCurrent()
			segment.SetValue(atomic.NewInt64(0))
			segment.SetMax(0)
			segment.SetStep(0)
			s.cache.Store(k, segmentBuffer)
		}
	}
	if len(dbTags) > 0 {
		for _, k := range dbTags {
			if _, ok := cacheTags[k]; !ok {
				removeTags = append(removeTags, k)
			}
		}
		if len(removeTags) > 0 && len(cacheTags) > 0 {
			for _, tag := range removeTags {
				s.cache.Delete(tag)
			}
		}
	}

}

func (s *SegmentIDGenImpl) GetAllLeafAllocs(ctx context.Context) (array []*dao.LeafAllocDao, err error) {
	return s.leafAllocDao.GetAllLeafAllocs(ctx)
}

func (s *SegmentIDGenImpl) GetCache(ctx context.Context) *sync.Map {
	return s.cache
}

func (s *SegmentIDGenImpl) GetDao(ctx context.Context) *dao.LeafAllocDao {
	return s.leafAllocDao
}

func (s *SegmentIDGenImpl) SetDao(dao *dao.LeafAllocDao) {
	s.leafAllocDao = dao
}

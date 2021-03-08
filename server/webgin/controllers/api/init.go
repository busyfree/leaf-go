package api

import (
	"github.com/busyfree/leaf-go/service"
)

var (
	segmentService   *service.SegmentIDGenImpl
	snowflakeService *service.SnowFlakeIdGenImpl
)

func Init(s *service.SegmentIDGenImpl, snowflake *service.SnowFlakeIdGenImpl) {
	segmentService = s
	snowflakeService = snowflake
}

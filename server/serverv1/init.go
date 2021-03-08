package serverv1

import "github.com/busyfree/leaf-go/service"

var (
	segmentService   *service.SegmentIDGenImpl
	snowflakeService *service.SnowFlakeIdGenImpl
)

func Init(segment *service.SegmentIDGenImpl, snowflake *service.SnowFlakeIdGenImpl) {
	segmentService = segment
	snowflakeService = snowflake
}

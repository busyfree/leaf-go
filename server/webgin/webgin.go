

package webgin

import (
	"context"
	"github.com/busyfree/leaf-go/service"

	"github.com/busyfree/leaf-go/util/log"
	"github.com/gin-gonic/gin"
)

const (
	BASEURL = "/web/"
)

var (
	GinRoute *gin.Engine
	logger   = log.Get(context.Background())
)

func InitWebGin(segment *service.SegmentIDGenImpl, snowflake *service.SnowFlakeIdGenImpl) {
	GinRoute = gin.New()
	initRoute(segment, snowflake)
}

package middlewares

import (
	"fmt"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/busyfree/leaf-go/util/ctxkit"
	"github.com/busyfree/leaf-go/util/log"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := log.Get(ctx.Request.Context())
		path := ctx.Request.URL.Path
		start := time.Now()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := ctx.Writer.Status()
		clientUserAgent := ctx.Request.UserAgent()
		referer := ctx.Request.Referer()
		dataLength := ctx.Writer.Size()
		uid, _ := ctx.Get(cast.ToString(ctxkit.UserIDKey))
		if dataLength < 0 {
			dataLength = 0
		}
		if len(ctx.Errors) > 0 {
			logger.Error(ctx.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("%v \"%s %s\" %d %d \"%s\" \"%s\" (%dms)", uid, ctx.Request.Method, path, statusCode, dataLength, referer, clientUserAgent, latency)
			if statusCode > 499 {
				logger.Error(msg)
			} else if statusCode > 399 {
				logger.Warn(msg)
			} else {
				logger.Info(msg)
			}
		}
		ctx.Next()
	}
}

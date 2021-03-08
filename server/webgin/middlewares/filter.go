package middlewares

import (
	"github.com/busyfree/leaf-go/util/ctxkit"

	"github.com/gin-gonic/gin"
)

func FilterMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request.Clone(ctx)
		rawCtx := req.Context()
		rawCtx = ctxkit.WithUserIP(rawCtx, ctx.ClientIP())
		req = req.WithContext(rawCtx)
		ctx.Request = req
		ctx.Next()
		return
	}
}

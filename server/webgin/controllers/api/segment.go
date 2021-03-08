package api

import (
	"github.com/gin-gonic/gin"
)

type SegmentController struct{}

func (c *SegmentController) Get(ctx *gin.Context) {
	key := ctx.Param("key")
	r := segmentService.Get(ctx, key)
	ctx.JSON(200, r)
	return
}

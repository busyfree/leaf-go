package api

import (
	"github.com/gin-gonic/gin"
)

type SnowFlakeController struct{}

func (c *SnowFlakeController) Get(ctx *gin.Context) {
	key := ctx.Param("key")
	r := snowflakeService.Get(ctx, key)
	ctx.JSON(200, r)
	return
}

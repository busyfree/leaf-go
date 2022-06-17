package acp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/busyfree/leaf-go/dao"
	"github.com/busyfree/leaf-go/models"
)

type MonitorController struct{}

func (c *MonitorController) Cache(ctx *gin.Context) {
	cacheMaps := segmentService.GetCache(ctx)
	cacheTags := make(map[string]*dao.SegmentBufferDao, 0)
	cacheMaps.Range(func(k, v interface{}) bool {
		cacheTags[k.(string)] = v.(*dao.SegmentBufferDao)
		return true
	})
	data := make([]*models.SegmentBufferView, 0, 0)
	if len(cacheTags) > 0 {
		for _, dao := range cacheTags {
			v := &models.SegmentBufferView{}
			v.InitOk = dao.IsInitOk()
			v.Key = dao.GetKey()
			v.Pos = dao.GetCurrentPos()
			v.NextReady = dao.IsNextReady()
			segments := dao.GetSegments()
			v.Max0 = segments[0].GetMax()
			v.Value0 = segments[0].GetValue().Load()
			v.Step0 = segments[0].GetStep()
			v.Max1 = segments[1].GetMax()
			v.Value1 = segments[1].GetValue().Load()
			v.Step1 = segments[1].GetStep()
			data = append(data, v)
		}
	}
	ctx.HTML(200, "cache.html", gin.H{"data": data})
	return
}

func (c *MonitorController) DB(ctx *gin.Context) {
	daos, err := segmentService.GetAllLeafAllocs(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.HTML(200, "db.html", gin.H{"daos": daos})
	return
}

func (c *MonitorController) Decode(ctx *gin.Context) {
	out := snowflakeService.DecodeSnowflakeId(ctx.Param("key"))
	ctx.JSON(200, out)
	return
}

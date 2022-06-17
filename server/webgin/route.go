package webgin

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/busyfree/leaf-go/server/webgin/controllers/acp"
	"github.com/busyfree/leaf-go/server/webgin/controllers/api"
	"github.com/busyfree/leaf-go/server/webgin/middlewares"
	"github.com/busyfree/leaf-go/service"
	"github.com/busyfree/leaf-go/util/conf"
)

func formatBool(t bool) string {
	if t {
		return "true"
	}
	return "false"
}

func initRoute(segmentService *service.SegmentIDGenImpl, snowflakeService *service.SnowFlakeIdGenImpl) {
	GinRoute.SetFuncMap(template.FuncMap{
		"formatBool": formatBool,
	})

	sessionMaxIdle := conf.GetInt("SESSION_REDIS_MAX_CONS")
	sessionIPPort := conf.GetString("SESSION_REDIS_IP_PORT")
	sessionDB := conf.GetString("SESSION_REDIS_DB")
	if sessionMaxIdle <= 0 {
		sessionMaxIdle = 10
	}
	if len(sessionIPPort) == 0 {
		sessionIPPort = "localhost:6379"
	}
	if len(sessionDB) == 0 {
		sessionDB = "0"
	}

	GinRoute.LoadHTMLGlob("views/***/**/*")
	fs := filepath.Join(conf.GetConfigPath(), "public", "static")
	GinRoute.StaticFS("/web/static", http.Dir(fs))
	GinRoute.Use(gin.Recovery())
	middlewares.FilterMiddleware()
	GinRoute.Use(middlewares.Logger())

	GinRoute.MaxMultipartMemory = 100 << 20 // 100M

	base := new(acp.BaseController)
	GinRoute.GET(BASEURL+"ping", base.Ping)

	v1 := GinRoute.Group(BASEURL + "v1")
	v1Dashboard := v1.Group("/acp")

	monitorAPIGroup := v1Dashboard.Group("/monitor")
	{
		monitor := new(acp.MonitorController)
		monitorAPIGroup.GET("/cache", monitor.Cache)
		monitorAPIGroup.GET("/db", monitor.DB)
		monitorAPIGroup.GET("/decode/:key", monitor.Decode)
	}

	v1Front := v1.Group("/api")

	segmentAPIGroup := v1Front.Group("/segment")
	{
		segment := new(api.SegmentController)
		segmentAPIGroup.GET("/get/:key", segment.Get)
		segmentAPIGroup.POST("/get/:key", segment.Get)
	}
	snowflakeAPIGroup := v1Front.Group("/snowflake")
	{
		snowflake := new(api.SnowFlakeController)
		snowflakeAPIGroup.GET("/get/:key", snowflake.Get)
		snowflakeAPIGroup.POST("/get/:key", snowflake.Get)
	}
	acp.Init(segmentService, snowflakeService)
	api.Init(segmentService, snowflakeService)
}

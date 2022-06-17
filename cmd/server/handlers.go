package server

import (
	"net/http"

	"github.com/withgame/twirp"

	"github.com/busyfree/leaf-go/cmd/server/hook"
	"github.com/busyfree/leaf-go/rpc/v1/public"
	"github.com/busyfree/leaf-go/server/serverv1"
	"github.com/busyfree/leaf-go/server/webgin"
	"github.com/busyfree/leaf-go/service"
	"github.com/busyfree/leaf-go/util/conf"
)

var hooks = twirp.ChainHooks(
	hook.NewRequestID(),
	hook.NewLog(),
)

var privateHooks = twirp.ChainHooks(
	hook.NewRequestID(),
	hook.NewLog(),
)

var allowGetHooks = twirp.ChainHooks(
	hook.NeAllowGet(),
	hook.NewRequestID(),
	hook.NewLog(),
)

func initMux(mux *http.ServeMux, isInternal bool) {
	snowflakePort := conf.GetInt("LEAF_SNOWFLAKE_PORT")
	leafSnowflakeTime := conf.GetTime("LEAF_SNOWFLAKE_START_TIME")
	leafSnowflakeTwepoch := leafSnowflakeTime.Unix() * 1000
	if leafSnowflakeTime.IsZero() {
		leafSnowflakeTwepoch = 1288834974657
	}
	snowflakeService := service.NewSnowFlakeIdGenImpl(snowflakePort, leafSnowflakeTwepoch)
	segmentService := service.NewSegmentIDGenImpl()
	{
		serverv1.Init(segmentService, snowflakeService)
		serverPublic := &serverv1.Public{}
		handler := public.NewServerServer(serverPublic, hooks)
		mux.Handle(public.ServerPathPrefix, handler)
	}
	{
		webgin.InitWebGin(segmentService, snowflakeService)
		mux.Handle(webgin.BASEURL, webgin.GinRoute)
	}
}

func initInternalMux(mux *http.ServeMux) {
}

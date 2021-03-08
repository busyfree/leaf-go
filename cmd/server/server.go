package server

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"sync"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/busyfree/leaf-go/service"
	"github.com/busyfree/leaf-go/util"
	"github.com/busyfree/leaf-go/util/conf"
	"github.com/busyfree/leaf-go/util/ctxkit"
	"github.com/busyfree/leaf-go/util/log"
	_ "github.com/busyfree/leaf-go/util/redis"
	"github.com/busyfree/leaf-go/util/trace"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cast"
)

var (
	server *http.Server
	logger = log.Get(context.Background())
)

type panicHandler struct {
	handler http.Handler
}

// 从 http 标准库搬来的
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	_ = tc.SetKeepAlive(true)
	_ = tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func startSpan(r *http.Request) (*http.Request, opentracing.Span) {
	operation := "ServerHTTP"

	ctx := r.Context()
	var span opentracing.Span

	tracer := opentracing.GlobalTracer()
	carrier := opentracing.HTTPHeadersCarrier(r.Header)

	if spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, carrier); err == nil {
		span = opentracing.StartSpan(operation, ext.RPCServerOption(spanCtx))
		ctx = opentracing.ContextWithSpan(ctx, span)
	} else {
		span, ctx = opentracing.StartSpanFromContext(ctx, operation)
	}

	ext.SpanKindRPCServer.Set(span)
	span.SetTag(string(ext.HTTPUrl), r.URL.Path)

	return r.WithContext(ctx), span
}

func (s panicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r, span := startSpan(r)
	defer func() {
		if rec := recover(); rec != nil {
			ctx := r.Context()
			ctx = ctxkit.WithTraceID(ctx, trace.GetTraceID(ctx))
			log.Get(ctx).Error(rec, string(debug.Stack()))
		}
		span.Finish()
	}()
	isReturn := initCROSHeaders(w, r)
	if isReturn {
		return
	}
	r = initRequestHeaders(w, r)
	s.handler.ServeHTTP(w, r)
}

func initTables() {
	autoCreate := conf.GetBool("DB_DEFAULT_AUTO_CREATE_TABLE")
	if autoCreate {
		service.SyncXORMTables()
	}
}

func initSentinel() {
	sentinelConfigFilePath := conf.GetString("SENTINEL_CONFIG_PATH")
	if len(sentinelConfigFilePath) == 0 {
		sentinelConfigFilePath = filepath.Join(conf.GetConfigPath(), "sentinel.yml")
	}
	err := sentinel.InitWithConfigFile(sentinelConfigFilePath)
	if err != nil {
		panic(fmt.Sprintf("missing sentinel config file:%s", sentinelConfigFilePath))
	}
	qpsMaps := conf.GetStrMapStr("SENTINEL_RES_QPS")
	if len(qpsMaps) > 0 {
		rules := make([]*flow.Rule, 0, 0)
		for resName, resQPSStr := range qpsMaps {
			rule := &flow.Rule{
				Resource:               resName,
				Threshold:              cast.ToFloat64(resQPSStr),
				TokenCalculateStrategy: flow.Direct,
				ControlBehavior:        flow.Reject,
			}
			rules = append(rules, rule)
		}
		_, err = flow.LoadRules(rules)
		if err != nil {
			panic(err)
		}
	}
}

func startServer() {
	logger.Info("start server")

	rand.Seed(int64(time.Now().Nanosecond()))

	mux := http.NewServeMux()

	timeout := 60 * time.Second
	initMux(mux, isInternal)
	if isInternal {
		initInternalMux(mux)

		if d := conf.GetDuration("INTERNAL_API_TIMEOUT"); d > 0 {
			timeout = d * time.Second
		}
	} else {
		if d := conf.GetDuration("OUTER_API_TIMEOUT"); d > 0 {
			timeout = d * time.Second
		}
	}

	handler := http.TimeoutHandler(panicHandler{handler: mux}, timeout, "timeout")

	http.Handle("/", handler)

	metricsHandler := promhttp.Handler()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		util.GatherMetrics()
		metricsHandler.ServeHTTP(w, r)
	})

	http.HandleFunc("/monitor/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	serverHttp := conf.GetString("SERVER_HTTP_IP")

	if len(serverHttp) == 0 {
		serverHttp = "127.0.0.1"
	}

	addr := fmt.Sprintf("%s:%d", serverHttp, port)
	server = &http.Server{
		IdleTimeout: 120 * time.Second,
	}

	// 配置下发可能会多次触发重启，必须等待 Listen() 调用成功
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		// 本段代码基本搬自 http 标准库
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			panic(err)
		}
		wg.Done()

		err = server.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
		if err != http.ErrServerClosed {
			panic(err)
		}
	}()

	wg.Wait()
}

func stopServer() {
	logger.Info("stop server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}
	util.Reset()
}

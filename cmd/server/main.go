package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/busyfree/leaf-go/util"
	"github.com/busyfree/leaf-go/util/conf"
	_ "github.com/busyfree/leaf-go/util/redis"
)

func main() {
	reload := make(chan int, 1)
	stop := make(chan os.Signal, 1)
	conf.OnConfigChange(func() { reload <- 1 })
	conf.WatchConfig()
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	initTables()
	initSentinel()
	startServer()

	for {
		select {
		case <-reload:
			util.Reset()
		case sg := <-stop:
			stopServer()
			// 仿 nginx 使用 HUP 信号重载配置
			if sg == syscall.SIGHUP {
				startServer()
			} else {
				util.Stop()
				return
			}
		}
	}
}

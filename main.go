

package main

import (
	_ "net/http/pprof" // 注册 pprof 接口

	"github.com/busyfree/leaf-go/cmd/job"
	"github.com/busyfree/leaf-go/cmd/server"
	"github.com/busyfree/leaf-go/cmd/version"
	"github.com/busyfree/leaf-go/util/conf"

	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs" // 根据容器配额设置 maxprocs
)

var (
	a string
	v string
	c string
	d string
)

func main() {
	conf.BinAppName = a
	conf.BinBuildCommit = c
	conf.BinBuildVersion = v
	conf.BinBuildDate = d
	root := cobra.Command{Use: "srv_treedom_web"}
	root.AddCommand(
		server.Cmd,
		job.Cmd,
		version.Cmd,
	)
	root.Execute()
}

package main

import (
	_ "net/http/pprof" // 注册 pprof 接口

	"github.com/busyfree/leaf-go/cmd/job"
	"github.com/busyfree/leaf-go/cmd/server"
	"github.com/busyfree/leaf-go/cmd/version"
	"github.com/busyfree/leaf-go/util/conf"

	"github.com/spf13/cobra"
	"go.uber.org/automaxprocs/maxprocs"
)

var (
	a string
	v string
	c string
	d string
)

func main() {
	nopLog := func(string, ...interface{}) {}
	maxprocs.Set(maxprocs.Logger(nopLog))
	conf.BinAppName = a
	conf.BinBuildCommit = c
	conf.BinBuildVersion = v
	conf.BinBuildDate = d
	root := cobra.Command{Use: "leaf_go"}
	root.AddCommand(
		server.Cmd,
		job.Cmd,
		version.Cmd,
	)
	root.Execute()
}

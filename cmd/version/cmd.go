package version

import (
	"fmt"

	"github.com/busyfree/leaf-go/util/conf"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "version",
	Long:  `version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("AppName   : %s\n", conf.BinAppName)
		fmt.Printf("Version   : %s\n", conf.BinBuildVersion)
		fmt.Printf("Commit    : %s\n", conf.BinBuildCommit)
		fmt.Printf("BuildDate : %s\n", conf.BinBuildDate)
	},
}

package cmd

import (
	"fmt"

	"github.com/infraboard/mcube/v2/ioc/server"
	"github.com/spf13/cobra"
)

// Root represents the base command when called without any subcommands
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动服务",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("xx")
		cobra.CheckErr(server.Run(cmd.Context()))
	},
}

func init() {
	Root.AddCommand(startCmd)
}

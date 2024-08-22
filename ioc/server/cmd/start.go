package cmd

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc/server"
	"github.com/spf13/cobra"
)

// Root represents the base command when called without any subcommands
var startCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(server.Run(context.Background()))
	},
}

package initial

import (
	"github.com/spf13/cobra"
)

// Cmd represents the start command
var Cmd = &cobra.Command{
	Use:   "init",
	Short: "mpaas 服务初始化",
	Long:  "mpaas 服务初始化",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

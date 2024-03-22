package module

import (
	"github.com/spf13/cobra"
)

// InitCmd 初始化系统
var Cmd = &cobra.Command{
	Use:   "init",
	Short: "初始化",
	Long:  `初始化一个mcube项目`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

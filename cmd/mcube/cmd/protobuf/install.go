package protobuf

import (
	"github.com/spf13/cobra"
)

// InitCmd 初始化系统
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "安装项目依赖的protobuf文件",
	Long:  `安装项目依赖的protobuf文件`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	Cmd.AddCommand(installCmd)
}

package project

import (
	"github.com/spf13/cobra"
)

// ProjectCmd 初始化系统
var ProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "项目初始化工具",
	Long:  `项目初始化工具`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

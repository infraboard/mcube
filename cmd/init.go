package cmd

import (
	"github.com/spf13/cobra"

	"github.com/infraboard/mcube/cmd/project"
)

// InitCmd 初始化系统
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化",
	Long:  `初始化一个mcube项目`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := project.LoadConfigFromCLI()
		if err != nil {
			return err
		}

		err = p.Init()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(InitCmd)
}

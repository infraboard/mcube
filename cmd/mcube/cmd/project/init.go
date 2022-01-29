package project

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/infraboard/mcube/cmd/mcube/project"
)

// InitCmd 初始化系统
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化",
	Long:  `初始化一个mcube项目`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := project.LoadConfigFromCLI()
		if err != nil {
			if err == terminal.InterruptErr {
				fmt.Println("项目初始化取消")
				return nil
			}
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
	Cmd.AddCommand(initCmd)
}

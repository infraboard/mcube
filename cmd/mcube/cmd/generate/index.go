package generate

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "generate",
	Short: "代码生成器",
	Long:  `代码生成器`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

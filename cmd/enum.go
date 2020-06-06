package cmd

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/infraboard/mcube/cmd/enum"
)

// EnumCmd 枚举生成器
var EnumCmd = &cobra.Command{
	Use:   "enum",
	Short: "枚举生成器",
	Long:  `枚举生成器`,
	RunE: func(cmd *cobra.Command, args []string) error {
		code, err := enum.G.Generate()
		if err != nil {
			return err
		}

		fname := os.Getenv("GOFILE")

		genFile := fname[0:len(fname)-len(path.Ext(fname))] + "_generate.go"
		ioutil.WriteFile(genFile, code, 0644)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(EnumCmd)
}

func init() {
	EnumCmd.PersistentFlags().BoolVarP(&enum.G.Marshal, "marshal", "m", false, "is generate json MarshalJSON and UnmarshalJSON method")
}

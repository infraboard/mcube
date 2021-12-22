package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/infraboard/mcube/cmd/mcube/protoc_tools/inject_tag"
	"github.com/spf13/cobra"
)

// EnumCmd 枚举生成器
var InjectTagCmd = &cobra.Command{
	Use:   "inject-tag",
	Short: "tag 注入",
	Long:  `tag 注入`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("input file is mandatory, see: -help")
		}

		matchedFiles := []string{}
		for _, v := range args {
			files, err := filepath.Glob(v)
			if err != nil {
				return err
			}
			matchedFiles = append(matchedFiles, files...)
		}

		if len(matchedFiles) == 0 {
			return fmt.Errorf("no file matched")
		}

		for _, path := range matchedFiles {
			// It should end with ".go" at a minimum.
			if !strings.HasSuffix(path, ".go") {
				continue
			}

			areas, err := inject_tag.ParseFile(path, []string{})
			if err != nil {
				log.Fatal(err)
			}

			if err = inject_tag.WriteFile(path, areas); err != nil {
				log.Fatal(err)
			}
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(InjectTagCmd)
}

func init() {
}

package cmd

// RootTemplate todo
const RootTemplate = `package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	
	"{{.PKG}}/version"
)

var vers bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "{{.Name}}",
	Short: "{{.Description}}",
	Long:  "{{.Description}}",
	RunE: func(cmd *cobra.Command, args []string) error {
		if vers {
			fmt.Println(version.FullVersion())
			return nil
		}
		return errors.New("no flags find")
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&vers, "version", "v", false, "the {{.Name}} version")
}`

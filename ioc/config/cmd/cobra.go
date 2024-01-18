package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
)

var (
	confType string
	confFile string
	vers     bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   application.Get().AppName,
	Short: application.Get().AppDescription,
	Long:  application.Get().AppDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		if vers {
			fmt.Println(application.FullVersion())
			return nil
		}
		return cmd.Help()
	},
}

// 初始化Ioc
func initail() {
	req := ioc.NewLoadConfigRequest()
	switch confType {
	case "file":
		req.ConfigFile.Enabled = true
		req.ConfigFile.Path = confFile
	default:
		req.ConfigEnv.Enabled = true
	}
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// 初始化设置
	cobra.OnInitialize(initail)
	err := RootCmd.Execute()
	cobra.CheckErr(err)
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&confType, "config-type", "t", "file", "the service config type [file/env]")
	RootCmd.PersistentFlags().StringVarP(&confFile, "config-file", "f", "etc/application.toml", "the service config from file")
	RootCmd.PersistentFlags().BoolVarP(&vers, "version", "v", false, "the mcenter version")
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/server"
)

var (
	confType string
	confFile string
	vers     bool
)

// Root represents the base command when called without any subcommands
var Root = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		if vers {
			fmt.Println(application.FullVersion())
			return nil
		}
		return cmd.Help()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.OnInitialize(func() {
		req := server.DefaultConfig
		switch confType {
		case "file":
			req.ConfigFile.Enabled = true
			req.ConfigFile.Path = confFile
		default:
			req.ConfigEnv.Enabled = true
		}

		// 初始化ioc
		cobra.CheckErr(ioc.ConfigIocObject(server.DefaultConfig))

		// 补充Root命令信息
		Root.Use = application.Get().AppName
		Root.Short = application.Get().AppDescription
		Root.Long = application.Get().AppDescription
	})
	cobra.CheckErr(Root.Execute())
}

// 补充启动命令的CLI
func Start() {
	Root.AddCommand(startCmd)
	Execute()
}

func init() {
	Root.PersistentFlags().StringVarP(&confType, "config-type", "t", "file", "the service config type [file/env]")
	Root.PersistentFlags().StringVarP(&confFile, "config-file", "f", "etc/application.toml", "the service config from file")
	Root.PersistentFlags().BoolVarP(&vers, "version", "v", false, "the service version")
}

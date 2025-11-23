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
	debug    bool
)

// Root represents the base command when called without any subcommands
var Root = &cobra.Command{
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if vers {
			return nil
		}
		req := server.DefaultConfig
		switch confType {
		case "file":
			req.ConfigFile.Enabled = true
			req.ConfigFile.Path = confFile
		default:
			req.ConfigEnv.Enabled = true
		}

		// 设置debug
		ioc.SetDebug(debug)

		// 初始化ioc
		cobra.CheckErr(ioc.ConfigIocObject(server.DefaultConfig))

		// 补充Root命令信息
		cmd.Use = application.Get().AppName
		cmd.Short = application.Get().AppDescription
		cmd.Long = application.Get().AppDescription

		return nil
	},
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
	cobra.CheckErr(Root.Execute())
}

// 补充启动命令的CLI
func Start() {
	Execute()
}

func init() {
	Root.PersistentFlags().StringVarP(&confType, "config-type", "t", "file", "the service config type [file/env]")
	Root.PersistentFlags().StringVarP(&confFile, "config-file", "f", "etc/application.toml", "the service config from file")
	Root.PersistentFlags().BoolVarP(&vers, "version", "v", false, "the service version")
	Root.PersistentFlags().BoolVarP(&debug, "debug", "d", true, "debug mode")
}

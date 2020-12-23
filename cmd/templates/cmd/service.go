package cmd

// ServiceTemplate todo
const ServiceTemplate = `package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/spf13/cobra"

	"{{.PKG}}/api"
	"{{.PKG}}/conf"
	"{{.PKG}}/pkg"
	
	_ "github.com/go-sql-driver/mysql"
	// 加载服务模块
)

var (
	// pusher service config option
	confType string
	confFile string
	confEtcd string
)

// startCmd represents the start command
var serviceCmd = &cobra.Command{
	Use:   "service [start/stop/reload/restart]",
	Short: "{{.Name}} API服务",
	Long:  "{{.Name}} API服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("[start] are required")
		}
		// 初始化全局变量
		if err := loadGlobalConfig(confType); err != nil {
			return err
		}
		// 初始化全局日志配置
		if err := loadGlobalLogger(); err != nil {
			return err
		}
		// 初始化服务层
		if err := pkg.InitService(); err != nil {
			return err
		}
		conf := conf.C()
		switch args[0] {
		case "start":
			// 启动服务
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
			// 初始化服务
			svr, err := newService(conf)
			if err != nil {
				return err
			}
			// 等待信号处理
			go svr.waitSign(ch)
			// 启动服务
			if err := svr.start(); err != nil {
				if !strings.Contains(err.Error(), "http: Server closed") {
					return err
				}
			}
		default:
			return errors.New("not support argument, support [start]")
		}
		return nil
	},
}

func newService(cnf *conf.Config) (*service, error) {
	http := api.NewHTTPService()
	svr := &service{
		http: http,
		log:  zap.L().Named("CLI"),
	}
	return svr, nil
}

type service struct {
	http *api.HTTPService
	log  logger.Logger
	stop context.CancelFunc
}

func (s *service) start() error {
	s.log.Infof("loaded services: %v", pkg.LoadedService())
	return s.http.Start()
}

// config 为全局变量, 只需要load 即可全局可用户
func loadGlobalConfig(configType string) error {
	// 配置加载
	switch configType {
	case "file":
		err := conf.LoadConfigFromToml(confFile)
		if err != nil {
			return err
		}
	case "env":
		err := conf.LoadConfigFromEnv()
		if err != nil {
			return err
		}
	case "etcd":
		return errors.New("not implemented")
	default:
		return errors.New("unknown config type")
	}
	return nil
}

// log 为全局变量, 只需要load 即可全局可用户, 依赖全局配置先初始化
func loadGlobalLogger() error {
	var (
		logInitMsg string
		level      zap.Level
	)
	lc := conf.C().Log
	lv, err := zap.NewLevel(lc.Level)
	if err != nil {
		logInitMsg = fmt.Sprintf("%s, use default level INFO", err)
		level = zap.InfoLevel
	} else {
		level = lv
		logInitMsg = fmt.Sprintf("log level: %s", lv)
	}
	zapConfig := zap.DefaultConfig()
	zapConfig.Level = level
	switch lc.To {
	case conf.ToStdout:
		zapConfig.ToStderr = true
		zapConfig.ToFiles = false
	case conf.ToFile:
		zapConfig.Files.Name = "api.log"
		zapConfig.Files.Path = lc.PathDir
	}
	switch lc.Format {
	case conf.JSONFormat:
		zapConfig.JSON = true
	}
	if err := zap.Configure(zapConfig); err != nil {
		return err
	}
	zap.L().Named("INIT").Info(logInitMsg)
	return nil
}
func (s *service) waitSign(sign chan os.Signal) {
	for {
		select {
		case sg := <-sign:
			switch v := sg.(type) {
			default:
				s.log.Infof("receive signal '%v', start graceful shutdown", v.String())
				if err := s.http.Stop(); err != nil {
					s.log.Errorf("graceful shutdown err: %s, force exit", err)
				}
				s.log.Infof("service stop complete")
				return
			}
		}
	}
}
func init() {
	serviceCmd.Flags().StringVarP(&confType, "config-type", "t", "file", "the service config type [file/env/etcd]")
	serviceCmd.Flags().StringVarP(&confFile, "config-file", "f", "etc/{{.Name}}.toml", "the service config from file")
	serviceCmd.Flags().StringVarP(&confEtcd, "config-etcd", "e", "127.0.0.1:2379", "the service config from etcd")
	RootCmd.AddCommand(serviceCmd)
}`

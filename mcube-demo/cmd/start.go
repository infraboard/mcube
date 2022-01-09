package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/cache"
	"github.com/infraboard/mcube/cache/memory"
	"github.com/infraboard/mcube/cache/redis"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"gitee.com/go-course/mcube-demo/conf"
	"gitee.com/go-course/mcube-demo/protocol"

	// 注册所有服务
	_ "gitee.com/go-course/mcube-demo/apps/all"
)

// startCmd represents the start command
var serviceCmd = &cobra.Command{
	Use:   "start",
	Short: "mcube-demo API服务",
	Long:  "mcube-demo API服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 初始化全局变量
		if err := loadGlobalConfig(confType); err != nil {
			return err
		}

		// 初始化全局日志配置
		if err := loadGlobalLogger(); err != nil {
			return err
		}

		// 加载缓存
		if err := loadCache(); err != nil {
			return err
		}

		// 初始化全局app
		if err := app.InitAllApp(); err != nil {
			return err
		}

		conf := conf.C()
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

		return nil
	},
}

func newService(cnf *conf.Config) (*service, error) {
	http := protocol.NewHTTPService()
	grpc := protocol.NewGRPCService()
	svr := &service{
		http: http,
		grpc: grpc,
		log:  zap.L().Named("CLI"),
	}

	return svr, nil
}

type service struct {
	http *protocol.HTTPService
	grpc *protocol.GRPCService

	log  logger.Logger
	stop context.CancelFunc
}

func (s *service) start() error {
	s.log.Infof("loaded grpc app: %s", app.LoadedGrpcApp())
	s.log.Infof("loaded http app: %s", app.LoadedHttpApp())
	s.log.Infof("loaded internal app: %s", app.LoadedInternalApp())

	go s.grpc.Start()
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

func loadCache() error {
	l := zap.L().Named("INIT")
	c := conf.C()
	// 设置全局缓存
	switch c.Cache.Type {
	case "memory", "":
		ins := memory.NewCache(c.Cache.Memory)
		cache.SetGlobal(ins)
		l.Info("use cache in local memory")
	case "redis":
		ins := redis.NewCache(c.Cache.Redis)
		cache.SetGlobal(ins)
		l.Info("use redis to cache")
	default:
		return fmt.Errorf("unknown cache type: %s", c.Cache.Type)
	}

	return nil
}

func (s *service) waitSign(sign chan os.Signal) {
	for {
		select {
		case sg := <-sign:
			switch v := sg.(type) {
			default:
				s.log.Infof("receive signal '%v', start graceful shutdown", v.String())

				if err := s.grpc.Stop(); err != nil {
					s.log.Errorf("grpc graceful shutdown err: %s, force exit", err)
				} else {
					s.log.Info("grpc service stop complete")
				}
				
				if err := s.http.Stop(); err != nil {
					s.log.Errorf("http graceful shutdown err: %s, force exit", err)
				} else {
					s.log.Infof("http service stop complete")
				}
				return
			}
		}
	}
}

func init() {
	RootCmd.AddCommand(serviceCmd)
}
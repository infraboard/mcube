# mcube
![Build and Test](https://github.com/infraboard/mcube/v2/workflows/Build%20and%20Test/badge.svg)
[![codecov](https://codecov.io/gh/infraboard/mcube/branch/master/graph/badge.svg)](https://codecov.io/gh/infraboard/mcube)
[![Go Report Card](https://goreportcard.com/badge/github.com/infraboard/mcube/v2)](https://goreportcard.com/report/github.com/infraboard/mcube/v2)
[![Release](https://img.shields.io/github/release/infraboard/mcube.svg?style=flat-square)](https://github.com/infraboard/mcube/v2/releases)
[![MIT License](https://img.shields.io/github/license/infraboard/mcube.svg)](https://github.com/infraboard/mcube/v2/blob/master/LICENSE)

[官方文档](https://www.mcube.top/docs/framework/)

mcube是一款用于构建渐进式微服务(单体-->微服务)的框架, 让应用从单体无缝过渡到微服务, 同时提供丰富的配置即用的功能配置, 
只需简单配置就可拥有:
+ Log: 支持文件滚动和Trace的日志打印
+ Metric: 支持应用自定义指标监控
+ Trace: 集成支持完整的全链路追踪(HTTP Server/GRPC Server/数据库...)以及自定义埋点
+ CORS: 资源跨域共享
+ Health Check: HTTP 和 GRPC 健康检查
+ API DOC: 基于Swagger的 API 文档

除了上面这些功能配置，还会用到很多三方工具, 也是配置即用:
+ MySQL: Grom集成
+ MongoDB 官方驱动集成
+ Redis: go-redis集成
+ Kafka: kafka-go集成
+ 分布式缓存: 当前只适配了Redis
+ 分布式锁: 当前只适配了Redis

![框架架构](./docs/ioc/arch.png)

## 快速开始

下面是演示一个TestObject对象的注册与获取的基础功能:
```go
package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/ioc/server"
	"gorm.io/gorm"
)

func main() {
	// 注册HTTP接口类
	ioc.Api().Registry(&ApiHandler{})

	// 开启配置文件读取配置
	server.DefaultConfig.ConfigFile.Enabled = true
	server.DefaultConfig.ConfigFile.Path = "etc/application.toml"

	// 启动应用
	err := server.Run(context.Background())
	if err != nil {
		panic(err)
	}
}

type ApiHandler struct {
	// 继承自Ioc对象
	ioc.ObjectImpl

	// mysql db依赖
	db *gorm.DB
}

// 覆写对象的名称, 该名称名称会体现在API的路径前缀里面
// 比如: /simple/api/v1/module_a/db_stats
// 其中/simple/api/v1/module_a 就是对象API前缀, 命名规则如下:
// <service_name>/<path_prefix>/<object_version>/<object_name>
func (h *ApiHandler) Name() string {
	return "module_a"
}

// 初始化db属性, 从ioc的配置区域获取共用工具 gorm db对象
func (h *ApiHandler) Init() error {
	h.db = datasource.DB()
	return nil
}

// API路由
func (h *ApiHandler) Registry(r gin.IRouter) {
	r.GET("/db_stats", func(ctx *gin.Context) {
		db, _ := h.db.DB()
		ctx.JSON(http.StatusOK, gin.H{
			"data": db.Stats(),
		})
	})
}
```

## 应用开发

### 标准化工程配置

统一了项目的配置加载方式:

环境变量
配置文件
TOML
YAML
JSON
下面是项目配置文件(etc/application.toml)内容:

```toml
[app]
name = "simple"
key  = "this is your app key"

[http]
host = "127.0.0.1"
port = 8020

[datasource]
host = "127.0.0.1"
port = 3306
username = "root"
password = "123456"
database = "test"

[log]
level = "debug"

[log.file]
enable = true
file_path = "logs/app.log"
```

### 即插即用的组件

通过简单的配置就能为项目添加:

检查检查(Health Chcek)
应用指标监控(Metric)

```go
import (
  // 开启Health健康检查
  _ "github.com/infraboard/mcube/v2/ioc/apps/health/gin"
  // 开启Metric
  _ "github.com/infraboard/mcube/v2/ioc/apps/metric/gin"
)
```

启动过后, 在日志里就能看到这2个功能开启了:
```sh
2024-01-05T11:30:00+08:00 INFO   health/gin/check.go:52 > Get the Health using http://127.0.0.1:8020/healthz component:HEALTH_CHECK
2024-01-05T11:30:00+08:00 INFO   metric/gin/metric.go:51 > Get the Metric using http://127.0.0.1:8020/metrics component:METRIC
```

当然你也可以通过配置来修改功能的URL路径:
```toml
[health]
  path = "/healthz"

[metric]
  enable = true
  provider = "prometheus"
  endpoint = "/metrics"
```
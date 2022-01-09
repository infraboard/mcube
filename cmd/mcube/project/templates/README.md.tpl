# {{.Name}}

{{.Description}}

## 架构图

## 项目说明

{{.Backquote3}}
├── protocol                       # 脚手架功能: rpc / http 功能加载
│   ├── grpc.go              
│   └── http.go    
├── client                         # 脚手架功能: grpc 客户端实现 
│   ├── client.go              
│   └── config.go    
├── cmd                            # 脚手架功能: 处理程序启停参数，加载系统配置文件
│   ├── root.go             
│   └── start.go                
├── conf                           # 脚手架功能: 配置文件加载
│   ├── config.go                  # 配置文件定义
│   ├── load.go                    # 不同的配置加载方式
│   └── log.go                     # 日志配置文件
├── dist                           # 脚手架功能: 构建产物
├── etc                            # 配置文件
│   ├── xxx.env
│   └── xxx.toml
├── pkg                            # 具体业务场景的领域包
│   ├── all                    
│   │   └── all.go                 # 手动配置: 加载所有领域模块里面的http和rpc服务。 
│   ├── example                    # 具体业务场景领域服务 example
│   │   ├── http                   # http 
│   │   │    ├── example.go        # example 服务的http方法实现，请求参数处理、权限处理、数据响应等 
│   │   │    └── http.go           # 领域模块内的 http 路由处理，向系统层注册http服务
│   │   ├── impl                   # rpc
│   │   │    ├──  example.go       # example 服务的rpc方法实现，请求参数处理、权限处理、数据响应等 
│   │   │    └── impl.go           # 领域模块内的 rpc 服务注册 ，向系统层注册rpc服务
│   │   ├──  pb                    # protobuf 定义
│   │   │     ├── response.proto   # example 服务数据模型定义
│   │   │     ├── request.proto    # example 服务请求结构体定义
│   │   │     └── service.proto    # example 服务接口定义
│   │   ├── request_ext.go         # request 结构体初始化
│   │   ├── reponse_ext.go         # reponse 结构体初始化
│   │   ├── request.pb.go          # pb/request.proto 生成的结构体和默认方法
│   │   ├── reponse.pb.go          # pb/reponse.proto 生成的结构体和默认方法
│   │   ├── service_grpc.pb.go     # pb/service.proto 生成grpc方法
│   │   ├── service_http.pb.go     # pb/service.proto 生成http方法
│   │   └── service.pb.go          # pb/service.proto 生成方法定义
├── version                        # 程序版本信息
│   └── version.go                    
├── README.md                    
├── main.go                        # Go程序唯一入口
├── Makefile                       # make 命令定义
└── go.mod                         # go mod 依赖定义
{{.Backquote3}}


## 快速开发
make脚手架
{{.Backquote3}}sh
➜  {{.Name}} git:(master) ✗ make help
dep                            Get the dependencies
lint                           Lint Golang files
vet                            Run go vet
test                           Run unittests
test-coverage                  Run tests with coverage
build                          Local build
linux                          Linux build
run                            Run Server
clean                          Remove previous build
help                           Display this help screen
{{.Backquote3}}

1. 使用go mod下载项目依赖
{{.Backquote3}}sh
$ make dep
{{.Backquote3}}

2. 添加配置文件(默认读取位置: etc/{{.Name}}.toml)
{{.Backquote3}}sh
$ 编辑样例配置文件 etc/{{.Name}}.toml.example
$ mv etc/{{.Name}}.toml.example etc/{{.Name}}.toml
{{.Backquote3}}

3. 启动服务
{{.Backquote3}}sh
$ make run
{{.Backquote3}}

## 相关文档
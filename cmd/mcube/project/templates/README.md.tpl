# {{.Name}}

{{.Description}}

## 架构图

## 项目说明

```
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
│   │   |-- grpc.go                # 注册所有GRPC服务模块, 暴露给框架GRPC服务器加载, 注意 导入有先后顺序。  
│   │   |-- http.go                # 注册所有HTTP服务模块, 暴露给框架HTTP服务器加载。                    
│   │   └── internal.go            #  注册所有内部服务模块, 无须对外暴露的服务, 用于内部依赖。 
│   ├── book                       # 具体业务场景领域服务 book
│   │   ├── http                   # http 
│   │   │    ├── book.go           # book 服务的http方法实现，请求参数处理、权限处理、数据响应等 
│   │   │    └── http.go           # 领域模块内的 http 路由处理，向系统层注册http服务
│   │   ├── impl                   # rpc
│   │   │    ├── book.go          # book 服务的rpc方法实现，请求参数处理、权限处理、数据响应等 
│   │   │    └── impl.go           # 领域模块内的 rpc 服务注册 ，向系统层注册rpc服务
│   │   ├──  pb                    # protobuf 定义
│   │   │     └── book.proto       # book proto 定义文件
│   │   ├── app.go                 # book app 只定义扩展
│   │   ├── book.pb.go             # protobuf 生成的文件
│   │   └── book_grpc.pb.go        # pb/book.proto 生成方法定义
├── version                        # 程序版本信息
│   └── version.go                    
├── README.md                    
├── main.go                        # Go程序唯一入口
├── Makefile                       # make 命令定义
└── go.mod                         # go mod 依赖定义
```


## 快速开发
make脚手架
```sh
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
```

1. 使用go mod下载项目依赖
```sh
$ make dep
```

2. 添加配置文件(默认读取位置: etc/{{.Name}}.toml)
```sh
$ 编辑样例配置文件 etc/{{.Name}}.toml.book
$ mv etc/{{.Name}}.toml.book etc/{{.Name}}.toml
```

3. 启动服务
```sh
$ make run
```

## 相关文档
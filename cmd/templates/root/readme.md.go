package root

// ReadmeTemplate tooo
const ReadmeTemplate = `# {{.Name}}

{{.Description}}


## 架构图


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

## 相关文档`

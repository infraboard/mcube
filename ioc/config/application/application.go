package application

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/infraboard/mcube/v2/ioc"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Application{
	AppName:      "",
	AppGroup:     "default",
	Domain:       "localhost",
	Security:     false,
	EncryptKey:   "defualt app encrypt key",
	CipherPrefix: "@ciphered@",
}

type Application struct {
	ioc.ObjectImpl

	AppGroup       string `json:"group" yaml:"group" toml:"group" env:"GROUP"`
	AppName        string `json:"name" yaml:"name" toml:"name" env:"NAME"`
	AppDescription string `json:"description" yaml:"description" toml:"description" env:"DESCRIPTION"`
	Domain         string `json:"domain" yaml:"domain" toml:"domain" env:"DOMAIN"`
	Security       bool   `json:"security" yaml:"security" toml:"security" env:"SECURITY"`
	EncryptKey     string `json:"encrypt_key" yaml:"encrypt_key" toml:"encrypt_key" env:"ENCRYPT_KEY"`
	CipherPrefix   string `json:"cipher_prefix" yaml:"cipher_prefix" toml:"cipher_prefix" env:"CIPHER_PREFIX"`
}

func (i *Application) GetAppName() string {
	if i.AppName != "" {
		return i.AppName
	}

	// 获取当前可执行文件的路径
	executablePath, err := os.Executable()
	if err != nil {
		return ""
	}

	// 获取可执行文件的名称
	return filepath.Base(executablePath)
}

func (i *Application) IsInternalIP() bool {
	return i.Domain != "localhost"
}

func (i *Application) Endpoint() string {
	schema := "http"
	if i.Security {
		schema = "https"
	}
	return fmt.Sprintf("%s://%s", schema, i.Domain)
}

func (i *Application) Init() error {
	sn := os.Getenv("OTEL_SERVICE_NAME")
	if sn == "" {
		os.Setenv("OTEL_SERVICE_NAME", i.AppName)
	}
	return nil
}

func (i *Application) Name() string {
	return AppName
}

// 优先初始化, 以供后面的组件使用
func (i *Application) Priority() int {
	return 999
}

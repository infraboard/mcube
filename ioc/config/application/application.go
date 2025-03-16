package application

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/infraboard/mcube/v2/crypto/cbc"
	"github.com/infraboard/mcube/v2/ioc"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Application{
	AppName:      "",
	AppGroup:     "default",
	EncryptKey:   "defualt app encrypt key",
	CipherPrefix: "@ciphered@",
	Extras:       map[string]string{},
}

type Application struct {
	ioc.ObjectImpl

	AppGroup        string            `json:"group" yaml:"group" toml:"group" env:"GROUP"`
	AppName         string            `json:"name" yaml:"name" toml:"name" env:"NAME"`
	AppDescription  string            `json:"description" yaml:"description" toml:"description" env:"DESCRIPTION"`
	AppAddress      string            `json:"address" yaml:"address" toml:"address" env:"ADDRESS"`
	InternalAddress string            `json:"internal_address" yaml:"internal_address" toml:"internal_address" env:"INTERNAL_ADDRESS"`
	InternalToken   string            `json:"internal_token" yaml:"internal_token" toml:"internal_token" env:"INTERNAL_TOKEN"`
	EncryptKey      string            `json:"encrypt_key" yaml:"encrypt_key" toml:"encrypt_key" env:"ENCRYPT_KEY"`
	CipherPrefix    string            `json:"cipher_prefix" yaml:"cipher_prefix" toml:"cipher_prefix" env:"CIPHER_PREFIX"`
	Extras          map[string]string `json:"extras" yaml:"extras" toml:"extras" env:"EXTRAS"`

	appURL *url.URL
}

func (i *Application) GetExtras(key string) string {
	if i.Extras == nil {
		i.Extras = map[string]string{}
	}
	return i.Extras[key]
}

func (i *Application) Domain() string {
	if i.appURL != nil {
		return strings.Split(i.appURL.Host, ":")[0]
	}
	return "localhost"
}

func (i *Application) GenCBCEncryptKey() []byte {
	return cbc.GenKeyBySHA1Hash2(i.EncryptKey, cbc.AES_KEY_LEN_32)
}

func (i *Application) Host() string {
	if i.appURL != nil {
		return i.appURL.Host
	}
	return ""
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
	return i.Domain() != "localhost" && i.Domain() != "" && i.Domain() != "127.0.0.1"
}

func (i *Application) Init() error {
	sn := os.Getenv("OTEL_SERVICE_NAME")
	if sn == "" {
		os.Setenv("OTEL_SERVICE_NAME", i.AppName)
	}

	if i.AppAddress != "" {
		v, err := url.Parse(i.AppAddress)
		if err != nil {
			return err
		}
		i.appURL = v
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

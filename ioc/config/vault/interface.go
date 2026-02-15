package vault

import (
	"github.com/hashicorp/vault-client-go"
	"github.com/infraboard/mcube/v2/ioc"
)

const (
	AppName = "vault"
)

// Client 返回 vault 客户端实例
func Client() *vault.Client {
	return Get().client
}

// Get 获取 Vault 配置实例
func Get() *Vault {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*Vault)
}

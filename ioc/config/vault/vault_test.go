package vault_test

import (
	"testing"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/vault"
	"github.com/infraboard/mcube/v2/tools/file"
)

func TestVaultClient(t *testing.T) {
	client := vault.Client()
	if client != nil {
		t.Log("vault client initialized")
		// 示例：读取秘密
		// secret, err := client.Secrets.KvV2Read(context.Background(), "my-secret")
		// if err != nil {
		//     t.Fatal(err)
		// }
		// t.Log(secret.Data)
	}
}

func TestDefaultConfig(t *testing.T) {
	file.MustToToml(
		vault.AppName,
		ioc.Config().Get(vault.AppName),
		"test/default.toml",
	)
}

func init() {
	// 初始化 IOC 容器
	// 如果需要测试，请在 test/ 目录下创建 application.toml 配置文件
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Paths = []string{"test/application.toml"}
	req.ConfigFile.SkipIFNotExist = true

	err := ioc.ConfigIocObject(req)
	if err != nil {
		panic(err)
	}
}

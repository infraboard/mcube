package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault-client-go"
	"github.com/infraboard/mcube/v2/ioc"
	iocvault "github.com/infraboard/mcube/v2/ioc/config/vault"

	// 引入 vault 模块，使其自动注册
	_ "github.com/infraboard/mcube/v2/ioc/config/vault"
)

func main() {
	// 初始化 IOC 容器，会自动加载配置并初始化 vault 客户端
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Paths = []string{"etc/application.toml"}
	err := ioc.ConfigIocObject(req)
	if err != nil {
		panic(err)
	}

	// 获取 Vault 客户端
	client := iocvault.Client()
	fmt.Printf("Vault client initialized: %v\n", client != nil)
	fmt.Printf("KV Mount Path: %s\n", iocvault.KvMountPath())

	ctx := context.Background()

	// ==================== 方式 1: 使用辅助方法（推荐）====================
	fmt.Println("\n========== 方式 1: 使用辅助方法（自动使用配置的挂载点）==========")
	resp, err := iocvault.ReadSecret(ctx, "myapp/config")
	if err != nil {
		fmt.Printf("Error reading secret: %v\n", err)
		fmt.Println("\n提示：请先运行以下命令创建测试秘密：")
		fmt.Println("  docker-compose up -d")
		fmt.Println("  export VAULT_ADDR='http://127.0.0.1:8200'")
		fmt.Println("  export VAULT_TOKEN='myroot'")
		fmt.Println("  vault kv put secret/myapp/config database_url='postgres://localhost/db' api_key='secret'")
		return
	}

	fmt.Printf("Secret data: %+v\n", resp.Data.Data)
	fmt.Println("\nSecret fields:")
	for k, v := range resp.Data.Data {
		fmt.Printf("  %s: %v\n", k, v)
	}

	// ==================== 方式 2: 手动指定挂载点 ====================
	fmt.Println("\n========== 方式 2: 手动指定挂载点（灵活但繁琐）==========")
	resp2, err := client.Secrets.KvV2Read(
		ctx,
		"myapp/config",                // 秘密路径
		vault.WithMountPath("secret"), // 手动指定挂载点
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Same result: %v\n", resp2.Data.Data)

	// ==================== 写入秘密示例 ====================
	fmt.Println("\n========== 写入秘密示例 ==========")
	writeData := map[string]interface{}{
		"example_key": "example_value",
		"created_by":  "mcube_example",
		"created_at":  "2026-02-20",
	}

	_, err = iocvault.WriteSecret(ctx, "myapp/example", writeData)
	if err != nil {
		fmt.Printf("Error writing secret: %v\n", err)
	} else {
		fmt.Println("✓ 秘密写入成功: myapp/example")
	}

	// 读取刚写入的秘密验证
	verify, err := iocvault.ReadSecret(ctx, "myapp/example")
	if err == nil {
		fmt.Printf("✓ 验证读取成功: %+v\n", verify.Data.Data)
	}
}

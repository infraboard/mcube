package main

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/vault"

	// 引入 vault 模块，使其自动注册
	_ "github.com/infraboard/mcube/v2/ioc/config/vault"
)

func main() {
	// 初始化 IOC 容器，会自动加载配置并初始化 vault 客户端
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}

	// 获取 Vault 客户端
	client := vault.Client()
	fmt.Printf("Vault client initialized: %v\n", client != nil)

	// 读取秘密示例
	// 假设在 Vault 中有一个路径为 secret/data/myapp/config 的秘密
	ctx := context.Background()
	resp, err := client.Secrets.KvV2Read(ctx, "myapp/config")
	if err != nil {
		fmt.Printf("Error reading secret: %v\n", err)
		return
	}

	fmt.Printf("Secret data: %+v\n", resp.Data.Data)
	fmt.Printf("Secret metadata: %+v\n", resp.Data.Metadata)

	// 访问具体的字段
	fmt.Println("\nSecret fields:")
	for k, v := range resp.Data.Data {
		fmt.Printf("  %s: %v\n", k, v)
	}
}

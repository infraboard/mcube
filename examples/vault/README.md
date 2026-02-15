# Vault 配置模块示例

本示例展示如何使用 mcube 的 Vault 配置模块。

## 前置准备

### 1. 启动 Vault 服务

使用 Docker 启动一个开发模式的 Vault 服务器：

```bash
docker run -d --name=vault \
  --cap-add=IPC_LOCK \
  -p 8200:8200 \
  -e 'VAULT_DEV_ROOT_TOKEN_ID=myroot' \
  -e 'VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200' \
  hashicorp/vault:latest
```

### 2. 创建测试秘密

```bash
export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_TOKEN='myroot'

# 写入测试秘密
vault kv put secret/myapp/config \
  database_url="postgres://localhost:5432/mydb" \
  api_key="my-secret-api-key" \
  debug_mode="true"

# 验证秘密
vault kv get secret/myapp/config
```

## 配置说明

编辑 `etc/application.toml` 文件配置 Vault 连接：

```toml
[vault]
address = "http://127.0.0.1:8200"
auth_method = "token"
token = "myroot"
auto_renew = true
```

### 支持的认证方式

#### 1. Token 认证（推荐用于开发/测试）

```toml
auth_method = "token"
token = "hvs.xxxxxxxxxxxxx"
```

#### 2. AppRole 认证（推荐用于生产环境）

```toml
auth_method = "approle"
role_id = "your-role-id"
secret_id = "your-secret-id"
```

设置 AppRole：

```bash
# 启用 AppRole 认证
vault auth enable approle

# 创建角色
vault write auth/approle/role/myapp \
  secret_id_ttl=10m \
  token_num_uses=10 \
  token_ttl=20m \
  token_max_ttl=30m \
  secret_id_num_uses=40 \
  policies="myapp-policy"

# 获取 Role ID
vault read auth/approle/role/myapp/role-id

# 生成 Secret ID
vault write -f auth/approle/role/myapp/secret-id
```

#### 3. Kubernetes 认证（推荐用于 K8s 环境）

```toml
auth_method = "kubernetes"
k8s_role = "myapp"
k8s_token_path = "/var/run/secrets/kubernetes.io/serviceaccount/token"
```

## 运行示例

```bash
cd examples/vault
go run main.go
```

## 使用说明

### 基本使用

```go
import (
    "context"
    "github.com/infraboard/mcube/v2/ioc/config/vault"
    _ "github.com/infraboard/mcube/v2/ioc/config/vault" // 引入模块注册
)

// 获取 Vault 客户端
client := vault.Client()

// 读取 KV v2 秘密
resp, err := client.Secrets.KvV2Read(context.Background(), "myapp/config")
if err != nil {
    panic(err)
}

// 访问秘密数据
for key, value := range resp.Data.Data {
    fmt.Printf("%s: %v\n", key, value)
}

// 访问特定字段
if databaseURL, ok := resp.Data.Data["database_url"].(string); ok {
    fmt.Println("Database URL:", databaseURL)
}
```

### 写入秘密

```go
_, err := client.Secrets.KvV2Write(
    context.Background(),
    "myapp/config",
    schema.KvV2WriteRequest{
        Data: map[string]interface{}{
            "database_url": "postgres://localhost:5432/mydb",
            "api_key":      "new-api-key",
        },
    },
)
```

### 其他操作

vault-client-go 支持 Vault 的所有功能：

- KV 秘密引擎（v1 和 v2）
- Transit 加密
- PKI 证书管理
- Database 动态凭证
- Token 管理
- 策略管理

详细 API 文档：https://pkg.go.dev/github.com/hashicorp/vault-client-go

## TLS 配置

如果 Vault 启用了 HTTPS：

```toml
[vault]
address = "https://vault.example.com:8200"
auth_method = "token"
token = "hvs.xxxxxxxxxxxxx"
tls_enable = true
tls_ca = "/path/to/ca.crt"
# 可选：客户端证书认证
tls_cert = "/path/to/client.crt"
tls_key = "/path/to/client.key"
# 仅测试环境使用
# tls_skip_verify = true
```

## 企业版 Namespace

如果使用 Vault 企业版的 namespace 功能：

```toml
[vault]
namespace = "admin"
```

或通过环境变量：

```bash
export VAULT_NAMESPACE=admin
```

## 注意事项

1. **生产环境**：不要使用 token 认证，推荐 AppRole 或 Kubernetes 认证
2. **Token 续期**：设置 `auto_renew = true` 自动续期 token
3. **TLS**：生产环境必须启用 TLS
4. **秘密版本**：KV v2 支持秘密版本控制，推荐使用
5. **最小权限**：为应用分配最小必要的 Vault 策略权限

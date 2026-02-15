# Vault 配置模块

Vault 配置模块为 mcube 应用提供了与 HashiCorp Vault 集成的能力，支持从 Vault 读取秘密、管理凭证和加密数据。

## 特性

- ✅ **多种认证方式**: Token、AppRole、Kubernetes 认证
- ✅ **自动 Token 续期**: 支持自动续期，避免 token 过期
- ✅ **TLS 支持**: 支持 TLS 客户端证书认证
- ✅ **Namespace支持**: 支持企业版 Namespace 功能
- ✅ **灵活的客户端**: 直接返回 vault-client-go 客户端，支持所有 Vault 功能

## 快速开始

### 1. 引入模块

```go
import (
    _ "github.com/infraboard/mcube/v2/ioc/config/vault"
)
```

### 2. 配置

在 `application.toml` 中添加 Vault 配置：

```toml
[vault]
address = "http://127.0.0.1:8200"
auth_method = "token"
token = "myroot"
auto_renew = true
```

### 3. 使用

```go
import "github.com/infraboard/mcube/v2/ioc/config/vault"

// 获取 Vault 客户端
client := vault.Client()

// 读取秘密
ctx := context.Background()
resp, err := client.Secrets.KvV2Read(ctx, "myapp/config")
if err != nil {
    panic(err)
}

// 访问秘密数据
for key, value := range resp.Data.Data {
    fmt.Printf("%s: %v\n", key, value)
}
```

## 配置说明

### 基础配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `address` | string | http://127.0.0.1:8200 | Vault 服务器地址 |
| `auth_method` | string | token | 认证方式: token, approle, kubernetes |
| `auto_renew` | bool | true | 是否自动续期 token |
| `namespace` | string | - | Vault Namespace（企业版） |

### Token 认证

最简单的认证方式，适合开发和测试环境：

```toml
[vault]
auth_method = "token"
token = "hvs.xxxxxxxxxxxxx"
```

### AppRole 认证

推荐用于生产环境的应用程序：

```toml
[vault]
auth_method = "approle"
role_id = "your-role-id"
secret_id = "your-secret-id"
```

**设置 AppRole**:

```bash
# 启用 AppRole 认证
vault auth enable approle

# 创建角色和策略
vault write auth/approle/role/myapp \
  secret_id_ttl=10m \
  token_ttl=20m \
  token_max_ttl=30m \
  policies="myapp-policy"

# 获取凭证
vault read auth/approle/role/myapp/role-id
vault write -f auth/approle/role/myapp/secret-id
```

### Kubernetes 认证

适合在 Kubernetes 环境中运行的应用：

```toml
[vault]
auth_method = "kubernetes"
k8s_role = "myapp"
k8s_token_path = "/var/run/secrets/kubernetes.io/serviceaccount/token"
```

**配置 Kubernetes 认证**:

```bash
# 启用 Kubernetes 认证
vault auth enable kubernetes

# 配置 Kubernetes 认证后端
vault write auth/kubernetes/config \
  kubernetes_host="https://kubernetes.default.svc:443" \
  kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt \
  token_reviewer_jwt=@/var/run/secrets/kubernetes.io/serviceaccount/token

# 创建角色绑定到 ServiceAccount
vault write auth/kubernetes/role/myapp \
  bound_service_account_names=myapp \
  bound_service_account_namespaces=default \
  policies=myapp-policy \
  ttl=24h
```

### TLS 配置

生产环境强烈建议启用 TLS：

```toml
[vault]
address = "https://vault.example.com:8200"
tls_enable = true
tls_ca = "/path/to/ca.crt"
# 可选：客户端证书认证
tls_cert = "/path/to/client.crt"
tls_key = "/path/to/client.key"
# 仅在开发环境使用
# tls_skip_verify = true
```

## API 使用示例

### 读取秘密

```go
resp, err := client.Secrets.KvV2Read(ctx, "path/to/secret")
if err != nil {
    return err
}

// 访问秘密数据
password := resp.Data.Data["password"].(string)
```

### 写入秘密

```go
import "github.com/hashicorp/vault-client-go/schema"

_, err := client.Secrets.KvV2Write(
    ctx,
    "path/to/secret",
    schema.KvV2WriteRequest{
        Data: map[string]interface{}{
            "username": "admin",
            "password": "secret123",
        },
    },
)
```

### 列出秘密路径

```go
resp, err := client.Secrets.KvV2List(ctx, "path/to/")
if err != nil {
    return err
}

for _, key := range resp.Data.Keys {
    fmt.Println(key)
}
```

### 删除秘密

```go
// 删除最新版本
_, err := client.Secrets.KvV2Delete(ctx, "path/to/secret")

// 删除特定版本
_, err := client.Secrets.KvV2DeleteVersions(
    ctx,
    "path/to/secret",
    schema.KvV2DeleteVersionsRequest{
        Versions: []int64{1, 2},
    },
)
```

### 使用 Transit 加密

```go
import "github.com/hashicorp/vault-client-go/schema"

// 加密数据
encResp, err := client.Secrets.TransitEncrypt(
    ctx,
    "my-key",
    schema.TransitEncryptRequest{
        Plaintext: base64.StdEncoding.EncodeToString([]byte("secret data")),
    },
)

ciphertext := encResp.Data.Ciphertext

// 解密数据
decResp, err := client.Secrets.TransitDecrypt(
    ctx,
    "my-key",
    schema.TransitDecryptRequest{
        Ciphertext: ciphertext,
    },
)

plaintext, _ := base64.StdEncoding.DecodeString(decResp.Data.Plaintext)
```

### 动态数据库凭证

```go
// 获取数据库凭证
resp, err := client.Secrets.DatabaseGenerateCredentials(ctx, "my-role")
if err != nil {
    return err
}

username := resp.Data.Data["username"].(string)
password := resp.Data.Data["password"].(string)

// 使用凭证连接数据库...

// 凭证会在 lease_duration 后自动过期
```

## 环境变量

支持通过环境变量覆盖配置：

- `VAULT_ADDR`: Vault 服务器地址
- `VAULT_TOKEN`: Token 认证的 token
- `VAULT_NAMESPACE`: Vault namespace（企业版）
- `VAULT_CACERT`: CA 证书路径
- `VAULT_CLIENT_CERT`: 客户端证书路径
- `VAULT_CLIENT_KEY`: 客户端私钥路径
- `VAULT_SKIP_VERIFY`: 跳过 TLS 验证（值为 "true"）

## 最佳实践

1. **生产环境认证**: 使用 AppRole 或 Kubernetes 认证，避免使用 token 认证
2. **最小权限原则**: 为每个应用创建独立的策略，只授予必要的权限
3. **启用 TLS**: 生产环境必须启用 TLS 加密通信
4. **自动续期**: 开启 `auto_renew` 避免 token 过期
5. **动态秘密**: 优先使用动态秘密（如数据库凭证），而不是静态秘密
6. **秘密轮换**: 定期轮换秘密，使用 Vault 的版本控制功能
7. **审计日志**: 启用 Vault 审计日志，追踪秘密访问
8. **密钥管理**: 使用 Transit 引擎进行加密操作，而不是在应用中硬编码密钥

## 故障处理

### Token 过期

启用自动续期：

```toml
[vault]
auto_renew = true
```

模块会自动在 token 生命周期的 80% 时尝试续期。

### 连接失败

检查网络连接和 Vault 服务状态：

```bash
vault status
```

验证 TLS 配置是否正确。

### 权限不足

检查策略配置：

```bash
# 查看当前 token 的策略
vault token lookup

# 查看策略内容
vault policy read myapp-policy
```

确保策略包含必要的权限：

```hcl
# 读取 KV 秘密
path "secret/data/myapp/*" {
  capabilities = ["read", "list"]
}

# 写入 KV 秘密
path "secret/data/myapp/*" {
  capabilities = ["create", "update"]
}
```

## 相关资源

- [HashiCorp Vault 官方文档](https://www.vaultproject.io/docs)
- [vault-client-go API 文档](https://pkg.go.dev/github.com/hashicorp/vault-client-go)
- [Vault 最佳实践](https://learn.hashicorp.com/tutorials/vault/production-hardening)
- [示例代码](https://github.com/infraboard/mcube/tree/master/examples/vault)

## License

MIT

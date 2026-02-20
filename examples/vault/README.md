# Vault 配置模块示例

本示例展示如何使用 mcube 的 Vault 配置模块。

## ⚠️ 重要：路径使用说明

在使用 vault-client-go 时，需要正确指定挂载点：

```go
// ✅ 正确方式
resp, err := client.Secrets.KvV2Read(
    ctx,
    "myapp/config",                  // 秘密路径（不含挂载点）
    vault.WithMountPath("secret"),   // KV v2 引擎挂载点
)

// ❌ 错误方式
resp, err := client.Secrets.KvV2Read(ctx, "secret/myapp/config")
// 错误: 404 Not Found: no handler for route "kv-v2/data/secret/myapp/config"
```

**为什么？**
- Vault CLI: `vault kv get secret/myapp/config` ← 自动处理
- Vault API: 需要分离挂载点和秘密路径
  - 挂载点: `secret` (开发模式默认)
  - 秘密路径: `myapp/config`
  - 完整 API 路径: `secret/data/myapp/config` (自动添加 `/data/`)

详细说明请参考: [PATH_GUIDE.md](../../ioc/config/vault/PATH_GUIDE.md)

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
# 格式: vault kv put <挂载点>/<秘密路径> key=value
#                     ────┬───  ─────┬─────
#                       secret    myapp/config
vault kv put secret/myapp/config \
  database_url="postgres://localhost:5432/mydb" \
  api_key="my-secret-api-key" \
  debug_mode="true"

# 验证秘密
vault kv get secret/myapp/config
```

**路径说明**:
- `secret` = 挂载点（KV v2 引擎的挂载位置）
- `myapp/config` = 秘密路径（实际存储的路径）
- 完整 API 路径 = `secret/data/myapp/config`（自动添加 `/data/`）

## 配置说明

编辑 `etc/application.toml` 文件配置 Vault 连接：

```toml
[vault]
address = "http://127.0.0.1:8200"
auth_method = "token"
token = "myroot"
auto_renew = true

# 挂载点配置（与 Vault CLI 保持一致）
kv_mount_path = "secret"  # KV v2 引擎挂载点，对应 CLI 中的 secret/
```

**配置与 CLI 对应关系**:

| CLI 命令 | 配置项 | 代码调用 |
|----------|--------|----------|
| `vault kv get secret/myapp/config` | `kv_mount_path = "secret"` | `vault.ReadSecret(ctx, "myapp/config")` |
| `vault kv get kv/prod/db` | `kv_mount_path = "kv"` | `vault.ReadSecret(ctx, "prod/db")` |

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
## 使用说明

### 推荐方式：使用辅助方法

辅助方法自动使用配置的挂载点，与 Vault CLI 保持一致的使用体验：

```go
import (
    "context"
    "github.com/infraboard/mcube/v2/ioc/config/vault"
    _ "github.com/infraboard/mcube/v2/ioc/config/vault" // 引入模块注册
)

ctx := context.Background()

// 读取秘密
// 对应 CLI: vault kv get secret/myapp/config
resp, err := vault.ReadSecret(ctx, "myapp/config")
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

// 写入秘密
// 对应 CLI: vault kv put secret/myapp/config database_url="..." api_key="..."
data := map[string]interface{}{
    "database_url": "postgres://localhost:5432/mydb",
    "api_key":      "new-api-key",
}
_, err = vault.WriteSecret(ctx, "myapp/config", data)
```

**为什么推荐辅助方法？**
- ✅ 自动使用配置的挂载点，无需每次指定
- ✅ 与 Vault CLI 用法保持一致
- ✅ 代码更简洁，减少出错

### 高级用法：手动指定挂载点

如需更精细的控制或使用多个不同挂载点，可直接使用 vault-client-go 客户端：

```go
import (
    "context"
    "github.com/hashicorp/vault-client-go"
    "github.com/infraboard/mcube/v2/ioc/config/vault"
    _ "github.com/infraboard/mcube/v2/ioc/config/vault"
)

// 获取 Vault 客户端
client := vault.Client()

// 读取 KV v2 秘密 - 手动指定挂载点
// 对应 CLI: vault kv get secret/myapp/config
resp, err := client.Secrets.KvV2Read(
    context.Background(),
    "myapp/config",                  // 秘密路径（不包含挂载点）
    vault.WithMountPath("secret"),   // 挂载点
)
if err != nil {
    panic(err)
}

// 访问秘密数据
for key, value := range resp.Data.Data {
    fmt.Printf("%s: %v\n", key, value)
}

// 如果使用不同的挂载点
// 对应 CLI: vault kv get kv/prod/database
prodResp, err := client.Secrets.KvV2Read(
    context.Background(),
    "prod/database",              // 秘密路径
    vault.WithMountPath("kv"),    // 不同的挂载点
)
```

### 路径对应关系速查表

| Vault CLI | vault-client-go (手动) | vault 辅助方法 (推荐) |
|----------|----------------------|---------------------|
| `vault kv get secret/myapp/config` | `KvV2Read(ctx, "myapp/config", WithMountPath("secret"))` | `vault.ReadSecret(ctx, "myapp/config")` |
| `vault kv put secret/myapp/config` | `KvV2Write(ctx, "myapp/config", data, WithMountPath("secret"))` | `vault.WriteSecret(ctx, "myapp/config", data)` |
| `vault kv list secret/myapp` | `KvV2List(ctx, "myapp", WithMountPath("secret"))` | `vault.ListSecrets(ctx, "myapp")` |
| `vault kv delete secret/myapp/config` | `KvV2Delete(ctx, "myapp/config", WithMountPath("secret"))` | `vault.DeleteSecret(ctx, "myapp/config")` |

### 写入秘密

```go
import "github.com/hashicorp/vault-client-go/schema"

_, err := client.Secrets.KvV2Write(
    context.Background(),
    "myapp/config",
    schema.KvV2WriteRequest{
        Data: map[string]interface{}{
            "database_url": "postgres://localhost:5432/mydb",
            "api_key":      "new-api-key",
        },
    },
    vault.WithMountPath("secret"),  // 指定挂载点
)
```

## 完整示例

参考 [main.go](main.go) 查看完整的示例代码，包括：
- 辅助方法的简洁用法
- 手动指定挂载点的灵活用法
- 错误处理和数据访问

## 运行示例

```bash
cd examples/vault
go run main.go
```

## 其他功能

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

## 挂载点（Mount Path）概念说明

### Mount Path vs Namespace

**Mount Path（挂载点）** 和 **Namespace（命名空间）** 是两个不同概念：

| 概念 | 层级 | 支持版本 | 用途 |
|------|------|----------|------|
| **Mount Path** | 秘密引擎级别 | 所有版本 | 组织不同类型/用途的秘密 |
| **Namespace** | 租户级别 | 仅企业版 | 多租户逻辑隔离 |

### 生产环境挂载点策略

**`secret` 在生产环境中完全可以使用**，但可根据需求自定义：

#### 方案 1: 使用默认 `secret`（适合小型项目）
```toml
[vault]
kv_mount_path = "secret"
```
✅ 与 Vault 默认保持一致  
⚠️ 需通过 Vault Policy 精细控制权限

#### 方案 2: 按环境划分（推荐生产环境）
```bash
# 在 Vault 服务器上创建不同挂载点
vault secrets enable -path=dev-kv kv-v2
vault secrets enable -path=staging-kv kv-v2
vault secrets enable -path=prod-kv kv-v2
```

```toml
# 应用配置 - 开发环境
[vault]
kv_mount_path = "dev-kv"

# 应用配置 - 生产环境
[vault]
kv_mount_path = "prod-kv"
```
✅ 环境隔离清晰，降低误操作风险

#### 方案 3: 按应用/团队划分（适合多团队）
```bash
vault secrets enable -path=team-a-secrets kv-v2
vault secrets enable -path=team-b-secrets kv-v2
vault secrets enable -path=app-x-kv kv-v2
```

### 查看和管理挂载点

```bash
# 查看所有挂载点
vault secrets list

# 启用新的 KV v2 挂载点
vault secrets enable -path=my-secrets kv-v2

# 禁用挂载点（谨慎操作！）
vault secrets disable my-secrets/
```

## 注意事项

1. **生产环境认证**：不要使用 token 认证，推荐 AppRole 或 Kubernetes 认证
2. **Token 续期**：设置 `auto_renew = true` 自动续期 token
3. **TLS 安全**：生产环境必须启用 TLS
4. **秘密版本**：KV v2 支持秘密版本控制，推荐使用
5. **最小权限**：为应用分配最小必要的 Vault 策略权限
6. **挂载点规划**：生产环境建议按环境/应用划分挂载点，避免混用
7. **命名规范**：挂载点名称建议使用小写字母和连字符，如 `prod-kv`, `team-secrets`

# Vault 配置模块

Vault 配置模块为 mcube 应用提供了与 HashiCorp Vault 集成的能力，支持从 Vault 读取秘密、管理凭证和加密数据。

## 什么是 Vault？

[HashiCorp Vault](https://www.vaultproject.io/) 是一个统一的秘密管理平台，用于安全地存储和访问敏感数据。

### 为什么需要 Vault？

在传统开发中，我们经常遇到这些问题：

❌ **配置文件中硬编码密码**
```yaml
database:
  password: "mypassword123"  # 不安全！
```

❌ **环境变量泄露风险**
```bash
export DB_PASSWORD="secret"  # 进程列表可见
```

❌ **静态凭证难以轮换**
```go
// 密码写死在代码中，更换困难
conn := db.Connect("postgres://admin:oldpassword@localhost")
```

✅ **使用 Vault 后**
```go
// 从 Vault 动态获取最新凭证
resp, _ := client.Secrets.KvV2Read(ctx, "database/config")
password := resp.Data.Data["password"].(string)
```

### Vault 能做什么？

| 功能 | 场景 | 价值 |
|------|------|------|
| **秘密存储** | 存储 API 密钥、数据库密码、证书等敏感信息 | 集中管理，加密存储，访问控制 |
| **动态凭证** | 为数据库、云平台生成临时凭证 | 凭证自动过期，减少泄露风险 |
| **数据加密** | 加密敏感数据（信用卡号、身份证号） | 无需管理加密密钥，Vault 负责加密逻辑 |
| **证书管理** | 自动签发和续期 TLS 证书 | 简化 PKI 管理，自动化证书生命周期 |

### 核心概念

**1. 秘密（Secret）**
```
秘密是指任何需要严格控制访问的敏感数据，例如：
- API 密钥：第三方服务的认证令牌
- 数据库密码：MySQL、PostgreSQL 连接密码
- 加密密钥：用于数据加密的密钥材料
- TLS 证书：HTTPS 通信所需的证书和私钥
```

**2. 认证（Authentication）**
```
应用程序需要向 Vault 证明身份后才能访问秘密：
- Token：最简单，适合开发测试
- AppRole：基于角色的认证，适合服务器应用
- Kubernetes：利用 K8s ServiceAccount，适合容器环境
```

**3. 策略（Policy）**
```
定义谁可以访问哪些秘密，实现细粒度权限控制：

path "secret/data/myapp/*" {
  capabilities = ["read"]    # 只能读取
}

path "database/creds/readonly" {
  capabilities = ["read"]    # 只能读取数据库凭证
}
```

## 特性

- ✅ **多种认证方式**: Token、AppRole、Kubernetes 认证
- ✅ **自动 Token 续期**: 支持自动续期，避免 token 过期
- ✅ **TLS 支持**: 支持 TLS 客户端证书认证
- ✅ **Namespace支持**: 支持企业版 Namespace 功能
- ✅ **灵活的客户端**: 直接返回 vault-client-go 客户端，支持所有 Vault 功能

## 典型使用场景

### 场景 1：统一管理应用配置

**问题**: 多个微服务共享数据库密码，密码散落在各个配置文件中，更换密码需要重启所有服务。

**方案**: 将密码存储在 Vault，应用启动时从 Vault 读取最新密码。

```go
// 应用启动时读取数据库配置
resp, err := client.Secrets.KvV2Read(ctx, "myapp/database")
dbConfig := map[string]string{
    "host":     resp.Data.Data["host"].(string),
    "port":     resp.Data.Data["port"].(string),
    "username": resp.Data.Data["username"].(string),
    "password": resp.Data.Data["password"].(string),
}

// 连接数据库
db := connectDatabase(dbConfig)
```

**收益**:
- ✅ 集中管理：更新密码只需更新 Vault，无需修改每个服务的配置
- ✅ 访问审计：可追踪哪个服务在什么时候访问了密码
- ✅ 版本控制：Vault KV v2 支持秘密版本，可回滚到旧密码

### 场景 2：保护敏感数据

**问题**: 需要存储用户的信用卡号、身份证号等敏感信息，但不想自己管理加密密钥。

**方案**: 使用 Vault 的 Transit 加密引擎，在保存前加密，读取时解密。

```go
// 加密敏感数据
plaintext := base64.StdEncoding.EncodeToString([]byte("6222021234567890"))
encResp, _ := client.Secrets.TransitEncrypt(ctx, "credit-cards", 
    schema.TransitEncryptRequest{Plaintext: plaintext})

// 存储到数据库
db.Save(encResp.Data.Ciphertext)  // 存储: vault:v1:abcd1234...

// 从数据库读取并解密
ciphertext := db.Query(...)
decResp, _ := client.Secrets.TransitDecrypt(ctx, "credit-cards",
    schema.TransitDecryptRequest{Ciphertext: ciphertext})
original, _ := base64.StdEncoding.DecodeString(decResp.Data.Plaintext)
```

**收益**:
- ✅ 密钥集中管理：加密密钥存储在 Vault，应用无需管理密钥
- ✅ 密钥轮换：Vault 支持密钥轮换，旧数据仍可解密
- ✅ 合规要求：符合 PCI-DSS、GDPR 等数据保护规范

### 场景 3：动态数据库凭证

**问题**: 所有应用共享一个 MySQL 账号，一旦泄露影响范围大，且无法追踪是谁泄露的。

**方案**: 使用 Vault 动态生成临时数据库账号，每个应用实例使用独立凭证，凭证自动过期。

```go
// 请求临时数据库凭证（有效期 1 小时）
resp, _ := client.Secrets.DatabaseGenerateCredentials(ctx, "readonly-role")
username := resp.Data.Data["username"].(string)  // vault-token-abc123
password := resp.Data.Data["password"].(string)  // 随机生成

// 使用临时凭证连接数据库
dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/mydb", username, password)
db, _ := sql.Open("mysql", dsn)

// 1 小时后凭证自动失效，数据库账号被删除
```

**收益**:
- ✅ 最小权限：每个应用只获得所需的最小权限（如只读访问）
- ✅ 自动过期：凭证有生命周期，即使泄露影响也可控
- ✅ 审计追踪：可追踪每个凭证是谁在何时申请的

### 场景 4：自动化证书管理

**问题**: 微服务网格需要大量 TLS 证书，手动签发和续期证书繁琐且容易出错。

**方案**: 使用 Vault PKI 引擎自动签发和续期证书。

```go
// 申请证书（有效期 30 天）
resp, _ := client.Secrets.PkiIssue(ctx, "myapp-role",
    schema.PkiIssueRequest{
        CommonName: "service1.example.com",
        Ttl:        "720h",  // 30 days
    })

certificate := resp.Data.Certificate
privateKey := resp.Data.PrivateKey

// 配置 TLS 服务器
cert, _ := tls.X509KeyPair([]byte(certificate), []byte(privateKey))
```

**收益**:
- ✅ 自动化：无需手动生成 CSR、签名、分发证书
- ✅ 短周期：证书短期有效（如 30 天），减少泄露风险
- ✅ 零信任：结合服务网格实现服务间零信任通信

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

### 前置准备

为了方便功能验证，本示例将使用 Docker 启动一个 Vault 服务器。

#### 1. 启动 Vault 服务

使用 Docker 启动一个开发模式的 Vault 服务器：

```bash
docker run -itd --name=vault \
  --cap-add=IPC_LOCK \
  -p 8200:8200 \
  -e 'VAULT_DEV_ROOT_TOKEN_ID=myroot' \
  -e 'VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200' \
  hashicorp/vault:latest
```

#### 2. 创建测试秘密

在容器内执行命令（本地无需安装 Vault CLI）：

```bash
# 写入测试秘密
docker exec \
  -e VAULT_ADDR='http://127.0.0.1:8200' \
  -e VAULT_TOKEN='myroot' \
  vault vault kv put secret/myapp/config \
    database_url="postgres://localhost:5432/mydb" \
    api_key="my-secret-api-key" \
    debug_mode="true"

# 验证秘密
docker exec \
  -e VAULT_ADDR='http://127.0.0.1:8200' \
  -e VAULT_TOKEN='myroot' \
  vault vault kv get secret/myapp/config
# 输出
====== Secret Path ======
secret/data/myapp/config

======= Metadata =======
Key                Value
---                -----
created_time       2026-02-20T01:46:00.710838087Z
custom_metadata    <nil>
deletion_time      n/a
destroyed          false
version            1

======== Data ========
Key             Value
---             -----
api_key         my-secret-api-key
database_url    postgres://localhost:5432/mydb
debug_mode      true
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

package vault

import (
	"context"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
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

// KvMountPath 返回 KV v2 引擎的挂载点
func KvMountPath() string {
	return Get().KvMountPath
}

// TransitMountPath 返回 Transit 引擎的挂载点
func TransitMountPath() string {
	return Get().TransitMountPath
}

// DatabaseMountPath 返回 Database 引擎的挂载点
func DatabaseMountPath() string {
	return Get().DatabaseMountPath
}

// PkiMountPath 返回 PKI 引擎的挂载点
func PkiMountPath() string {
	return Get().PkiMountPath
}

// ==================== KV v2 辅助方法 ====================

// ReadSecret 读取 KV v2 秘密（自动使用配置的挂载点）
func ReadSecret(ctx context.Context, path string, options ...vault.RequestOption) (*vault.Response[schema.KvV2ReadResponse], error) {
	opts := append([]vault.RequestOption{vault.WithMountPath(KvMountPath())}, options...)
	return Client().Secrets.KvV2Read(ctx, path, opts...)
}

// WriteSecret 写入 KV v2 秘密（自动使用配置的挂载点）
func WriteSecret(ctx context.Context, path string, data map[string]interface{}, options ...vault.RequestOption) (*vault.Response[schema.KvV2WriteResponse], error) {
	opts := append([]vault.RequestOption{vault.WithMountPath(KvMountPath())}, options...)
	return Client().Secrets.KvV2Write(ctx, path, schema.KvV2WriteRequest{Data: data}, opts...)
}

// DeleteSecret 删除 KV v2 秘密最新版本（自动使用配置的挂载点）
func DeleteSecret(ctx context.Context, path string, options ...vault.RequestOption) error {
	opts := append([]vault.RequestOption{vault.WithMountPath(KvMountPath())}, options...)
	_, err := Client().Secrets.KvV2Delete(ctx, path, opts...)
	return err
}

// ListSecrets 列出 KV v2 秘密路径（自动使用配置的挂载点）
func ListSecrets(ctx context.Context, path string, options ...vault.RequestOption) (*vault.Response[schema.StandardListResponse], error) {
	opts := append([]vault.RequestOption{vault.WithMountPath(KvMountPath())}, options...)
	return Client().Secrets.KvV2List(ctx, path, opts...)
}

// ==================== Transit 辅助方法 ====================

// Encrypt 使用 Transit 引擎加密数据（自动使用配置的挂载点）
func Encrypt(ctx context.Context, keyName string, plaintext string, options ...vault.RequestOption) (*vault.Response[map[string]interface{}], error) {
	opts := append([]vault.RequestOption{vault.WithMountPath(TransitMountPath())}, options...)
	return Client().Secrets.TransitEncrypt(ctx, keyName, schema.TransitEncryptRequest{Plaintext: plaintext}, opts...)
}

// Decrypt 使用 Transit 引擎解密数据（自动使用配置的挂载点）
func Decrypt(ctx context.Context, keyName string, ciphertext string, options ...vault.RequestOption) (*vault.Response[map[string]interface{}], error) {
	opts := append([]vault.RequestOption{vault.WithMountPath(TransitMountPath())}, options...)
	return Client().Secrets.TransitDecrypt(ctx, keyName, schema.TransitDecryptRequest{Ciphertext: ciphertext}, opts...)
}

// ==================== Database 辅助方法 ====================

// GenerateDatabaseCredentials 生成动态数据库凭证（自动使用配置的挂载点）
func GenerateDatabaseCredentials(ctx context.Context, roleName string, options ...vault.RequestOption) (*vault.Response[map[string]interface{}], error) {
	opts := append([]vault.RequestOption{vault.WithMountPath(DatabaseMountPath())}, options...)
	return Client().Secrets.DatabaseGenerateCredentials(ctx, roleName, opts...)
}

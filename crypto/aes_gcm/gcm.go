package aesgcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// 便捷函数 - 直接使用字符串密钥

// NewAESGCMFromString 从字符串创建AESGCM实例
// 适用于从配置文件或环境变量中读取密钥的场景
// 注意：字符串形式的密钥安全性较低，建议使用字节数组密钥
func NewAESGCMFromString(key string, keySize KeySize) (*AESGCM, error) {
	if len(key) != int(keySize) {
		return nil, fmt.Errorf("key string length must be %d characters", keySize)
	}
	return NewAESGCM([]byte(key))
}

// 便捷函数 - 使用默认密钥长度(256位)

// NewAES256GCM 创建使用256位密钥的AESGCM实例
// 推荐用于大多数生产环境，提供最高级别的安全性
func NewAES256GCM(key []byte) (*AESGCM, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes for AES-256")
	}
	return NewAESGCM(key)
}

// NewAESGCM 创建一个新的AESGCM实例
// 适用于需要灵活指定任意有效密钥长度的场景
func NewAESGCM(key []byte) (*AESGCM, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, ErrInvalidKeySize
	}
	return &AESGCM{key: key}, nil
}

// AESGCM 封装AES-GCM加密解密功能
type AESGCM struct {
	key []byte
}

// Encrypt 加密明文数据
// 基础加密方法，适用于大多数简单的加密场景，不涉及额外的认证数据
// 返回格式: nonce + ciphertext + tag
func (a *AESGCM) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt 解密密文数据
// 基础解密方法，对应Encrypt方法的输出
// 输入格式: nonce + ciphertext + tag
func (a *AESGCM) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// 检查密文长度
	if len(ciphertext) < gcm.NonceSize() {
		return nil, ErrInvalidCiphertext
	}

	// 提取nonce和实际密文
	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}

// EncryptWithAdditionalData 使用附加数据进行加密
// 适用于需要同时加密和认证额外数据的场景，如：
// - 加密消息头信息
// - 保护关联的元数据
// - 需要额外上下文认证的场景
// 附加数据本身不会被加密，但会参与认证标签的计算
func (a *AESGCM) EncryptWithAdditionalData(plaintext, additionalData []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, additionalData)
	return ciphertext, nil
}

// DecryptWithAdditionalData 使用附加数据进行解密
// 对应EncryptWithAdditionalData方法的输出
// 必须提供与加密时相同的附加数据，否则认证会失败
func (a *AESGCM) DecryptWithAdditionalData(ciphertext, additionalData []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, ErrInvalidCiphertext
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, additionalData)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}

// EncryptToString 加密并返回base64编码的字符串
// 适用于需要将加密结果存储为字符串的场景，如：
// - 数据库存储
// - JSON/XML序列化
// - URL参数传递
// 返回的字符串可以直接存储或传输
func (a *AESGCM) EncryptToString(plaintext string) (string, error) {
	ciphertext, err := a.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptFromString 从base64字符串解密
// 对应EncryptToString方法的输出
// 适用于从字符串格式恢复加密数据的场景
func (a *AESGCM) DecryptFromString(encodedCiphertext string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	plaintext, err := a.Decrypt(ciphertext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

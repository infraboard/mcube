package aesgcm

import (
	"crypto/rand"
	"fmt"
	"io"
)

// KeySize 定义密钥大小的枚举类型
type KeySize int

const (
	AES128 KeySize = 16 // AES-128: 安全性足够，性能最佳
	AES192 KeySize = 24 // AES-192: 中等安全性，不常用
	AES256 KeySize = 32 // AES-256: 最高安全性，推荐用于敏感数据
)

// String 返回密钥大小的字符串表示
func (k KeySize) String() string {
	switch k {
	case AES128:
		return "AES-128"
	case AES192:
		return "AES-192"
	case AES256:
		return "AES-256"
	default:
		return fmt.Sprintf("Unknown(%d)", int(k))
	}
}

// Valid 检查密钥大小是否有效
func (k KeySize) Valid() bool {
	return k == AES128 || k == AES192 || k == AES256
}

// Generate256Key 生成256位随机密钥
// 推荐用于生成高安全性的密钥
func Generate256Key() ([]byte, error) {
	return GenerateKey(AES256)
}

// Generate128Key 生成128位随机密钥
// 适用于性能敏感但安全性要求稍低的场景
func Generate128Key() ([]byte, error) {
	return GenerateKey(AES128)
}

// Generate192Key 生成192位随机密钥
// 适用于中等安全性要求的场景（不常用）
func Generate192Key() ([]byte, error) {
	return GenerateKey(AES192)
}

// GenerateKey 生成指定长度的随机密钥
// 适用于密钥生成场景，使用枚举类型确保密钥长度正确
func GenerateKey(keySize KeySize) ([]byte, error) {
	if !keySize.Valid() {
		return nil, ErrInvalidKeySize
	}

	key := make([]byte, keySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}
	return key, nil
}

package aesgcm

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
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

// NewAESGCMFromBase64 从base64字符串创建AESGCM实例
// 适用于从配置文件读取base64编码密钥的场景
func NewAESGCMFromBase64(keyBase64 string) (*AESGCM, error) {
	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 key: %w", err)
	}
	return NewAESGCM(key)
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

// 适用于需要将加密结果存储为字符串的场景，如：
// - 数据库存储
// - JSON/XML序列化
// - URL参数传递
// 返回的字符串可以直接存储或传输
// EncryptToString 加密并返回格式化的字符串
// 格式: AES_{key_size}_GCM_V{version}:{key_fingerprint}:{base64_data}
// 示例: AES_256_GCM_V1:a1b2c3d4e5f6g7h8:Base64EncodedData...
func (a *AESGCM) EncryptToString(plaintext string) (string, error) {
	// 加密数据
	ciphertext, err := a.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}

	// 计算密钥指纹（前8字节的SHA256）
	keyFingerprint := sha256.Sum256(a.key)
	shortFingerprint := hex.EncodeToString(keyFingerprint[:8])

	// 确定密钥长度
	var keySizeStr string
	switch len(a.key) {
	case 16:
		keySizeStr = "128"
	case 24:
		keySizeStr = "192"
	case 32:
		keySizeStr = "256"
	default:
		return "", ErrInvalidKeySize
	}

	// 构建格式化字符串
	algorithm := fmt.Sprintf("AES_%s_GCM_V%d", keySizeStr, FormatVersion1)
	base64Data := base64.StdEncoding.EncodeToString(ciphertext)

	return fmt.Sprintf("%s:%s:%s", algorithm, shortFingerprint, base64Data), nil
}

// 对应EncryptToString方法的输出
// 适用于从字符串格式恢复加密数据的场景
// DecryptFromString 从格式化字符串解密
func (a *AESGCM) DecryptFromString(formatted string) (string, error) {
	parts := strings.Split(formatted, ":")
	if len(parts) != 3 {
		return "", errors.New("invalid formatted string: expected 3 parts")
	}

	algorithmPart := parts[0]
	keyFingerprint := parts[1]
	base64Data := parts[2]

	// 解析算法信息
	if !strings.HasPrefix(algorithmPart, "AES_") ||
		!strings.Contains(algorithmPart, "_GCM_V") {
		return "", errors.New("invalid algorithm format")
	}

	// 验证密钥长度匹配
	expectedKeySize := len(a.key)
	var expectedKeySizeStr string
	switch expectedKeySize {
	case 16:
		expectedKeySizeStr = "128"
	case 24:
		expectedKeySizeStr = "192"
	case 32:
		expectedKeySizeStr = "256"
	default:
		return "", ErrInvalidKeySize
	}

	if !strings.Contains(algorithmPart, "AES_"+expectedKeySizeStr+"_GCM") {
		return "", fmt.Errorf("key size mismatch: expected %s, got %s",
			expectedKeySizeStr, algorithmPart)
	}

	// 验证密钥指纹（可选，提供额外的安全性）
	currentFingerprint := sha256.Sum256(a.key)
	currentShortFingerprint := hex.EncodeToString(currentFingerprint[:8])
	if keyFingerprint != currentShortFingerprint {
		return "", errors.New("key fingerprint mismatch")
	}

	// 解码base64数据
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// 解密数据
	plaintext, err := a.Decrypt(data)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncryptToStringWithMetadata 加密并返回包含完整元数据的格式化字符串
// 格式: AES_{key_size}_GCM_V{version}:{key_fingerprint}:{base64_metadata}:{base64_ciphertext}
// 元数据包含: [版本1字节][算法1字节][密钥指纹8字节][nonce12字节]
func (a *AESGCM) EncryptToStringWithMetadata(plaintext string) (string, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 计算密钥指纹
	keyFingerprint := sha256.Sum256(a.key)
	shortFingerprint := hex.EncodeToString(keyFingerprint[:8])

	// 确定密钥长度
	var keySizeStr string
	switch len(a.key) {
	case 16:
		keySizeStr = "128"
	case 24:
		keySizeStr = "192"
	case 32:
		keySizeStr = "256"
	default:
		return "", ErrInvalidKeySize
	}

	// 构建元数据: 版本(1) + 算法(1) + 密钥指纹(8) + nonce(12)
	metadata := make([]byte, 0, 1+1+8+len(nonce))
	metadata = append(metadata, FormatVersion1)        // 版本
	metadata = append(metadata, AlgorithmAESGCM)       // 算法
	metadata = append(metadata, keyFingerprint[:8]...) // 密钥指纹
	metadata = append(metadata, nonce...)              // nonce

	// 加密数据
	ciphertextWithTag := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// 编码为base64
	metadataBase64 := base64.StdEncoding.EncodeToString(metadata)
	ciphertextBase64 := base64.StdEncoding.EncodeToString(ciphertextWithTag)

	// 构建格式化字符串
	algorithm := fmt.Sprintf("AES_%s_GCM_V%d", keySizeStr, FormatVersion1)

	return fmt.Sprintf("%s:%s:%s:%s", algorithm, shortFingerprint, metadataBase64, ciphertextBase64), nil
}

// DecryptFromStringWithMetadata 从包含完整元数据的格式化字符串解密
func (a *AESGCM) DecryptFromStringWithMetadata(formatted string) (string, error) {
	parts := strings.Split(formatted, ":")
	if len(parts) != 4 {
		return "", errors.New("invalid formatted string: expected 4 parts")
	}

	algorithmPart := parts[0]
	keyFingerprint := parts[1]
	metadataBase64 := parts[2]
	ciphertextBase64 := parts[3]

	// 验证算法格式
	if !strings.HasPrefix(algorithmPart, "AES_") ||
		!strings.Contains(algorithmPart, "_GCM_V") {
		return "", errors.New("invalid algorithm format")
	}

	// 验证密钥指纹
	currentFingerprint := sha256.Sum256(a.key)
	currentShortFingerprint := hex.EncodeToString(currentFingerprint[:8])
	if keyFingerprint != currentShortFingerprint {
		return "", errors.New("key fingerprint mismatch")
	}

	// 解码元数据
	metadata, err := base64.StdEncoding.DecodeString(metadataBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode metadata base64: %w", err)
	}

	if len(metadata) < 1+1+8+12 { // 版本+算法+指纹+nonce
		return "", errors.New("invalid metadata length")
	}

	// 解析元数据
	version := metadata[0]
	algorithm := metadata[1]
	storedFingerprint := metadata[2:10]
	nonce := metadata[10:22]

	// 验证版本和算法
	if version != FormatVersion1 {
		return "", errors.New("unsupported format version")
	}
	if algorithm != AlgorithmAESGCM {
		return "", errors.New("unsupported algorithm")
	}

	// 验证存储的密钥指纹
	if !bytes.Equal(storedFingerprint, currentFingerprint[:8]) {
		return "", errors.New("stored key fingerprint mismatch")
	}

	// 解码密文
	ciphertextWithTag, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext base64: %w", err)
	}

	// 解密数据
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertextWithTag, nil)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	return string(plaintext), nil
}

// 辅助函数：从格式化字符串中提取信息
func ParseFormattedString(formatted string) (algorithm, keyFingerprint string, err error) {
	parts := strings.Split(formatted, ":")
	if len(parts) < 2 {
		return "", "", errors.New("invalid formatted string")
	}
	return parts[0], parts[1], nil
}

// 辅助函数：验证格式化字符串是否由当前密钥加密
func (a *AESGCM) IsEncryptedByMe(formatted string) bool {
	_, fingerprint, err := ParseFormattedString(formatted)
	if err != nil {
		return false
	}

	currentFingerprint := sha256.Sum256(a.key)
	currentShortFingerprint := hex.EncodeToString(currentFingerprint[:8])

	return fingerprint == currentShortFingerprint
}

// GetKeyInfo 获取当前密钥的信息
func (a *AESGCM) GetKeyInfo() (keySize string, fingerprint string) {
	// 确定密钥长度
	switch len(a.key) {
	case 16:
		keySize = "128"
	case 24:
		keySize = "192"
	case 32:
		keySize = "256"
	default:
		keySize = "unknown"
	}

	// 计算密钥指纹（前8字节的SHA256）
	keyFingerprint := sha256.Sum256(a.key)
	fingerprint = hex.EncodeToString(keyFingerprint[:8])

	return keySize, fingerprint
}

// GetKeyInfoDetailed 获取更详细的密钥信息
func (a *AESGCM) GetKeyInfoDetailed() map[string]any {
	info := make(map[string]any)

	// 基本密钥信息
	info["key_length_bytes"] = len(a.key)

	switch len(a.key) {
	case 16:
		info["key_size"] = "128"
		info["key_algorithm"] = "AES-128"
	case 24:
		info["key_size"] = "192"
		info["key_algorithm"] = "AES-192"
	case 32:
		info["key_size"] = "256"
		info["key_algorithm"] = "AES-256"
	default:
		info["key_size"] = "unknown"
		info["key_algorithm"] = "AES-Unknown"
	}

	// 计算完整指纹和短指纹
	fullFingerprint := sha256.Sum256(a.key)
	info["key_fingerprint_full"] = hex.EncodeToString(fullFingerprint[:])
	info["key_fingerprint_short"] = hex.EncodeToString(fullFingerprint[:8])

	// 密钥前几个字节（用于调试，不包含完整密钥）
	if len(a.key) >= 4 {
		info["key_prefix"] = hex.EncodeToString(a.key[:4])
	}

	// 支持的加密模式
	info["supported_modes"] = []string{"GCM", "GCM with additional data"}

	return info
}

// GetKeySize 获取密钥大小
func (a *AESGCM) GetKeySize() int {
	return len(a.key)
}

// GetKeySizeBits 获取密钥位数
func (a *AESGCM) GetKeySizeBits() int {
	return len(a.key) * 8
}

// GetKeyFingerprint 获取密钥指纹
func (a *AESGCM) GetKeyFingerprint() string {
	keyFingerprint := sha256.Sum256(a.key)
	return hex.EncodeToString(keyFingerprint[:8])
}

// GetKeyFingerprintFull 获取完整密钥指纹
func (a *AESGCM) GetKeyFingerprintFull() string {
	keyFingerprint := sha256.Sum256(a.key)
	return hex.EncodeToString(keyFingerprint[:])
}

// ValidateKey 验证密钥是否有效
func (a *AESGCM) ValidateKey() error {
	if len(a.key) != 16 && len(a.key) != 24 && len(a.key) != 32 {
		return ErrInvalidKeySize
	}
	return nil
}

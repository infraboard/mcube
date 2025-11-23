// comparison_test.go
package aesgcm_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	aesgcm "github.com/infraboard/mcube/v2/crypto/aes_gcm"
)

// BenchmarkGCMvsCBC 性能对比测试
func BenchmarkGCMvsCBC(b *testing.B) {
	key, _ := aesgcm.GenerateKey(aesgcm.AES256)
	plaintext := make([]byte, 1024) // 1KB数据
	rand.Read(plaintext)

	// AES-GCM
	gcm, _ := aesgcm.NewAESGCM(key)

	b.Run("AES-GCM", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := gcm.Encrypt(plaintext)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	// AES-CBC + HMAC (模拟传统方式)
	b.Run("AES-CBC+HMAC", func(b *testing.B) {
		block, _ := aes.NewCipher(key)
		iv := make([]byte, aes.BlockSize)
		macKey := make([]byte, 32) // HMAC密钥

		for i := 0; i < b.N; i++ {
			// 生成随机IV
			rand.Read(iv)

			// 加密
			mode := cipher.NewCBCEncrypter(block, iv)
			ciphertext := make([]byte, len(plaintext))
			mode.CryptBlocks(ciphertext, plaintext)

			// 计算HMAC
			h := hmac.New(sha256.New, macKey)
			h.Write(ciphertext)
			mac := h.Sum(nil)

			_ = append(iv, append(ciphertext, mac...)...)
		}
	})
}

// TestGCMAuthentication 测试GCM认证功能
func TestGCMAuthentication(t *testing.T) {
	key, _ := aesgcm.GenerateKey(aesgcm.AES256)
	crypto, _ := aesgcm.NewAESGCM(key)

	plaintext := "重要数据"
	ciphertext, _ := crypto.Encrypt([]byte(plaintext))

	// 测试篡改检测
	tampered := make([]byte, len(ciphertext))
	copy(tampered, ciphertext)
	if len(tampered) > 20 {
		tampered[20] ^= 0x01 // 修改一个字节
	}

	_, err := crypto.Decrypt(tampered)
	if err == nil {
		t.Error("GCM应该检测到数据篡改，但没有报错")
	} else {
		fmt.Printf("GCM成功检测到数据篡改: %v\n", err)
	}
}

// Example_compareWithCBC 展示GCM与CBC的对比
func Example_compareWithCBC() {
	key, _ := aesgcm.GenerateKey(aesgcm.AES256)
	data := "需要加密的数据"

	fmt.Println("=== AES-GCM 方式 ===")
	// GCM - 一步完成加密和认证
	gcm, _ := aesgcm.NewAESGCM(key)
	gcmCiphertext, _ := gcm.Encrypt([]byte(data))
	gcmDecrypted, _ := gcm.Decrypt(gcmCiphertext)
	fmt.Printf("GCM解密: %s\n", string(gcmDecrypted))

	fmt.Println("=== AES-CBC+HMAC 方式 ===")
	// CBC + HMAC - 需要多个步骤
	block, _ := aes.NewCipher(key)

	// 加密
	iv := make([]byte, aes.BlockSize)
	rand.Read(iv)
	mode := cipher.NewCBCEncrypter(block, iv)
	paddedData := pkcs7Pad([]byte(data), aes.BlockSize)
	cbcCiphertext := make([]byte, len(paddedData))
	mode.CryptBlocks(cbcCiphertext, paddedData)

	// 计算HMAC（需要额外的密钥）
	macKey := make([]byte, 32)
	rand.Read(macKey)
	h := hmac.New(sha256.New, macKey)
	h.Write(cbcCiphertext)
	mac := h.Sum(nil)

	// 组合结果: IV + ciphertext + MAC
	finalCiphertext := append(iv, append(cbcCiphertext, mac...)...)

	// 解密过程同样复杂...
	fmt.Printf("CBC+HMAC完成，结果长度: %d\n", len(finalCiphertext))

	fmt.Println("=== 总结 ===")
	fmt.Println("GCM优势:")
	fmt.Println("  - 单一步骤完成加密和认证")
	fmt.Println("  - 无需填充处理")
	fmt.Println("  - 代码更简洁")
	fmt.Println("  - 性能通常更好")
}

// pkcs7Pad PKCS7填充函数
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// Example_nonceUniqueness 展示Nonce唯一性的重要性
func Example_nonceUniqueness() {
	key, _ := aesgcm.GenerateKey(aesgcm.AES256)
	crypto, _ := aesgcm.NewAESGCM(key)

	// 正确用法：每次加密自动生成新的随机nonce
	message1 := "第一条消息"
	message2 := "第二条消息"

	ciphertext1, _ := crypto.Encrypt([]byte(message1))
	ciphertext2, _ := crypto.Encrypt([]byte(message2))

	// 提取nonce进行比较（仅用于演示）
	nonce1 := ciphertext1[:12] // GCM通常使用12字节nonce
	nonce2 := ciphertext2[:12]

	fmt.Printf("Nonce1: %s\n", hex.EncodeToString(nonce1))
	fmt.Printf("Nonce2: %s\n", hex.EncodeToString(nonce2))
	fmt.Printf("Nonce是否唯一: %t\n", !hmac.Equal(nonce1, nonce2))

	// 安全提示
	fmt.Println("安全提示: 永远不要重复使用nonce!")
	fmt.Println("重复使用nonce会严重破坏GCM的安全性")
}

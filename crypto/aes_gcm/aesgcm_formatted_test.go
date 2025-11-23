// aesgcm_formatted_test.go
package aesgcm_test

import (
	"fmt"
	"strings"
	"testing"

	aesgcm "github.com/infraboard/mcube/v2/crypto/aes_gcm"
)

// TestEncryptToString 测试格式化字符串加密
func TestEncryptToString(t *testing.T) {
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)

	plaintext := "测试格式化字符串加密"
	formatted, err := crypto.EncryptToString(plaintext)
	if err != nil {
		t.Fatalf("EncryptToString failed: %v", err)
	}

	// 验证格式
	parts := strings.Split(formatted, ":")
	if len(parts) != 3 {
		t.Errorf("Expected 3 parts, got %d", len(parts))
	}

	// 验证算法格式
	if !strings.HasPrefix(parts[0], "AES_") {
		t.Errorf("Algorithm part should start with 'AES_', got: %s", parts[0])
	}

	if !strings.Contains(parts[0], "_GCM_V") {
		t.Errorf("Algorithm part should contain '_GCM_V', got: %s", parts[0])
	}

	// 验证密钥指纹长度
	if len(parts[1]) != 16 { // 8字节的hex编码是16字符
		t.Errorf("Key fingerprint should be 16 characters, got: %d", len(parts[1]))
	}

	// 验证解密
	decrypted, err := crypto.DecryptFromString(formatted)
	if err != nil {
		t.Fatalf("DecryptFromString failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Decrypted text doesn't match original. Got: %s, Want: %s", decrypted, plaintext)
	}
}

// TestEncryptToStringWithMetadata 测试带元数据的格式化字符串加密
func TestEncryptToStringWithMetadata(t *testing.T) {
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)

	plaintext := "测试带元数据的格式化字符串加密"
	formatted, err := crypto.EncryptToStringWithMetadata(plaintext)
	if err != nil {
		t.Fatalf("EncryptToStringWithMetadata failed: %v", err)
	}

	// 验证格式
	parts := strings.Split(formatted, ":")
	if len(parts) != 4 {
		t.Errorf("Expected 4 parts, got %d", len(parts))
	}

	// 验证解密
	decrypted, err := crypto.DecryptFromStringWithMetadata(formatted)
	if err != nil {
		t.Fatalf("DecryptFromStringWithMetadata failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Decrypted text doesn't match original. Got: %s, Want: %s", decrypted, plaintext)
	}
}

// TestDifferentKeySizesFormatted 测试不同密钥长度的格式化字符串
func TestDifferentKeySizesFormatted(t *testing.T) {
	testCases := []struct {
		name        string
		keySize     aesgcm.KeySize
		expectedStr string
	}{
		{"AES-128", aesgcm.AES128, "128"},
		{"AES-192", aesgcm.AES192, "192"},
		{"AES-256", aesgcm.AES256, "256"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key, _ := aesgcm.GenerateKey(tc.keySize)
			crypto, _ := aesgcm.NewAESGCM(key)

			plaintext := "测试数据"
			formatted, err := crypto.EncryptToString(plaintext)
			if err != nil {
				t.Fatalf("EncryptToFormattedString failed: %v", err)
			}

			// 验证算法部分包含正确的密钥长度
			parts := strings.Split(formatted, ":")
			algorithm := parts[0]
			expectedKeySize := fmt.Sprintf("AES_%s_GCM", tc.expectedStr)
			if !strings.Contains(algorithm, expectedKeySize) {
				t.Errorf("Algorithm should contain %s, got: %s", expectedKeySize, algorithm)
			}

			// 验证解密
			decrypted, err := crypto.DecryptFromString(formatted)
			if err != nil {
				t.Fatalf("DecryptFromFormattedString failed: %v", err)
			}

			if decrypted != plaintext {
				t.Errorf("Decrypted text doesn't match original")
			}
		})
	}
}

// TestParseFormattedString 测试解析格式化字符串
func TestParseFormattedString(t *testing.T) {
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)

	plaintext := "测试解析功能"
	formatted, _ := crypto.EncryptToString(plaintext)

	algorithm, fingerprint, err := aesgcm.ParseFormattedString(formatted)
	if err != nil {
		t.Fatalf("ParseFormattedString failed: %v", err)
	}

	if !strings.HasPrefix(algorithm, "AES_") {
		t.Errorf("Parsed algorithm should start with 'AES_', got: %s", algorithm)
	}

	if len(fingerprint) != 16 {
		t.Errorf("Parsed fingerprint should be 16 characters, got: %d", len(fingerprint))
	}
}

// TestIsEncryptedByMe 测试密钥匹配验证
func TestIsEncryptedByMe(t *testing.T) {
	key1, _ := aesgcm.Generate256Key()
	crypto1, _ := aesgcm.NewAESGCM(key1)

	key2, _ := aesgcm.Generate256Key()
	crypto2, _ := aesgcm.NewAESGCM(key2)

	plaintext := "测试密钥匹配"
	formatted, _ := crypto1.EncryptToString(plaintext)

	// 使用加密的密钥验证
	if !crypto1.IsEncryptedByMe(formatted) {
		t.Error("IsEncryptedByMe should return true for correct key")
	}

	// 使用不同的密钥验证
	if crypto2.IsEncryptedByMe(formatted) {
		t.Error("IsEncryptedByMe should return false for different key")
	}
}

// TestFormattedStringErrorHandling 测试错误处理
func TestFormattedStringErrorHandling(t *testing.T) {
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)

	// 测试无效格式
	_, err := crypto.DecryptFromString("invalid:format")
	if err == nil {
		t.Error("Should return error for invalid format")
	}

	// 测试部分格式
	_, err = crypto.DecryptFromString("AES_256_GCM_V1:fingerprint")
	if err == nil {
		t.Error("Should return error for incomplete format")
	}

	// 测试错误的算法格式
	_, err = crypto.DecryptFromString("INVALID_ALGO:fingerprint:data")
	if err == nil {
		t.Error("Should return error for invalid algorithm")
	}

	// 测试错误的密钥指纹
	wrongFormatted := "AES_256_GCM_V1:wrongfingerprint:base64data"
	_, err = crypto.DecryptFromString(wrongFormatted)
	if err == nil {
		t.Error("Should return error for wrong key fingerprint")
	}
}

// TestFormattedStringWithMetadataErrorHandling 测试带元数据的错误处理
func TestFormattedStringWithMetadataErrorHandling(t *testing.T) {
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)

	// 测试无效格式
	_, err := crypto.DecryptFromStringWithMetadata("invalid:format:with:too:many:parts")
	if err == nil {
		t.Error("Should return error for invalid format")
	}

	// 测试无效的base64数据
	_, err = crypto.DecryptFromStringWithMetadata("AES_256_GCM_V1:fingerprint:invalid-base64:data")
	if err == nil {
		t.Error("Should return error for invalid base64")
	}
}

// BenchmarkFormattedStringEncrypt 格式化字符串加密性能测试
func BenchmarkFormattedStringEncrypt(b *testing.B) {
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)
	plaintext := "This is a test plaintext for benchmarking formatted string encryption"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := crypto.EncryptToString(plaintext)
		if err != nil {
			b.Fatalf("EncryptToString failed: %v", err)
		}
	}
}

// BenchmarkFormattedStringDecrypt 格式化字符串解密性能测试
func BenchmarkFormattedStringDecrypt(b *testing.B) {
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)
	plaintext := "This is a test plaintext for benchmarking formatted string decryption"
	formatted, _ := crypto.EncryptToString(plaintext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := crypto.DecryptFromString(formatted)
		if err != nil {
			b.Fatalf("DecryptFromString failed: %v", err)
		}
	}
}

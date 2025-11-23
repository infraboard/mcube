package aesgcm_test

import (
	"strings"
	"testing"

	aesgcm "github.com/infraboard/mcube/v2/crypto/aes_gcm"
)

// TestKeyGeneration 测试密钥生成
func TestKeyGeneration(t *testing.T) {
	tests := []struct {
		name    string
		keySize aesgcm.KeySize
		wantErr bool
		errMsg  string
	}{
		{
			name:    "generate 128 bit key",
			keySize: 16,
			wantErr: false,
		},
		{
			name:    "generate 192 bit key",
			keySize: 24,
			wantErr: false,
		},
		{
			name:    "generate 256 bit key",
			keySize: 32,
			wantErr: false,
		},
		{
			name:    "generate invalid key size",
			keySize: 20,
			wantErr: true,
			errMsg:  aesgcm.ErrInvalidKeySize.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := aesgcm.GenerateKey(tt.keySize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err.Error() != tt.errMsg {
					t.Errorf("GenerateKey() error message = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if len(key) != int(tt.keySize) {
				t.Errorf("GenerateKey() key length = %d, want %d", len(key), tt.keySize)
			}
		})
	}
}

// TestNewAESGCM 测试创建AESGCM实例
func TestNewAESGCM(t *testing.T) {
	tests := []struct {
		name    string
		key     []byte
		wantErr bool
	}{
		{
			name:    "valid 128 bit key",
			key:     make([]byte, 16),
			wantErr: false,
		},
		{
			name:    "valid 192 bit key",
			key:     make([]byte, 24),
			wantErr: false,
		},
		{
			name:    "valid 256 bit key",
			key:     make([]byte, 32),
			wantErr: false,
		},
		{
			name:    "invalid key size",
			key:     make([]byte, 20),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := aesgcm.NewAESGCM(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAESGCM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestEncryptDecrypt 测试加密解密功能
func TestEncryptDecrypt(t *testing.T) {
	key, err := aesgcm.GenerateKey(32)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	crypto, err := aesgcm.NewAESGCM(key)
	if err != nil {
		t.Fatalf("Failed to create AESGCM: %v", err)
	}

	testCases := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "short text",
			plaintext: "hello",
		},
		{
			name:      "long text",
			plaintext: strings.Repeat("This is a longer text for testing. ", 10),
		},
		{
			name:      "special characters",
			plaintext: "Hello, 世界! 🎉 特殊字符测试",
		},
		{
			name:      "binary data",
			plaintext: string([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 测试字节数组加密解密
			ciphertext, err := crypto.Encrypt([]byte(tc.plaintext))
			if err != nil {
				t.Errorf("Encrypt() failed: %v", err)
				return
			}

			decrypted, err := crypto.Decrypt(ciphertext)
			if err != nil {
				t.Errorf("Decrypt() failed: %v", err)
				return
			}

			if string(decrypted) != tc.plaintext {
				t.Errorf("Decrypted text doesn't match original. Got: %s, Want: %s",
					string(decrypted), tc.plaintext)
			}

			// 测试字符串接口
			encoded, err := crypto.EncryptToString(tc.plaintext)
			if err != nil {
				t.Errorf("EncryptToString() failed: %v", err)
				return
			}

			decryptedStr, err := crypto.DecryptFromString(encoded)
			if err != nil {
				t.Errorf("DecryptFromString() failed: %v", err)
				return
			}

			if decryptedStr != tc.plaintext {
				t.Errorf("Decrypted string doesn't match original. Got: %s, Want: %s",
					decryptedStr, tc.plaintext)
			}
		})
	}
}

// TestEncryptDecryptWithAdditionalData 测试带附加数据的加密解密
func TestEncryptDecryptWithAdditionalData(t *testing.T) {
	key, err := aesgcm.GenerateKey(32)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	crypto, err := aesgcm.NewAESGCM(key)
	if err != nil {
		t.Fatalf("Failed to create AESGCM: %v", err)
	}

	plaintext := "sensitive data"
	additionalData := []byte("authentication data")

	// 使用正确的附加数据加密解密
	ciphertext, err := crypto.EncryptWithAdditionalData([]byte(plaintext), additionalData)
	if err != nil {
		t.Fatalf("EncryptWithAdditionalData() failed: %v", err)
	}

	decrypted, err := crypto.DecryptWithAdditionalData(ciphertext, additionalData)
	if err != nil {
		t.Fatalf("DecryptWithAdditionalData() failed: %v", err)
	}

	if string(decrypted) != plaintext {
		t.Errorf("Decrypted text doesn't match original. Got: %s, Want: %s",
			string(decrypted), plaintext)
	}

	// 使用错误的附加数据应该失败
	wrongAdditionalData := []byte("wrong authentication data")
	_, err = crypto.DecryptWithAdditionalData(ciphertext, wrongAdditionalData)
	if err == nil {
		t.Error("Expected decryption to fail with wrong additional data, but it succeeded")
	}
}

// TestDifferentKeys 测试不同密钥的隔离性
func TestDifferentKeys(t *testing.T) {
	key1, _ := aesgcm.GenerateKey(32)
	key2, _ := aesgcm.GenerateKey(32)

	crypto1, _ := aesgcm.NewAESGCM(key1)
	crypto2, _ := aesgcm.NewAESGCM(key2)

	plaintext := "test data"

	// 使用密钥1加密
	ciphertext, err := crypto1.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("Encrypt with key1 failed: %v", err)
	}

	// 使用密钥2解密应该失败
	_, err = crypto2.Decrypt(ciphertext)
	if err == nil {
		t.Error("Expected decryption to fail with different key, but it succeeded")
	}

	// 使用密钥1解密应该成功
	decrypted, err := crypto1.Decrypt(ciphertext)
	if err != nil {
		t.Errorf("Decryption with key1 failed: %v", err)
	}

	if string(decrypted) != plaintext {
		t.Errorf("Decrypted text doesn't match original")
	}
}

// TestTamperedCiphertext 测试篡改密文的情况
func TestTamperedCiphertext(t *testing.T) {
	key, _ := aesgcm.GenerateKey(32)
	crypto, _ := aesgcm.NewAESGCM(key)

	plaintext := "test data"
	ciphertext, err := crypto.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// 篡改密文
	tampered := make([]byte, len(ciphertext))
	copy(tampered, ciphertext)
	if len(tampered) > 10 {
		tampered[10] ^= 0x01 // 修改一个字节
	}

	// 解密应该失败
	_, err = crypto.Decrypt(tampered)
	if err == nil {
		t.Error("Expected decryption to fail with tampered ciphertext, but it succeeded")
	}
}

// TestInvalidCiphertext 测试无效密文
func TestInvalidCiphertext(t *testing.T) {
	key, _ := aesgcm.GenerateKey(32)
	crypto, _ := aesgcm.NewAESGCM(key)

	// 测试过短的密文
	shortCiphertext := []byte{0x01, 0x02, 0x03}
	_, err := crypto.Decrypt(shortCiphertext)
	if err == nil {
		t.Error("Expected decryption to fail with short ciphertext, but it succeeded")
	}

	// 测试无效的base64字符串
	_, err = crypto.DecryptFromString("invalid base64!!!")
	if err == nil {
		t.Error("Expected decryption to fail with invalid base64, but it succeeded")
	}
}

// TestAES256GCM 测试AES-256专用函数
func TestAES256GCM(t *testing.T) {
	// 测试生成256位密钥
	key, err := aesgcm.Generate256Key()
	if err != nil {
		t.Fatalf("Generate256Key failed: %v", err)
	}

	if len(key) != 32 {
		t.Errorf("Generate256Key returned key of length %d, want 32", len(key))
	}

	// 测试创建AES-256实例
	_, err = aesgcm.NewAES256GCM(key)
	if err != nil {
		t.Fatalf("NewAES256GCM failed: %v", err)
	}

	// 测试使用错误长度的密钥
	_, err = aesgcm.NewAES256GCM([]byte("short key"))
	if err == nil {
		t.Error("Expected NewAES256GCM to fail with short key, but it succeeded")
	}
}

// TestNewAESGCMFromString 测试从字符串创建实例
func TestNewAESGCMFromString(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		keySize aesgcm.KeySize
		wantErr bool
	}{
		{
			name:    "valid 128 bit key string",
			key:     "1234567890123456", // 16 bytes
			keySize: 16,
			wantErr: false,
		},
		{
			name:    "valid 256 bit key string",
			key:     "12345678901234567890123456789012", // 32 bytes
			keySize: 32,
			wantErr: false,
		},
		{
			name:    "invalid key string length",
			key:     "short",
			keySize: 32,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crypto, err := aesgcm.NewAESGCMFromString(tt.key, tt.keySize)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAESGCMFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 测试加密解密功能
				plaintext := "test data"
				ciphertext, err := crypto.Encrypt([]byte(plaintext))
				if err != nil {
					t.Errorf("Encrypt failed: %v", err)
					return
				}

				decrypted, err := crypto.Decrypt(ciphertext)
				if err != nil {
					t.Errorf("Decrypt failed: %v", err)
					return
				}

				if string(decrypted) != plaintext {
					t.Errorf("Decrypted text doesn't match original")
				}
			}
		})
	}
}

// TestMultipleEncryptions 测试多次加密产生不同的结果
func TestMultipleEncryptions(t *testing.T) {
	key, _ := aesgcm.GenerateKey(32)
	crypto, _ := aesgcm.NewAESGCM(key)

	plaintext := "same plaintext"

	// 多次加密相同明文
	ciphertext1, err := crypto.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("First encryption failed: %v", err)
	}

	ciphertext2, err := crypto.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("Second encryption failed: %v", err)
	}

	// 由于随机nonce，两次加密结果应该不同
	if string(ciphertext1) == string(ciphertext2) {
		t.Error("Multiple encryptions of same plaintext produced identical results")
	}

	// 但两次解密结果应该相同
	decrypted1, err := crypto.Decrypt(ciphertext1)
	if err != nil {
		t.Errorf("First decryption failed: %v", err)
	}

	decrypted2, err := crypto.Decrypt(ciphertext2)
	if err != nil {
		t.Errorf("Second decryption failed: %v", err)
	}

	if string(decrypted1) != string(decrypted2) {
		t.Error("Decryptions of different ciphertexts produced different results")
	}
}

// BenchmarkEncrypt 加密性能测试
func BenchmarkEncrypt(b *testing.B) {
	key, _ := aesgcm.GenerateKey(32)
	crypto, _ := aesgcm.NewAESGCM(key)
	plaintext := []byte("This is a test plaintext for benchmarking")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := crypto.Encrypt(plaintext)
		if err != nil {
			b.Fatalf("Encrypt failed: %v", err)
		}
	}
}

// BenchmarkDecrypt 解密性能测试
func BenchmarkDecrypt(b *testing.B) {
	key, _ := aesgcm.GenerateKey(32)
	crypto, _ := aesgcm.NewAESGCM(key)
	plaintext := []byte("This is a test plaintext for benchmarking")
	ciphertext, _ := crypto.Encrypt(plaintext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := crypto.Decrypt(ciphertext)
		if err != nil {
			b.Fatalf("Decrypt failed: %v", err)
		}
	}
}

func TestKeySizeEnum(t *testing.T) {
	tests := []struct {
		name     string
		keySize  aesgcm.KeySize
		valid    bool
		expected string
	}{
		{"AES128", aesgcm.AES128, true, "AES-128"},
		{"AES192", aesgcm.AES192, true, "AES-192"},
		{"AES256", aesgcm.AES256, true, "AES-256"},
		{"Invalid", aesgcm.KeySize(20), false, "Unknown(20)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试有效性检查
			if valid := tt.keySize.Valid(); valid != tt.valid {
				t.Errorf("Valid() = %v, want %v", valid, tt.valid)
			}

			// 测试字符串表示
			if str := tt.keySize.String(); str != tt.expected {
				t.Errorf("String() = %v, want %v", str, tt.expected)
			}
		})
	}
}

func TestGenerateKeyWithEnum(t *testing.T) {
	tests := []struct {
		name    string
		keySize aesgcm.KeySize
		wantErr bool
	}{
		{"AES128", aesgcm.AES128, false},
		{"AES192", aesgcm.AES192, false},
		{"AES256", aesgcm.AES256, false},
		{"Invalid", aesgcm.KeySize(20), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := aesgcm.GenerateKey(tt.keySize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(key) != int(tt.keySize) {
				t.Errorf("GenerateKey() key length = %d, want %d", len(key), tt.keySize)
			}
		})
	}
}

// 测试新的便捷函数
func TestConvenienceKeyGeneration(t *testing.T) {
	// 测试256位密钥生成
	key256, err := aesgcm.Generate256Key()
	if err != nil {
		t.Errorf("Generate256Key failed: %v", err)
	}
	if len(key256) != 32 {
		t.Errorf("Generate256Key length = %d, want 32", len(key256))
	}

	// 测试128位密钥生成
	key128, err := aesgcm.Generate128Key()
	if err != nil {
		t.Errorf("Generate128Key failed: %v", err)
	}
	if len(key128) != 16 {
		t.Errorf("Generate128Key length = %d, want 16", len(key128))
	}

	// 测试192位密钥生成
	key192, err := aesgcm.Generate192Key()
	if err != nil {
		t.Errorf("Generate192Key failed: %v", err)
	}
	if len(key192) != 24 {
		t.Errorf("Generate192Key length = %d, want 24", len(key192))
	}
}

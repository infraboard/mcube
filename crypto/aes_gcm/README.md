# AES-GCM 加密库

## 使用示例

```go
// package_example.go
package aesgcm

// 使用指南示例代码，展示不同方法的使用场景

// ExampleUsage 展示不同加密方法的使用场景
func ExampleUsage() {
	// 场景1: 基础加密 - 适用于大多数简单加密需求
	func ExampleBasicEncryption() {
		key, _ := GenerateKey(AES256)
		crypto, _ := NewAESGCM(key)
		
		// 加密普通数据
		plaintext := "敏感数据"
		ciphertext, _ := crypto.Encrypt([]byte(plaintext))
		
		// 解密
		decrypted, _ := crypto.Decrypt(ciphertext)
		_ = string(decrypted)
	}

	// 场景2: 字符串加密 - 适用于需要字符串输出的场景
	func ExampleStringEncryption() {
		key, _ := GenerateKey(AES256)
		crypto, _ := NewAESGCM(key)
		
		// 加密为base64字符串，适合存储到数据库
		encoded, _ := crypto.EncryptToString("用户密码")
		
		// 从base64字符串解密
		decrypted, _ := crypto.DecryptFromString(encoded)
		_ = decrypted
	}

	// 场景3: 带附加数据的加密 - 适用于需要认证额外元数据的场景
	func ExampleAdditionalDataEncryption() {
		key, _ := GenerateKey(AES256)
		crypto, _ := NewAESGCM(key)
		
		// 加密消息体，同时认证消息头
		message := "消息内容"
		headers := []byte("消息类型:私密;发送者:Alice")
		
		ciphertext, _ := crypto.EncryptWithAdditionalData(
			[]byte(message), 
			headers,
		)
		
		// 解密时必须提供相同的消息头，否则会失败
		decrypted, _ := crypto.DecryptWithAdditionalData(ciphertext, headers)
		_ = string(decrypted)
	}

	// 场景4: 从配置创建实例 - 适用于从环境变量读取密钥
	func ExampleFromConfig() {
		// 假设从环境变量读取密钥
		keyStr := "my-secret-key-32-bytes-long-keys!"
		crypto, _ := NewAESGCMFromString(keyStr, AES256)
		
		ciphertext, _ := crypto.EncryptToString("配置数据")
		_ = ciphertext
	}
}
```


## 选择加密方法的指南

1. 基础需求 - 使用 Encrypt/Decrypt:
   - 简单的数据加密
   - 不需要额外的认证数据
   - 处理二进制数据

2. 字符串存储 - 使用 EncryptToString/DecryptFromString:
   - 加密结果需要存储到数据库
   - 需要在JSON/XML中传输加密数据
   - URL参数传递

3. 高级认证 - 使用 EncryptWithAdditsionalData/DecryptWithAdditionalData:
   - 需要保护消息头或元数据
   - 加密数据有关联的上下文信息
   - 需要额外的数据完整性验证

4. 密钥生成建议:
   - 高安全性: 使用 Generate256Key() 或 GenerateKey(AES256)
   - 性能优先: 使用 Generate128Key() 或 GenerateKey(AES128)
   - 兼容性: 使用 GenerateKey() 并指定合适的密钥长度
// example_test.go
package aesgcm_test

import (
	"fmt"
	"log"
	"strings"

	aesgcm "github.com/infraboard/mcube/v2/crypto/aes_gcm"
)

// Example_basicUsage 展示基础用法
func Example_basicUsage() {
	// 生成256位密钥（推荐用于生产环境）
	key, err := aesgcm.GenerateKey(aesgcm.AES256)
	if err != nil {
		log.Fatal(err)
	}

	// 创建加密器实例
	crypto, err := aesgcm.NewAESGCM(key)
	if err != nil {
		log.Fatal(err)
	}

	// 加密数据
	plaintext := "这是一段敏感数据，需要安全存储"
	ciphertext, err := crypto.Encrypt([]byte(plaintext))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("加密成功\n")

	// 解密数据
	decrypted, err := crypto.Decrypt(ciphertext)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("解密成功: %s\n", string(decrypted))
	// Output:
	// 加密成功
	// 解密成功: 这是一段敏感数据，需要安全存储
}

// Example_stringEncryption 展示字符串加密
func Example_stringEncryption() {
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)

	// 加密为base64字符串，适合存储到数据库
	userData := "用户私密信息"
	encoded, err := crypto.EncryptToString(userData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("加密成功\n")

	// 从数据库读取后解密
	decrypted, err := crypto.DecryptFromString(encoded)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("解密成功: %s\n", decrypted)
	// Output:
	// 加密成功
	// 解密成功: 用户私密信息
}

// Example_additionalData 展示附加数据的使用
func Example_additionalData() {
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)

	// 消息内容
	message := "转账给Bob: 1000元"
	// 消息头（不会被加密，但参与认证）
	headers := []byte("交易类型:转账;时间戳:2024-01-01")

	// 使用附加数据加密
	ciphertext, err := crypto.EncryptWithAdditionalData([]byte(message), headers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("加密成功，包含消息头认证\n")

	// 解密时必须提供相同的消息头
	decrypted, err := crypto.DecryptWithAdditionalData(ciphertext, headers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("解密成功: %s\n", string(decrypted))

	// 尝试使用错误的头信息解密（会失败）
	wrongHeaders := []byte("交易类型:查询")
	_, err = crypto.DecryptWithAdditionalData(ciphertext, wrongHeaders)
	if err != nil {
		fmt.Printf("认证失败: %v\n", err)
	}
	// Output:
	// 加密成功，包含消息头认证
	// 解密成功: 转账给Bob: 1000元
	// 认证失败: decryption failed: authentication failed
}

// Example_differentKeySizes 展示不同密钥长度的使用
func Example_differentKeySizes() {
	// 场景1: 高安全性需求 - 使用AES-256
	key256, _ := aesgcm.GenerateKey(aesgcm.AES256)
	crypto256, _ := aesgcm.NewAESGCM(key256)
	_, err := crypto256.EncryptToString("高度敏感数据")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("AES-256加密成功\n")

	// 场景2: 性能敏感 - 使用AES-128
	key128, _ := aesgcm.GenerateKey(aesgcm.AES128)
	crypto128, _ := aesgcm.NewAESGCM(key128)
	_, err = crypto128.EncryptToString("普通敏感数据")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("AES-128加密成功\n")

	// 场景3: 从字符串密钥创建 - 使用确切的32字节字符串
	keyStr := "0123456789abcdef0123456789abcdef" // 确切的32字节
	cryptoFromStr, err := aesgcm.NewAESGCMFromString(keyStr, aesgcm.AES256)
	if err != nil {
		log.Fatal(err)
	}
	_, err = cryptoFromStr.EncryptToString("配置数据")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("字符串密钥加密成功\n")
	// Output:
	// AES-256加密成功
	// AES-128加密成功
	// 字符串密钥加密成功
}

// Example_securityBestPractices 展示安全最佳实践
func Example_securityBestPractices() {
	// 1. 使用足够长度的密钥
	key, _ := aesgcm.Generate256Key()

	// 2. 创建加密器
	crypto, _ := aesgcm.NewAESGCM(key)

	// 3. 加密数据
	sensitiveData := "信用卡号: 1234-5678-9012-3456"
	ciphertext, _ := crypto.Encrypt([]byte(sensitiveData))

	// 4. 安全处理解密错误
	_, err := crypto.Decrypt(ciphertext)
	if err != nil {
		log.Printf("安全警告: 数据完整性验证失败 - %v", err)
		return
	}

	fmt.Printf("解密成功\n")
	// Output:
	// 解密成功
}

// Example_performanceConsiderations 展示性能考虑
func Example_performanceConsiderations() {
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)

	// 批量加密
	messages := []string{
		"第一条消息",
		"第二条消息",
		"测试消息",
	}

	successCount := 0
	for _, msg := range messages {
		_, err := crypto.Encrypt([]byte(msg))
		if err != nil {
			continue
		}
		successCount++
	}

	fmt.Printf("成功加密 %d 条消息\n", successCount)
	// Output:
	// 成功加密 3 条消息
}

// Example_encryptionWorkflow 展示完整的加密工作流程
func Example_encryptionWorkflow() {
	// 步骤1: 生成密钥
	key, err := aesgcm.Generate256Key()
	if err != nil {
		log.Fatal("密钥生成失败:", err)
	}
	fmt.Println("✓ 密钥生成成功")

	// 步骤2: 创建加密器
	crypto, err := aesgcm.NewAESGCM(key)
	if err != nil {
		log.Fatal("加密器创建失败:", err)
	}
	fmt.Println("✓ 加密器创建成功")

	// 步骤3: 加密数据
	plaintext := "Hello, AES-GCM!"
	ciphertext, err := crypto.Encrypt([]byte(plaintext))
	if err != nil {
		log.Fatal("加密失败:", err)
	}
	fmt.Println("✓ 数据加密成功")

	// 步骤4: 解密数据
	decrypted, err := crypto.Decrypt(ciphertext)
	if err != nil {
		log.Fatal("解密失败:", err)
	}
	fmt.Println("✓ 数据解密成功")

	// 步骤5: 验证数据完整性
	if string(decrypted) == plaintext {
		fmt.Println("✓ 数据完整性验证通过")
	} else {
		fmt.Println("✗ 数据完整性验证失败")
	}
	// Output:
	// ✓ 密钥生成成功
	// ✓ 加密器创建成功
	// ✓ 数据加密成功
	// ✓ 数据解密成功
	// ✓ 数据完整性验证通过
}

// Example_errorHandling 展示错误处理
func Example_errorHandling() {
	// 测试无效密钥
	_, err := aesgcm.NewAESGCM([]byte("short-key"))
	if err != nil {
		fmt.Println("捕获无效密钥错误")
	}

	// 测试无效密钥长度
	_, err = aesgcm.GenerateKey(aesgcm.KeySize(20))
	if err != nil {
		fmt.Println("捕获无效密钥长度错误")
	}

	// 测试无效密文
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)
	_, err = crypto.Decrypt([]byte("short"))
	if err != nil {
		fmt.Println("捕获无效密文错误")
	}

	// 测试篡改的密文
	plaintext := "测试数据"
	ciphertext, _ := crypto.Encrypt([]byte(plaintext))

	// 篡改密文
	tampered := make([]byte, len(ciphertext))
	copy(tampered, ciphertext)
	if len(tampered) > 20 {
		tampered[20] ^= 0x01
	}

	_, err = crypto.Decrypt(tampered)
	if err != nil {
		fmt.Println("捕获数据篡改错误")
	}
	// Output:
	// 捕获无效密钥错误
	// 捕获无效密钥长度错误
	// 捕获无效密文错误
	// 捕获数据篡改错误
}

// Example_base64Encoding 展示Base64编码的使用场景
func Example_base64Encoding() {
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)

	// 场景1: 存储到数据库
	userPassword := "mySecretPassword123"
	encryptedForDB, _ := crypto.EncryptToString(userPassword)

	// 场景2: URL安全传输
	apiKey := "api_key_123456789"
	encryptedForURL, _ := crypto.EncryptToString(apiKey)

	// 场景3: JSON序列化
	configData := `{"server": "api.example.com", "port": 443}`
	encryptedForJSON, _ := crypto.EncryptToString(configData)

	// 验证解密功能
	decrypted1, _ := crypto.DecryptFromString(encryptedForDB)
	decrypted2, _ := crypto.DecryptFromString(encryptedForURL)
	decrypted3, _ := crypto.DecryptFromString(encryptedForJSON)

	success := strings.Contains(decrypted1, "Password") &&
		strings.Contains(decrypted2, "api_key") &&
		strings.Contains(decrypted3, "example.com")

	fmt.Printf("Base64编码加解密验证: %v\n", success)
	// Output:
	// Base64编码加解密验证: true
}

// Example_keyManagement 展示密钥管理
func Example_keyManagement() {
	// 生成不同长度的密钥
	key128, _ := aesgcm.GenerateKey(aesgcm.AES128)
	key192, _ := aesgcm.GenerateKey(aesgcm.AES192)
	key256, _ := aesgcm.GenerateKey(aesgcm.AES256)

	fmt.Printf("支持AES-128: %v\n", len(key128) == 16)
	fmt.Printf("支持AES-192: %v\n", len(key192) == 24)
	fmt.Printf("支持AES-256: %v\n", len(key256) == 32)

	// 使用便捷函数
	key256Easy, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAES256GCM(key256Easy)

	plaintext := "测试数据"
	ciphertext, _ := crypto.Encrypt([]byte(plaintext))
	decrypted, _ := crypto.Decrypt(ciphertext)

	fmt.Printf("加解密功能正常: %v\n", string(decrypted) == plaintext)
	// Output:
	// 支持AES-128: true
	// 支持AES-192: true
	// 支持AES-256: true
	// 加解密功能正常: true
}

// Example_comprehensive 综合示例
func Example_comprehensive() {
	// 1. 密钥生成和管理
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)

	// 2. 基础加密解密
	data1 := "普通数据"
	encrypted1, _ := crypto.Encrypt([]byte(data1))
	decrypted1, _ := crypto.Decrypt(encrypted1)

	// 3. 字符串加密解密
	data2 := "字符串数据"
	encryptedStr, _ := crypto.EncryptToString(data2)
	decryptedStr, _ := crypto.DecryptFromString(encryptedStr)

	// 4. 带附加数据的加密
	data3 := "敏感数据"
	additionalData := []byte("元数据")
	encryptedWithAD, _ := crypto.EncryptWithAdditionalData([]byte(data3), additionalData)
	decryptedWithAD, _ := crypto.DecryptWithAdditionalData(encryptedWithAD, additionalData)

	// 验证所有功能正常
	allSuccess := string(decrypted1) == data1 &&
		decryptedStr == data2 &&
		string(decryptedWithAD) == data3

	fmt.Printf("AES-GCM所有功能测试: %v\n", allSuccess)
	// Output:
	// AES-GCM所有功能测试: true
}

// Example_stringKeyUsage 专门展示字符串密钥的使用
func Example_stringKeyUsage() {
	// 方法1: 使用字节数组（推荐）- 使用确切的32字节
	keyBytes := []byte("0123456789abcdef0123456789abcdef") // 确切的32字节
	crypto1, err := aesgcm.NewAESGCM(keyBytes)
	if err != nil {
		log.Fatal(err)
	}
	_, err = crypto1.EncryptToString("使用方法1加密的数据")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("字节数组密钥加密成功\n")

	// 方法2: 使用字符串（需要精确长度）
	keyStr := "abcdefghijklmnopqrstuvwxyz123456" // 32字节字符串
	crypto2, err := aesgcm.NewAESGCMFromString(keyStr, aesgcm.AES256)
	if err != nil {
		log.Fatal(err)
	}
	_, err = crypto2.EncryptToString("使用方法2加密的数据")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("字符串密钥加密成功\n")
	// Output:
	// 字节数组密钥加密成功
	// 字符串密钥加密成功
}

// Example_newAESGCMFromBase64 展示base64密钥的使用示例
func Example_newAESGCMFromBase64() {
	// 从配置文件读取base64编码的密钥
	keyBase64 := "4pSf+zOOucXZrnnNamN/HFcUy55bwCcw1HmCi5U5S9w="

	// 创建加密器
	crypto, err := aesgcm.NewAESGCMFromBase64(keyBase64)
	if err != nil {
		log.Fatal(err)
	}

	// 获取密钥信息
	keySize, fingerprint := crypto.GetKeyInfo()

	// 验证密钥长度
	if keySize == "256" {
		fmt.Println("密钥长度: 256位")
	}

	// 验证指纹格式（16字符的hex）
	if len(fingerprint) == 16 {
		fmt.Println("密钥指纹格式正确")
	}

	// 加密数据
	plaintext := "敏感数据"
	formatted, err := crypto.EncryptToString(plaintext)
	if err != nil {
		log.Fatal(err)
	}

	// 验证加密字符串格式
	if strings.HasPrefix(formatted, "AES_256_GCM_V1:") && len(formatted) > 30 {
		fmt.Println("加密成功")
	}

	// 解密数据
	decrypted, err := crypto.DecryptFromString(formatted)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("解密成功: %v\n", decrypted == plaintext)
	// Output:
	// 密钥长度: 256位
	// 密钥指纹格式正确
	// 加密成功
	// 解密成功: true
}

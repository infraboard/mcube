package aesgcm_test

import (
	"fmt"
	"log"
	"strings"

	aesgcm "github.com/infraboard/mcube/v2/crypto/aes_gcm"
)

// Example_EncryptToString 展示格式化字符串加密
func Example_encryptToString() {
	// 生成密钥
	key, err := aesgcm.Generate256Key()
	if err != nil {
		log.Fatal(err)
	}

	// 创建加密器
	crypto, err := aesgcm.NewAESGCM(key)
	if err != nil {
		log.Fatal(err)
	}

	// 加密为格式化字符串
	plaintext := "这是一段需要加密的敏感数据"
	formatted, err := crypto.EncryptToString(plaintext)
	if err != nil {
		log.Fatal(err)
	}

	// 只验证格式前缀，不验证具体内容
	if strings.HasPrefix(formatted, "AES_256_GCM_V1:") {
		fmt.Println("格式化加密字符串格式正确")
	}

	// 解密格式化字符串
	decrypted, err := crypto.DecryptFromString(formatted)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("解密成功: %v\n", decrypted == plaintext)
	// Output:
	// 格式化加密字符串格式正确
	// 解密成功: true
}

// Example_EncryptToStringWithMetadata 展示带元数据的格式化字符串加密
func Example_encryptToStringWithMetadata() {
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)

	// 加密为带元数据的格式化字符串
	plaintext := "这是一段需要加密的敏感数据"
	formatted, err := crypto.EncryptToStringWithMetadata(plaintext)
	if err != nil {
		log.Fatal(err)
	}

	// 只验证格式前缀
	if strings.HasPrefix(formatted, "AES_256_GCM_V1:") {
		fmt.Println("带元数据的格式化字符串格式正确")
	}

	// 解密
	decrypted, err := crypto.DecryptFromStringWithMetadata(formatted)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("解密成功: %v\n", decrypted == plaintext)
	// Output:
	// 带元数据的格式化字符串格式正确
	// 解密成功: true
}

// Example_parseFormattedString 展示解析格式化字符串
func Example_parseFormattedString() {
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)

	// 加密数据
	formatted, _ := crypto.EncryptToString("测试数据")

	// 解析格式化字符串
	algorithm, fingerprint, err := aesgcm.ParseFormattedString(formatted)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("算法: %s\n", algorithm)
	// 只验证指纹长度，不验证具体值
	if len(fingerprint) == 16 {
		fmt.Println("密钥指纹长度正确")
	}

	// 验证密钥匹配
	isMatch := crypto.IsEncryptedByMe(formatted)
	fmt.Printf("密钥匹配: %v\n", isMatch)
	// Output:
	// 算法: AES_256_GCM_V1
	// 密钥指纹长度正确
	// 密钥匹配: true
}

// Example_differentKeySizesFormatted 展示不同密钥长度的格式化字符串
func Example_differentKeySizesFormatted() {
	// AES-128
	key128, _ := aesgcm.Generate128Key()
	crypto128, _ := aesgcm.NewAESGCM(key128)
	formatted128, _ := crypto128.EncryptToString("AES-128加密数据")

	// AES-256
	key256, _ := aesgcm.Generate256Key()
	crypto256, _ := aesgcm.NewAESGCM(key256)
	formatted256, _ := crypto256.EncryptToString("AES-256加密数据")

	// 解析算法信息
	algo128, _, _ := aesgcm.ParseFormattedString(formatted128)
	algo256, _, _ := aesgcm.ParseFormattedString(formatted256)

	fmt.Printf("AES-128算法: %s\n", algo128)
	fmt.Printf("AES-256算法: %s\n", algo256)

	// 验证解密
	decrypted128, _ := crypto128.DecryptFromString(formatted128)
	decrypted256, _ := crypto256.DecryptFromString(formatted256)

	fmt.Printf("AES-128解密成功: %v\n", decrypted128 == "AES-128加密数据")
	fmt.Printf("AES-256解密成功: %v\n", decrypted256 == "AES-256加密数据")
	// Output:
	// AES-128算法: AES_128_GCM_V1
	// AES-256算法: AES_256_GCM_V1
	// AES-128解密成功: true
	// AES-256解密成功: true
}

// Example_keyFingerprintVerification 展示密钥指纹验证
func Example_keyFingerprintVerification() {
	// 创建两个不同的密钥
	key1, _ := aesgcm.Generate256Key()
	crypto1, _ := aesgcm.NewAESGCM(key1)

	key2, _ := aesgcm.Generate256Key()
	crypto2, _ := aesgcm.NewAESGCM(key2)

	// 使用密钥1加密
	formatted, _ := crypto1.EncryptToString("敏感数据")

	// 验证密钥匹配
	match1 := crypto1.IsEncryptedByMe(formatted)
	match2 := crypto2.IsEncryptedByMe(formatted)

	fmt.Printf("密钥1匹配: %v\n", match1)
	fmt.Printf("密钥2匹配: %v\n", match2)

	// 尝试使用错误密钥解密
	_, err := crypto2.DecryptFromString(formatted)
	if err != nil {
		fmt.Printf("密钥不匹配错误: %v\n", err)
	}
	// Output:
	// 密钥1匹配: true
	// 密钥2匹配: false
	// 密钥不匹配错误: key fingerprint mismatch
}

// Example_formattedStringErrorHandling 展示格式化字符串的错误处理
func Example_formattedStringErrorHandling() {
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)

	// 测试无效格式
	_, err := crypto.DecryptFromString("invalid-format")
	if err != nil {
		fmt.Printf("无效格式错误: %v\n", err)
	}

	// 测试错误的算法
	_, err = crypto.DecryptFromString("WRONG_ALGO:fingerprint:data")
	if err != nil {
		fmt.Printf("错误算法错误: %v\n", err)
	}

	// 测试错误的密钥指纹 - 创建一个有效的格式但使用错误的密钥长度标识
	validFormatted, _ := crypto.EncryptToString("测试数据")
	parts := strings.Split(validFormatted, ":")
	// 修改算法部分为错误的密钥长度
	wrongAlgorithm := "AES_128_GCM_V1"
	wrongFormatted := wrongAlgorithm + ":" + parts[1] + ":" + parts[2]
	_, err = crypto.DecryptFromString(wrongFormatted)
	if err != nil {
		fmt.Printf("密钥长度不匹配错误: %v\n", err)
	}
	// Output:
	// 无效格式错误: invalid formatted string: expected 3 parts
	// 错误算法错误: invalid algorithm format
	// 密钥长度不匹配错误: key size mismatch: expected 256, got AES_128_GCM_V1
}

// Example_formattedStringUsage 展示格式化字符串的实际使用场景
func Example_formattedStringUsage() {
	key, _ := aesgcm.Generate256Key()
	crypto, _ := aesgcm.NewAESGCM(key)

	// 场景1: 数据库存储
	userData := "用户敏感信息"
	dbFormatted, _ := crypto.EncryptToString(userData)
	// 只验证格式前缀
	if strings.HasPrefix(dbFormatted, "AES_256_GCM_V1:") {
		fmt.Println("数据库存储格式正确")
	}

	// 场景2: 配置文件
	configData := "数据库密码"
	configFormatted, _ := crypto.EncryptToString(configData)
	if strings.HasPrefix(configFormatted, "AES_256_GCM_V1:") {
		fmt.Println("配置文件格式正确")
	}

	// 场景3: API传输
	apiData := "API密钥"
	apiFormatted, _ := crypto.EncryptToString(apiData)
	if strings.HasPrefix(apiFormatted, "AES_256_GCM_V1:") {
		fmt.Println("API传输格式正确")
	}

	// 验证所有数据都能正确解密
	dbDecrypted, _ := crypto.DecryptFromString(dbFormatted)
	configDecrypted, _ := crypto.DecryptFromString(configFormatted)
	apiDecrypted, _ := crypto.DecryptFromString(apiFormatted)

	allSuccess := dbDecrypted == userData &&
		configDecrypted == configData &&
		apiDecrypted == apiData

	fmt.Printf("所有场景加解密验证: %v\n", allSuccess)
	// Output:
	// 数据库存储格式正确
	// 配置文件格式正确
	// API传输格式正确
	// 所有场景加解密验证: true
}

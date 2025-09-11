package trace_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	// 修改导入路径
)

func TestGenerateTraceparentLikeJS(t *testing.T) {
	// 仿照JavaScript版本实现
	traceparent := generateTraceparent()
	t.Logf("Generated traceparent: %s", traceparent)
}

func generateTraceparent() string {
	const version = "00"

	// 生成32字符的随机十六进制作为Trace ID
	traceId := generateValidRandomHex(32)

	// 生成16字符的随机十六进制作为Span ID
	spanId := generateValidRandomHex(16)

	// 01=已采样，00=未采样
	flags := "01"

	return fmt.Sprintf("%s-%s-%s-%s", version, traceId, spanId, flags)
}

// generateValidRandomHex 生成指定长度的随机十六进制字符串，确保不全为0
func generateValidRandomHex(length int) string {
	const hexChars = "0123456789abcdef"
	result := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 生成随机十六进制字符串
	for i := range result {
		result[i] = hexChars[r.Intn(len(hexChars))]
	}

	// 确保不全为0（至少有一个非0字符）
	allZeros := true
	for _, b := range result {
		if b != '0' {
			allZeros = false
			break
		}
	}

	// 如果全为0，则将第一个字符设为1
	if allZeros {
		result[0] = '1'
	}

	return string(result)
}

package cbc_test

import (
	"fmt"
	"testing"

	"github.com/infraboard/mcube/v2/crypto/cbc"
)

func TestPKCS7Padding(t *testing.T) {
	data := []byte("Hello World!")
	blockSize := 16 // PKCS#7 适用于 1 到 255 字节的块，这里以 16 字节为例

	// 填充
	paddedData := cbc.PKCS7Padding(data, blockSize)
	fmt.Printf("Padded Data (PKCS#7): %x\n", paddedData)

	// 去填充
	unpaddedData := cbc.PKCS7Unpadding(paddedData)
	fmt.Printf("Unpadded Data: %x, %s\n", unpaddedData, unpaddedData)
}

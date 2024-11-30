/*CBC加密 按照golang标准库的例子代码
不过里面没有填充的部分,所以补上
*/

package cbc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

func NewAESCBCCihper(key []byte) *AESCBCCihper {
	return &AESCBCCihper{
		key:    key,
		encode: DATA_ENCODE_TYPE_HEX,
	}
}

type AESCBCCihper struct {
	key    []byte
	encode DATA_ENCODE_TYPE
}

func (c *AESCBCCihper) EncryptToString(rawData []byte) (string, error) {
	cipherData, err := c.Encrypt(rawData)
	if err != nil {
		return "", err
	}
	switch c.encode {
	case DATA_ENCODE_TYPE_BASE64:
		return base64.StdEncoding.EncodeToString(cipherData), nil
	default:
		return hex.EncodeToString(cipherData), nil
	}
}

func (c *AESCBCCihper) DecryptFromString(cipherText string) (string, error) {
	var (
		cipherData []byte
		err        error
	)
	switch c.encode {
	case DATA_ENCODE_TYPE_BASE64:
		cipherData, err = base64.StdEncoding.DecodeString(cipherText)
	default:
		cipherData, err = hex.DecodeString(cipherText)
	}
	if err != nil {
		return "", err
	}

	planData, err := c.Decrypt(cipherData)
	if err != nil {
		return "", err
	}
	return string(planData), nil
}

func (c *AESCBCCihper) EncryptFromString(rawString string) (string, error) {
	return c.EncryptToString([]byte(rawString))
}

func (c *AESCBCCihper) Encrypt(rawData []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	//填充原文
	blockSize := block.BlockSize()
	rawData = PKCS7Padding(rawData, blockSize)
	//初始向量IV必须是唯一，但不需要保密
	cipherText := make([]byte, blockSize+len(rawData))
	//block大小 16
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	//block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], rawData)

	return cipherText, nil
}

func (c *AESCBCCihper) Decrypt(encryptData []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()

	if len(encryptData) < blockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := encryptData[:blockSize]
	encryptData = encryptData[blockSize:]

	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encryptData, encryptData)
	//解填充
	encryptData = PKCS7Unpadding(encryptData)
	return encryptData, nil
}

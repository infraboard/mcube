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

func MustNewAESCBCCihper(key []byte) *AESCBCCihper {
	c, err := NewAESCBCCihper(key)
	if err != nil {
		panic(err)
	}
	return c
}

func NewAESCBCCihper(key []byte) (*AESCBCCihper, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &AESCBCCihper{
		encode: DATA_ENCODE_TYPE_HEX,
		block:  block,
	}, nil
}

type AESCBCCihper struct {
	encode DATA_ENCODE_TYPE
	block  cipher.Block
}

func (c *AESCBCCihper) SetDataEncodeType(v DATA_ENCODE_TYPE) *AESCBCCihper {
	c.encode = v
	return c
}

func (c *AESCBCCihper) EncryptFromString(rawString string) (string, error) {
	return c.EncryptToString(rawString)
}

func (c *AESCBCCihper) EncryptToString(rawData string) (string, error) {
	cipherData, err := c.Encrypt([]byte(rawData))
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

func (c *AESCBCCihper) DecryptFromCipherText(cipherText string) (string, error) {
	cipherData, err := c.DecodeCipherText(cipherText)
	if err != nil {
		return "", err
	}

	planData, err := c.Decrypt(cipherData)
	if err != nil {
		return "", err
	}
	return string(planData), nil
}

func (c *AESCBCCihper) DecodeCipherText(cipherText string) (cipherData []byte, err error) {
	switch c.encode {
	case DATA_ENCODE_TYPE_BASE64:
		cipherData, err = base64.StdEncoding.DecodeString(cipherText)
	default:
		cipherData, err = hex.DecodeString(cipherText)
	}
	return
}

// 初始向量IV必须是唯一，但不需要保密, 这里随机生成一个block 作为IV向量
func (c *AESCBCCihper) GenRandomIv(cipherText []byte) ([]byte, error) {
	iv := cipherText[:c.block.BlockSize()]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	return iv, nil
}

func (c *AESCBCCihper) GetIv(encryptData []byte) []byte {
	return encryptData[:c.block.BlockSize()]
}

func (c *AESCBCCihper) GetIvFromCipherText(cipherString string) ([]byte, error) {
	data, err := c.DecodeCipherText(cipherString)
	if err != nil {
		return nil, err
	}
	return data[:c.block.BlockSize()], nil
}

func (c *AESCBCCihper) MustGetIvFromCipherText(cipherString string) []byte {
	data, _ := c.GetIvFromCipherText(cipherString)
	return data
}

func (c *AESCBCCihper) Encrypt(rawData []byte) ([]byte, error) {
	//填充原文
	blockSize := c.block.BlockSize()
	rawData = PKCS7Padding(rawData, blockSize)

	cipherText := make([]byte, blockSize+len(rawData))
	iv, err := c.GenRandomIv(cipherText)
	if err != nil {
		return nil, err
	}

	//block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(c.block, iv)
	mode.CryptBlocks(cipherText[blockSize:], rawData)
	return cipherText, nil
}

func (c *AESCBCCihper) Decrypt(encryptData []byte) ([]byte, error) {
	blockSize := c.block.BlockSize()
	if len(encryptData) < blockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := c.GetIv(encryptData)
	encryptData = encryptData[blockSize:]

	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(c.block, iv)
	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encryptData, encryptData)
	//解填充
	encryptData = PKCS7Unpadding(encryptData)
	return encryptData, nil
}

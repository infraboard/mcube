/*CBC加密 按照golang标准库的例子代码
不过里面没有填充的部分,所以补上
*/

package cbc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
)

func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// aesCBCEncrypt aes加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func aesCBCEncrypt(rawData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//填充原文
	blockSize := block.BlockSize()
	rawData = pkcs7Padding(rawData, blockSize)
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

func aesCBCDecrypt(encryptData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
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
	encryptData = pkcs7UnPadding(encryptData)
	return encryptData, nil
}

// 采用hmac进行2次hash, 取32位
func sha1Hash2(key []byte) []byte {
	h := sha1.New()
	h.Write(key)
	hashData := h.Sum(nil)
	keyBuffer := bytes.NewBuffer(hashData)

	h.Reset()
	h.Write(hashData)
	keyBuffer.Write(h.Sum(nil))

	return keyBuffer.Bytes()[:32]
}

// Encrypt aes cbc加密
func Encrypt(data, key []byte) ([]byte, error) {
	return aesCBCEncrypt(data, sha1Hash2(key))
}

// Decrypt aes cbc解密
func Decrypt(data, key []byte) ([]byte, error) {
	return aesCBCDecrypt(data, sha1Hash2(key))
}

// Encrypt aes cbc加密
func EncryptToString(plan string, key []byte) (string, error) {
	data, err := aesCBCEncrypt([]byte(plan), sha1Hash2(key))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// Decrypt aes cbc解密
func DecryptFromString(encrypt string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", err
	}

	plan, err := aesCBCDecrypt(data, sha1Hash2(key))
	if err != nil {
		return "", err
	}
	return string(plan), nil
}

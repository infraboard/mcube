package cbc

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
)

// aesCBCEncrypt aes加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
type AES_KEY_LEN int

const (
	AES_KEY_LEN_16 AES_KEY_LEN = 16
	AES_KEY_LEN_24 AES_KEY_LEN = 24
	AES_KEY_LEN_32 AES_KEY_LEN = 32
)

// 采用hmac进行2次hash, 取32位, 这样对key没有长度要求
func GenKeyBySHA1Hash2(key string, keyLen AES_KEY_LEN) []byte {
	h := sha1.New()
	h.Write([]byte(key))
	hashData := h.Sum(nil)
	keyBuffer := bytes.NewBuffer(hashData)

	h.Reset()
	h.Write(hashData)
	keyBuffer.Write(h.Sum(nil))
	return keyBuffer.Bytes()[:keyLen]
}

func GenRandomKey(keyLen AES_KEY_LEN) ([]byte, error) {
	key := make([]byte, keyLen)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	return key, nil
}

func MustGenRandomKey(keyLen AES_KEY_LEN) []byte {
	v, _ := GenRandomKey(keyLen)
	return v
}

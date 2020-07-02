package ecies

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"io"

	"golang.org/x/crypto/curve25519"
)

// NewCurve25519ECDH 使用密码学家Daniel J. Bernstein的椭圆曲线算法:Curve25519来创建ECDH实例
// 因为Curve25519独立于NIST之外, 没在标准库实现, 需要单独为期实现一套接口来支持

// GenerateKey 基于curve25519椭圆曲线算法生成秘钥对(32位)
func GenerateKey() (priv, pub [32]byte, err error) {
	_, err = io.ReadFull(rand.Reader, priv[:])
	if err != nil {
		return
	}

	priv[0] &= 248
	priv[31] &= 127
	priv[31] |= 64

	curve25519.ScalarBaseMult(&pub, &priv)

	return priv, pub, nil
}

// Encrypt 加密
func Encrypt(plainText []byte, publicKey [32]byte) ([]byte, error) {
	var (
		r, R, S, kb [32]byte
	)

	if _, err := rand.Read(r[:]); err != nil {
		return nil, err
	}

	r[0] &= 248
	r[31] &= 127
	r[31] |= 64

	kb = publicKey
	curve25519.ScalarBaseMult(&R, &r)
	curve25519.ScalarMult(&S, &r, &kb)

	ke := make([]byte, 0, len(plainText))
	for {
		if len(ke) > len(plainText) {
			break
		}

		for _, b := range sha512.Sum512(S[:]) {
			ke = append(ke, b)
		}
	}

	cipherText := make([]byte, 32+len(plainText))
	copy(cipherText[:32], R[:])
	for i := 0; i < len(plainText); i++ {
		cipherText[32+i] = plainText[i] ^ ke[i]
	}

	return cipherText, nil
}

// Decrypt 解密
func Decrypt(cipherText []byte, privateKey [32]byte) ([]byte, error) {
	if len(cipherText) <= 32 {
		return nil, fmt.Errorf("not an validate cipher data, must > 32byte")
	}

	var R, S, kb [32]byte
	copy(R[:], cipherText[:32])

	kb = privateKey

	curve25519.ScalarMult(&S, &kb, &R)

	plainText := make([]byte, len(cipherText)-32)
	ke := make([]byte, 0, len(plainText))
	for {
		if len(ke) < len(plainText) {
			for _, b := range sha512.Sum512(S[:]) {
				ke = append(ke, b)
			}
		} else {
			break
		}
	}

	for i := 0; i < len(plainText); i++ {
		plainText[i] = cipherText[32+i] ^ ke[i]
	}

	return plainText, nil
}

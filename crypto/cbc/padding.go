package cbc

import "bytes"

// PKCS7 Padding 和 PKCS5 Padding 都是用于对称加密中填充数据的方式，但它们之间有一些细微的区别：
// 填充块大小：
// PKCS5：只适用于块大小为 8 字节的加密算法（如 DES）。它的填充方式是将填充字节的值设为填充的字节数。
// PKCS7：可以用于块大小为 1 到 255 字节的加密算法（如 AES，块大小为 16 字节）。它的填充方式与 PKCS#5 类似，但适用于更广泛的块大小。

// PKCS7Padding 对数据进行 PKCS#7 填充
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7Unpadding 对数据进行 PKCS#7 去填充
func PKCS7Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

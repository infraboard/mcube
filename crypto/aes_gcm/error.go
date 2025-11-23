package aesgcm

import "errors"

// 错误定义
var (
	ErrInvalidKeySize    = errors.New("invalid key size: must be 16, 24, or 32 bytes")
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
	ErrDecryptionFailed  = errors.New("decryption failed: authentication failed")
)

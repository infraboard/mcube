package ecies

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/decred/dcrd/dcrec/secp256k1/v3"
)

// Encrypt encrypts a passed message with a receiver public key, returns ciphertext or encryption error
func Encrypt(pubkey *PublicKey, msg []byte) ([]byte, error) {
	var ct bytes.Buffer

	// Generate ephemeral key
	ek, err := GenerateKey()
	if err != nil {
		return nil, err
	}

	ct.Write(ek.PublicKey.Bytes(false))

	// Derive shared secret
	ss, err := ek.Encapsulate(pubkey)
	if err != nil {
		return nil, err
	}

	// AES encryption
	block, err := aes.NewCipher(ss)
	if err != nil {
		return nil, fmt.Errorf("cannot create new aes block: %w", err)
	}

	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("cannot read random bytes for nonce: %w", err)
	}

	ct.Write(nonce)

	aesgcm, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		return nil, fmt.Errorf("cannot create aes gcm: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, msg, nil)

	tag := ciphertext[len(ciphertext)-aesgcm.NonceSize():]
	ct.Write(tag)
	ciphertext = ciphertext[:len(ciphertext)-len(tag)]
	ct.Write(ciphertext)

	return ct.Bytes(), nil
}

// Decrypt decrypts a passed message with a receiver private key, returns plaintext or decryption error
func Decrypt(privkey *PrivateKey, msg []byte) ([]byte, error) {
	// Message cannot be less than length of public key (65) + nonce (16) + tag (16)
	if len(msg) <= (1 + 32 + 32 + 16 + 16) {
		return nil, fmt.Errorf("invalid length of message")
	}

	// Ephemeral sender public key
	ethPubkey := &PublicKey{
		Curve: secp256k1.S256(),
		X:     new(big.Int).SetBytes(msg[1:33]),
		Y:     new(big.Int).SetBytes(msg[33:65]),
	}

	// Shift message
	msg = msg[65:]

	// Derive shared secret
	ss, err := ethPubkey.Decapsulate(privkey)
	if err != nil {
		return nil, err
	}

	// AES decryption part
	nonce := msg[:16]
	tag := msg[16:32]

	// Create Golang-accepted ciphertext
	ciphertext := bytes.Join([][]byte{msg[32:], tag}, nil)

	block, err := aes.NewCipher(ss)
	if err != nil {
		return nil, fmt.Errorf("cannot create new aes block: %w", err)
	}

	gcm, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		return nil, fmt.Errorf("cannot create gcm cipher: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot decrypt ciphertext: %w", err)
	}

	return plaintext, nil
}

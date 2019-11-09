package crypto_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/infraboard/mcube/crypto"
)

func TestEncrypt(t *testing.T) {
	should := require.New(t)

	privKey, pubKey, err := crypto.GenerateKey()

	should.NoError(err)

	plainText := []byte("This is message for ecies test.")

	cipherText, err := crypto.Encrypt(plainText, pubKey)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	plainText2, err := crypto.Decrypt(cipherText, privKey)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	if bytes.Compare(plainText, plainText2) != 0 {
		t.Fatalf("result did not match:\nGot:\n%s\nExpected:\n%s\n", hex.Dump(plainText2), hex.Dump(plainText))
	}
}

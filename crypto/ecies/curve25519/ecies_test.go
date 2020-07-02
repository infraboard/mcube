package ecies_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	ecies "github.com/infraboard/mcube/crypto/ecies/curve25519"
)

func TestEncrypt(t *testing.T) {
	should := require.New(t)

	privKey, pubKey, err := ecies.GenerateKey()

	should.NoError(err)

	plainText := []byte("This is message for ecies test.")

	cipherText, err := ecies.Encrypt(plainText, pubKey)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	plainText2, err := ecies.Decrypt(cipherText, privKey)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	if bytes.Compare(plainText, plainText2) != 0 {
		t.Fatalf("result did not match:\nGot:\n%s\nExpected:\n%s\n", hex.Dump(plainText2), hex.Dump(plainText))
	}
}

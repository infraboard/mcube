package cbc_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/infraboard/mcube/crypto/cbc"
)

func TestAESCBC(t *testing.T) {
	data := []byte("abcdefg")
	key := []byte("123456")

	should := require.New(t)

	cipherData, err := cbc.Encrypt(data, key)
	should.NoError(err)
	t.Logf("cipher data: %s", cipherData)

	rawData, err := cbc.Decrypt(cipherData, key)
	should.NoError(err)
	t.Logf("raw data: %s", rawData)

	should.Equal(data, []byte(rawData))
}

func TestEncryptString(t *testing.T) {
	data := "abcdefg"
	key := []byte("123456")

	should := require.New(t)

	cipherData, err := cbc.EncryptToString(data, key)
	should.NoError(err)
	t.Logf("cipher data: %s", cipherData)

	rawData, err := cbc.DecryptFromString(cipherData, key)
	should.NoError(err)
	t.Logf("raw data: %s", rawData)

	should.Equal(data, rawData)
}

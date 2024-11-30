package cbc_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/infraboard/mcube/v2/crypto/cbc"
)

var (
	data      = "abcdefg"
	key       = cbc.MustGenRandomKey(cbc.AES_KEY_LEN_32)
	cipher, _ = cbc.NewAESCBCCihper(key)
)

func TestEncryptAndDecrypt(t *testing.T) {
	should := require.New(t)

	cipherText, err := cipher.EncryptFromString(data)
	should.NoError(err)
	t.Logf("cipher text: %s", cipherText)

	t.Log(cipher.MustGetIvFromCipherText(cipherText))

	rawData, err := cipher.DecryptFromCipherText(cipherText)
	should.NoError(err)
	t.Logf("raw text: %s", rawData)
	should.Equal(data, rawData)

}

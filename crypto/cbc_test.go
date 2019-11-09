package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/infraboard/mcube/crypto"
)

func TestAESCBC(t *testing.T) {
	data := []byte("abcdefg")
	key := []byte("123456")

	should := require.New(t)

	cipherData, err := crypto.AESCBCEncrypt(data, key)
	should.NoError(err)
	t.Logf("cipher data: %s", cipherData)

	rawData, err := crypto.AESCBCDecrypt(cipherData, key)
	should.NoError(err)
	t.Logf("raw data: %s", rawData)

	should.Equal(data, []byte(rawData))
}

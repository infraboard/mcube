package cbc_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/infraboard/mcube/v2/crypto/cbc"
)

var (
	data   = "abcdefg"
	key    = cbc.MustGenRandomKey(cbc.AES_KEY_LEN_32)
	cipher = cbc.NewAESCBCCihper(key)
)

func TestAESCBC(t *testing.T) {
	should := require.New(t)

	cipherData, err := cipher.EncryptFromString(data)
	should.NoError(err)
	t.Logf("cipher data: %s", cipherData)

	rawData, err := cipher.DecryptFromString(cipherData)
	should.NoError(err)
	t.Logf("raw data: %s", rawData)
	should.Equal(data, rawData)
}

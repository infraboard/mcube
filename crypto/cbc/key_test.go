package cbc_test

import (
	"testing"

	"github.com/infraboard/mcube/v2/crypto/cbc"
)

func TestGenKeyBySHA1Hash2(t *testing.T) {
	aesCihperKey := cbc.GenKeyBySHA1Hash2("test", cbc.AES_KEY_LEN_32)
	t.Log(aesCihperKey)
}

func TestMustGenRandomKey(t *testing.T) {
	aesCihperKey := cbc.MustGenRandomKey(cbc.AES_KEY_LEN_32)
	t.Log(aesCihperKey)
}

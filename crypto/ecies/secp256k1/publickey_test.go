package ecies_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	ecies "github.com/infraboard/mcube/crypto/ecies/secp256k1"
)

func TestNewPublicKeyFromHex(t *testing.T) {
	_, err := ecies.NewPublicKeyFromHex(testingReceiverPubkeyHex)
	assert.NoError(t, err)
}

func TestPublicKey_Equals(t *testing.T) {
	privkey, err := ecies.GenerateKey()
	if !assert.NoError(t, err) {
		return
	}

	assert.True(t, privkey.PublicKey.Equals(privkey.PublicKey))
}

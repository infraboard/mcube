package ecies_test

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/hkdf"

	ecies "github.com/infraboard/mcube/crypto/ecies/secp256k1"
)

const testingMessage = "helloworld"
const testingReceiverPubkeyHex = "0498afe4f150642cd05cc9d2fa36458ce0a58567daeaf5fde7333ba9b403011140a4e28911fcf83ab1f457a30b4959efc4b9306f514a4c3711a16a80e3b47eb58b"
const testingReceiverPrivkeyHex = "95d3c5e483e9b1d4f5fc8e79b2deaf51362980de62dbb082a9a4257eef653d7d"
const pythonBackend = "https://eciespy.herokuapp.com/"

var testingReceiverPrivkey = []byte{51, 37, 145, 156, 66, 168, 189, 189, 176, 19, 177, 30, 148, 104, 25, 140, 155, 42, 248, 190, 121, 110, 16, 174, 143, 148, 72, 129, 94, 113, 219, 58}

func TestGenerateKey(t *testing.T) {
	_, err := ecies.GenerateKey()
	assert.NoError(t, err)
}

func TestEncryptAndDecrypt(t *testing.T) {
	privkey := ecies.NewPrivateKeyFromBytes(testingReceiverPrivkey)

	ciphertext, err := ecies.Encrypt(privkey.PublicKey, []byte(testingMessage))
	if !assert.NoError(t, err) {
		return
	}

	plaintext, err := ecies.Decrypt(privkey, ciphertext)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, testingMessage, string(plaintext))
}

func TestPublicKeyDecompression(t *testing.T) {
	// Generate public key
	privkey, err := ecies.GenerateKey()
	if !assert.NoError(t, err) {
		return
	}

	// Drop Y part and restore it
	pubkey, err := ecies.NewPublicKeyFromHex(privkey.PublicKey.Hex(true))
	if !assert.NoError(t, err) {
		return
	}

	// Check that point is still at curve
	assert.True(t, privkey.IsOnCurve(pubkey.X, pubkey.Y))
}

func TestKEM(t *testing.T) {
	derived := make([]byte, 32)
	kdf := hkdf.New(sha256.New, []byte("secret"), nil, nil)
	if _, err := io.ReadFull(kdf, derived); err != nil {
		if !assert.Equal(
			t,
			hex.EncodeToString(derived),
			"2f34e5ff91ec85d53ca9b543683174d0cf550b60d5f52b24c97b386cfcf6cbbf",
		) {
			return
		}
	}

	k1 := ecies.NewPrivateKeyFromBytes(new(big.Int).SetInt64(2).Bytes())
	k2 := ecies.NewPrivateKeyFromBytes(new(big.Int).SetInt64(3).Bytes())

	sk1, err := k1.Encapsulate(k2.PublicKey)
	if !assert.NoError(t, err) {
		return
	}
	sk2, err := k1.PublicKey.Decapsulate(k2)
	if !assert.NoError(t, err) {
		return
	}

	if !assert.Equal(t, sk1, sk2) {
		return
	}

	assert.Equal(
		t,
		hex.EncodeToString(sk1),
		"6f982d63e8590c9d9b5b4c1959ff80315d772edd8f60287c9361d548d5200f82",
	)
}

func TestDecryptAgainstPythonVersion(t *testing.T) {
	prv, err := ecies.NewPrivateKeyFromHex(testingReceiverPrivkeyHex)
	if !assert.NoError(t, err) {
		return
	}

	form := make(url.Values)
	form.Set("data", testingMessage)
	form.Set("pub", testingReceiverPubkeyHex)

	resp, err := http.PostForm(pythonBackend, form)
	if !assert.NoError(t, err) {
		return
	}

	defer resp.Body.Close()

	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}

	hexBytes, err := ioutil.ReadAll(resp.Body)
	if !assert.NoError(t, err) {
		return
	}

	ciphertext, err := hex.DecodeString(string(hexBytes))
	if !assert.NoError(t, err) {
		return
	}

	plaintext, err := ecies.Decrypt(prv, ciphertext)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, testingMessage, string(plaintext))
}

func TestEncryptAgainstPythonVersion(t *testing.T) {
	prv, err := ecies.NewPrivateKeyFromHex(testingReceiverPrivkeyHex)
	if !assert.NoError(t, err) {
		return
	}

	ciphertext, err := ecies.Encrypt(prv.PublicKey, []byte(testingMessage))
	if !assert.NoError(t, err) {
		return
	}

	form := make(url.Values)
	form.Set("data", hex.EncodeToString(ciphertext))
	form.Set("prv", testingReceiverPrivkeyHex)

	resp, err := http.PostForm(pythonBackend, form)
	if !assert.NoError(t, err) {
		return
	}

	defer resp.Body.Close()

	if !assert.Equal(t, http.StatusOK, resp.StatusCode) {
		return
	}

	plaintext, err := ioutil.ReadAll(resp.Body)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, string(plaintext), testingMessage)
}

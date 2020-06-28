// Package kafka https://github.com/Shopify/sarama/blob/master/examples/sasl_scram_client/scram_client.go
package kafka

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"

	"github.com/xdg/scram"
)

var (
	// SHA256 todo
	SHA256 scram.HashGeneratorFcn = func() hash.Hash { return sha256.New() }
	// SHA512 todo
	SHA512 scram.HashGeneratorFcn = func() hash.Hash { return sha512.New() }
)

// XDGSCRAMClient todo
type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

// Begin todo
func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

// Step todo
func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

// Done todo
func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}

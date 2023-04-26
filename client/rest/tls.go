package rest

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

func NewTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	ca := x509.NewCertPool()
	file, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	ca.AppendCertsFromPEM(file)

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		RootCAs:            ca,
		ClientCAs:          ca,
	}, nil
}

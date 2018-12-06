package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

//TLSCredentialFromFile loads certificate from file end returns transport credentials
func TLSCredentialFromFile(cacert, svcert, key string, mutual bool) (credentials.TransportCredentials, error) {

	cert, err := tls.LoadX509KeyPair(svcert, key)
	if err != nil {
		return nil, err
	}

	rawCaCert, err := ioutil.ReadFile(cacert)
	if err != nil {
		return nil, err
	}

	return tlsCredential(rawCaCert, cert, mutual), nil
}

//TLSCredentialFromData returns transport credentials
func TLSCredentialFromData(cacert, svcert, key []byte, mutual bool) (credentials.TransportCredentials, error) {

	cert, err := tls.X509KeyPair(svcert, key)

	if err != nil {
		return nil, err
	}

	return tlsCredential(cacert, cert, mutual), nil
}

func tlsCredential(cacert []byte, cert tls.Certificate, mutual bool) credentials.TransportCredentials {
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cacert)

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		RootCAs:      caCertPool,
	}

	if mutual {
		tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return credentials.NewTLS(tlsCfg)
}

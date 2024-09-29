package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"os"
)

// NewServerTLS creates a new TLS configuration for a server using the provided certificate, key, and CA certificate file paths
func NewServerTLS(cert, key, ca string) (*tls.Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not get working directory: %w", err)
	}

	serverCert, err := tls.LoadX509KeyPair(wd+cert, wd+key)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate and key: %w", err)
	}

	file, err := os.Open(wd + ca)
	if err != nil {
		return nil, fmt.Errorf("failed to open ca: %w", err)
	}
	defer file.Close()

	caCert, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read ca: %w", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return tlsConfig, nil
}

// NewClientTLS creates a new TLS configuration for a client using the provided certificate, key, and CA certificate file paths
func NewClientTLS(cert, key, ca string) (*tls.Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not get working directory: %w", err)
	}

	clientCert, err := tls.LoadX509KeyPair(wd+cert, wd+key)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate and key: %w", err)
	}

	file, err := os.Open(wd + ca)
	if err != nil {
		return nil, fmt.Errorf("failed to open ca: %w", err)
	}
	defer file.Close()

	caCert, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read ca: %w", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	}

	return tlsConfig, nil
}

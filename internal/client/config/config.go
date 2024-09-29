package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration settings, including logging level and client certificates.
type Config struct {
	EnvLogLevel string // Log level for the application
	GRPCServer  string // Address of the gRPC server
	ClientCert  string // Path to the client certificate file
	ClientKey   string // Path to the client private key file
	ClientCa    string // Path to the client CA certificate file
}

// New loads the configuration from the "client.env" file using environment variables
// and returns an instance of Config. If the environment file cannot be loaded, it returns an error.
func New() (*Config, error) {
	err := godotenv.Load("client.env")
	if err != nil {
		return nil, fmt.Errorf("new load .env: %w", err)
	}

	// Initialize Config and assign environment variables
	config := &Config{
		EnvLogLevel: os.Getenv("LOG_LEVEL"),
		GRPCServer:  os.Getenv("GRPC_SERVER"),
		ClientCert:  os.Getenv("CLIENT_CERT_FILE"),
		ClientKey:   os.Getenv("CLIENT_KEY_FILE"),
		ClientCa:    os.Getenv("CLIENT_CA_FILE"),
	}
	return config, nil
}

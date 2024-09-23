package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	EnvLogLevel string
	GRPCServer  string
	ClientCert  string
	ClientKey   string
	ClientCa    string
}

func New() (*Config, error) {
	err := godotenv.Load("client.env")
	if err != nil {
		return nil, fmt.Errorf("new load .env: %w", err)
	}

	config := &Config{}

	config.EnvLogLevel = os.Getenv("LOG_LEVEL")
	config.GRPCServer = os.Getenv("GRPC_SERVER")
	config.ClientCert = os.Getenv("CLIENT_CERT_FILE")
	config.ClientKey = os.Getenv("CLIENT_KEY_FILE")
	config.ClientCa = os.Getenv("CLIENT_CA_FILE")

	return config, nil
}

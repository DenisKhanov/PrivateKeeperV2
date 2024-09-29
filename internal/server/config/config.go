package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the application configuration parameters.
// Each field corresponds to an expected environment variable.
type Config struct {
	EnvLogLevel     string // Log level for the application (e.g., DEBUG, INFO)
	DatabaseURI     string // Connection string for the database
	GRPCServer      string // Address of the gRPC server
	TokenName       string // Name of the authentication token
	TokenSecret     string // Secret key for signing tokens
	TokenExpHours   int    // Token expiration time in hours
	ServerCert      string // Path to the server's SSL certificate
	ServerKey       string // Path to the server's SSL key
	ServerCa        string // Path to the server's CA file
	RedisURL        string // URL of the Redis server
	RedisPassword   string // Password for the Redis server
	RedisDB         int    // Redis database number
	RedisTimeoutSec int    // Timeout for Redis operations in seconds
}

// New initializes a new Config instance by loading environment variables from a .env file.
// It returns a pointer to the Config struct and an error if any of the environment variables are missing or invalid.
func New() (*Config, error) {
	err := godotenv.Load("server.env")
	if err != nil {
		return nil, fmt.Errorf("new load .env: %w", err)
	}

	config := &Config{}
	config.EnvLogLevel = os.Getenv("LOG_LEVEL")
	config.DatabaseURI = os.Getenv("DATABASE_URI")
	config.GRPCServer = os.Getenv("GRPC_SERVER")
	config.TokenName = os.Getenv("TOKEN_NAME")
	expHours, err := strconv.Atoi(os.Getenv("TOKEN_EXP_HOURS"))
	if err != nil {
		return nil, fmt.Errorf("atoi TOKEN_EXP_HOURS: %w", err)
	}
	config.TokenExpHours = expHours
	config.TokenSecret = os.Getenv("TOKEN_SECRET")
	config.ServerCert = os.Getenv("SERVER_CERT_FILE")
	config.ServerKey = os.Getenv("SERVER_KEY_FILE")
	config.ServerCa = os.Getenv("SERVER_CA_FILE")

	config.RedisURL = os.Getenv("REDIS_URL")
	config.RedisPassword = os.Getenv("REDIS_PASSWORD")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, fmt.Errorf("atoi REDIS_DB: %w", err)
	}
	config.RedisDB = db

	config.RedisTimeoutSec, err = strconv.Atoi(os.Getenv("REDIS_TIMEOUT_SEC"))
	if err != nil {
		return nil, fmt.Errorf("atoi REDIS_TIMEOUT_SEC: %w", err)
	}

	return config, nil
}

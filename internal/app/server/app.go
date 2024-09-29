package server

import (
	"context"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/logcfg"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	tlsCreds "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/proto/binary_data"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credentials"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credit_card"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/proto/text_data"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/proto/user"
	binaryDataGRPCHandlers "github.com/DenisKhanov/PrivateKeeperV2/internal/server/binary_data/api/v1/grpchandlers"
	binaryDataValidation "github.com/DenisKhanov/PrivateKeeperV2/internal/server/binary_data/api/v1/validation"
	binaryDataService "github.com/DenisKhanov/PrivateKeeperV2/internal/server/binary_data/service"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/cache"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/config"
	credentialsGRPCHandlers "github.com/DenisKhanov/PrivateKeeperV2/internal/server/credentials/api/v1/grpchandlers"
	credentialsValidation "github.com/DenisKhanov/PrivateKeeperV2/internal/server/credentials/api/v1/validation"
	credentialsService "github.com/DenisKhanov/PrivateKeeperV2/internal/server/credentials/service"
	creditCardGRPCHandlers "github.com/DenisKhanov/PrivateKeeperV2/internal/server/credit_card/api/v1/grpchandlers"
	creditCardValidation "github.com/DenisKhanov/PrivateKeeperV2/internal/server/credit_card/api/v1/validation"
	creditCardService "github.com/DenisKhanov/PrivateKeeperV2/internal/server/credit_card/service"
	repository "github.com/DenisKhanov/PrivateKeeperV2/internal/server/data_repository"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/encryption"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/auth"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/keyextraction"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/storage/postgresql"
	textDataGRPCHandlers "github.com/DenisKhanov/PrivateKeeperV2/internal/server/text_data/api/v1/grpchandlers"
	textDataValidation "github.com/DenisKhanov/PrivateKeeperV2/internal/server/text_data/api/v1/validation"
	textDataService "github.com/DenisKhanov/PrivateKeeperV2/internal/server/text_data/service"
	userGRPCHandlers "github.com/DenisKhanov/PrivateKeeperV2/internal/server/user/api/v1/grpchandlers"
	userValidation "github.com/DenisKhanov/PrivateKeeperV2/internal/server/user/api/v1/validation"
	userRepository "github.com/DenisKhanov/PrivateKeeperV2/internal/server/user/repository"
	userService "github.com/DenisKhanov/PrivateKeeperV2/internal/server/user/service"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/tlsconfig"
	"github.com/DenisKhanov/PrivateKeeperV2/pkg/jwtmanager"
)

// Run initializes and starts the gRPC server for the application.
//
// This function performs the following steps:
// - Loads the configuration for the server.
// - Configures logging for the server to a log file.
// - Initializes Redis for caching, cryptographic services for encryption, and Postgres for database operations.
// - Sets up JWT authentication for securing API requests.
// - Initializes various service components including user, credit card, text data, credentials, and binary data services.
// - Creates validators for input data for each service.
// - Configures and starts the gRPC server with TLS encryption and authentication middleware.
// - Registers the gRPC services (user, credit card, text data, credentials, binary data) with the server.
// - Sets up a TCP listener and serves the gRPC server, blocking until an error occurs or the server shuts down.
func Run() {

	cfg, err := config.New()
	if err != nil {
		log.Println("Failed to initialize config", err.Error())
		os.Exit(1)
	}
	logFileName := "keeperServer.log"
	logcfg.RunLoggerConfig(cfg.EnvLogLevel, logFileName)

	redis, err := cache.NewRedis(cfg.RedisURL, cfg.RedisPassword, cfg.RedisDB, cfg.RedisTimeoutSec)
	if err != nil {
		logrus.WithError(err).Error("Failed to initialize redis")
		os.Exit(1)
	}

	cryptService, err := encryption.New([]byte("master-key"))
	if err != nil {
		logrus.WithError(err).Error("Failed to initialize crypt service")
		os.Exit(1)
	}

	postgresPool, err := initPostgresPool(cfg.DatabaseURI)
	if err != nil {
		logrus.WithError(err).Error("Failed to initialize postgres pool")
		os.Exit(1)
	}

	jwtManager := jwtmanager.New(cfg.TokenName, cfg.TokenSecret, cfg.TokenExpHours)

	userRepo := userRepository.New(postgresPool)
	dataRepo := repository.New(postgresPool)

	userServ := userService.New(userRepo, cryptService, jwtManager, redis)
	creditCardServ := creditCardService.New(dataRepo, cryptService, jwtManager)
	textDataServ := textDataService.New(dataRepo, cryptService, jwtManager)
	credentialServ := credentialsService.New(dataRepo, cryptService, jwtManager)
	binaryDataServ := binaryDataService.New(dataRepo, cryptService, jwtManager)

	tls, err := tlsconfig.NewServerTLS(cfg.ServerCert, cfg.ServerKey, cfg.ServerCa)
	if err != nil {
		logrus.WithError(err).Error("Failed to initialize tls")
		os.Exit(1)
	}

	validate := validator.New()
	creditCardValidator, err := creditCardValidation.New(validate)
	if err != nil {
		logrus.WithError(err).Error("Failed to initialize credit card validator")
		os.Exit(1)
	}
	textDataValidator, err := textDataValidation.New(validate)
	if err != nil {
		logrus.WithError(err).Error("Failed to initialize text data validator")
		os.Exit(1)
	}
	credentialsValidator, err := credentialsValidation.New(validate)
	if err != nil {
		logrus.WithError(err).Error("Failed to initialize credentials validator")
		os.Exit(1)
	}
	binaryDataValidator, err := binaryDataValidation.New(validate)
	if err != nil {
		logrus.WithError(err).Error("Failed to initialize binary data validator")
		os.Exit(1)
	}

	jwtAuth := auth.New(jwtManager)
	userKeyExtractor := keyextraction.New(cryptService, userRepo, redis)

	grpcServer := grpc.NewServer(grpc.Creds(tlsCreds.NewTLS(tls)),
		grpc.ChainUnaryInterceptor(jwtAuth.GRPCJWTAuth, userKeyExtractor.ExtractUserKey))

	user.RegisterUserServiceServer(grpcServer, userGRPCHandlers.New(userServ, userValidation.New(validate)))
	credit_card.RegisterCreditCardServiceServer(grpcServer, creditCardGRPCHandlers.New(creditCardServ, creditCardValidator))
	text_data.RegisterTextDataServiceServer(grpcServer, textDataGRPCHandlers.New(textDataServ, textDataValidator))
	credentials.RegisterCredentialsServiceServer(grpcServer, credentialsGRPCHandlers.New(credentialServ, credentialsValidator))
	binary_data.RegisterBinaryDataServiceServer(grpcServer, binaryDataGRPCHandlers.New(binaryDataServ, binaryDataValidator))

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", cfg.GRPCServer)
	if err != nil {
		logrus.WithError(err).Error("Unable to create listener")
		os.Exit(1)
	}

	if err = grpcServer.Serve(listener); err != nil {
		logrus.WithError(err).Error("Unable to start gRPC server")
		os.Exit(1)
	}
}

// initPostgresPool initializes a connection to the PostgreSQL database using the provided URI.
// It also applies any pending database migrations. If the connection or migrations fail,
// an error is returned.
func initPostgresPool(databaseURI string) (*postgresql.PostgresPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	postgresPool, err := postgresql.NewPool(ctx, databaseURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	migrations, err := postgresql.NewMigrations(postgresPool)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrations: %w", err)
	}

	err = migrations.Up()
	if err != nil {
		return nil, fmt.Errorf("failed to up migrations: %w", err)
	}
	logrus.Info("Connected to database")

	return postgresPool, nil
}

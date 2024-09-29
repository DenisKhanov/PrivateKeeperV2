package client

import (
	"bufio"
	"context"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/logcfg"
	"github.com/sirupsen/logrus"
	"log"
	"os"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	binarypb "github.com/DenisKhanov/PrivateKeeperV2/internal/client/binary_data/pbclient"
	binaryservice "github.com/DenisKhanov/PrivateKeeperV2/internal/client/binary_data/service"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/config"
	credentialspb "github.com/DenisKhanov/PrivateKeeperV2/internal/client/credentials/pbclient"
	credentialsservice "github.com/DenisKhanov/PrivateKeeperV2/internal/client/credentials/service"
	creditcardpb "github.com/DenisKhanov/PrivateKeeperV2/internal/client/credit_card/pbclient"
	creditcardservice "github.com/DenisKhanov/PrivateKeeperV2/internal/client/credit_card/service"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/state"
	textdatapb "github.com/DenisKhanov/PrivateKeeperV2/internal/client/text_data/pbclient"
	textdataservice "github.com/DenisKhanov/PrivateKeeperV2/internal/client/text_data/service"
	userpb "github.com/DenisKhanov/PrivateKeeperV2/internal/client/user/pbclient"
	userservice "github.com/DenisKhanov/PrivateKeeperV2/internal/client/user/service"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/proto/binary_data"
	credGrpc "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credentials"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credit_card"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/proto/text_data"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/proto/user"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/tlsconfig"
)

// Run initializes the client application, configures necessary services,
// and provides an interactive command-line interface for the user.
//
// The function performs the following steps:
// - Loads the configuration for the application.
// - Configures logging to a specified log file.
// - Initializes TLS for secure gRPC communication with the server.
// - Establishes a gRPC client connection for communicating with various services.
// - Sets up client-side state management and initializes service clients (user, credit card, text data, credentials, and binary data).
// - Enters an interactive loop where the user can issue commands to perform various actions such as login, register, save, and load data.
//
// The user can interact with the application via the console input where different numbered options correspond to different functionalities.
// The application will continuously run until the user enters the quit command.
func Run() {

	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		log.Println("Failed to initialize config", err.Error())
		os.Exit(1)
	}

	logFileName := "keeperClient.log"
	logcfg.RunLoggerConfig(cfg.EnvLogLevel, logFileName)

	tls, err := tlsconfig.NewClientTLS(cfg.ClientCert, cfg.ClientKey, cfg.ClientCa)
	if err != nil {
		logrus.WithError(err).Error("Failed to initialize tls")
		os.Exit(1)
	}

	grpcClient, err := grpc.NewClient(cfg.GRPCServer, grpc.WithTransportCredentials(credentials.NewTLS(tls)))
	if err != nil {
		logrus.WithError(err).Error("Failed to initialize grpcClient")
		os.Exit(1)
	}

	clientState := state.NewClientState()

	userClient := userpb.NewUserPBClient(user.NewUserServiceClient(grpcClient))
	userService := userservice.NewUserService(userClient, clientState)

	creditCardClient := creditcardpb.NewCreditCardPBClient(credit_card.NewCreditCardServiceClient(grpcClient))
	creditCardService := creditcardservice.NewUserService(creditCardClient, clientState)

	textDataClient := textdatapb.NewTextDataPBClient(text_data.NewTextDataServiceClient(grpcClient))
	textDataService := textdataservice.NewTextDataService(textDataClient, clientState)

	credentialsClient := credentialspb.NewCredentialsPBClient(credGrpc.NewCredentialsServiceClient(grpcClient))
	credentialsService := credentialsservice.NewCredentialsService(credentialsClient, clientState)

	binaryClient := binarypb.NewBinaryDataPBClient(binary_data.NewBinaryDataServiceClient(grpcClient))
	binaryService := binaryservice.NewBinaryDataService(binaryClient, clientState)

	scanner := bufio.NewScanner(os.Stdin)

	blue := color.New(color.FgBlue).SprintFunc()

	for {
		if clientState.IsAuthorized() {
			fmt.Printf("You are authorized as %s\n", blue(clientState.GetLogin()))
		} else {
			fmt.Printf("You are not authorized, please login or register\n")
		}

		if clientState.GetDirPath() == "" {
			fmt.Printf("Working directory is not set \n")
		} else {
			fmt.Printf("Working directory is set to %s\n", blue(clientState.GetDirPath()))
		}
		//TODO реализовать подпункты
		fmt.Println("Input command number to proceed")
		fmt.Println("[1] - login")
		fmt.Println("[2] - register")
		fmt.Println("---------------------------------------------")
		fmt.Println("[3] - save credit card")
		fmt.Println("[4] - load all credit cards information")
		fmt.Println("[5] - load credit card data")
		fmt.Println("---------------------------------------------")
		fmt.Println("[6] - save text data")
		fmt.Println("[7] - load all text data information")
		fmt.Println("[8] - load text data")
		fmt.Println("---------------------------------------------")
		fmt.Println("[9] - save credentials")
		fmt.Println("[10] - load all credentials information")
		fmt.Println("[11] - load credentials data")
		fmt.Println("---------------------------------------------")
		fmt.Println("[12] - save binary file")
		fmt.Println("[13] - load all binary files information")
		fmt.Println("[14] - load binary file")
		fmt.Println("---------------------------------------------")
		fmt.Println("[15] - set working directory")
		fmt.Println("------------")
		fmt.Println("|[0] - quit|")
		fmt.Println("------------")
		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "1":
			userService.LoginUser(ctx)
		case "2":
			userService.RegisterUser(ctx)
		case "3":
			creditCardService.Save(ctx)
		case "4":
			creditCardService.LoadAllInfo(ctx)
		case "5":
			creditCardService.LoadData(ctx)
		case "6":
			textDataService.Save(ctx)
		case "7":
			textDataService.LoadAllInfo(ctx)
		case "8":
			textDataService.LoadData(ctx)
		case "9":
			credentialsService.Save(ctx)
		case "10":
			credentialsService.LoadAllInfo(ctx)
		case "11":
			credentialsService.LoadData(ctx)
		case "12":
			binaryService.Save(ctx)
		case "13":
			binaryService.LoadAllInfo(ctx)
		case "14":
			binaryService.LoadData(ctx)
		case "15":
			clientState.SetWorkingDirectory()
		case "0":
			fmt.Println("Application shutdown.")
			return
		default:
			fmt.Println("Unknown command, please try again")
		}
	}
}

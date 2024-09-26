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

	textDataClient := textdatapb.NewCreditCardPBClient(text_data.NewTextDataServiceClient(grpcClient))
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

		fmt.Println("Input command number to proceed")
		fmt.Println("[1] - login")
		fmt.Println("[2] - register")
		fmt.Println("[3] - save credit card")
		fmt.Println("[4] - load credit cards")
		fmt.Println("[5] - save text data")
		fmt.Println("[6] - load text data")
		fmt.Println("[7] - save credentials")
		fmt.Println("[8] - load credentials")
		fmt.Println("[9] - save binary file")
		fmt.Println("[10] - load binary files")
		fmt.Println("[11] - set working directory")
		fmt.Println("[0] - quit")
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
			creditCardService.Load(ctx)
		case "5":
			textDataService.Save(ctx)
		case "6":
			textDataService.Load(ctx)
		case "7":
			credentialsService.Save(ctx)
		case "8":
			credentialsService.Load(ctx)
		case "9":
			binaryService.Save(ctx)
		case "10":
			binaryService.Load(ctx)
		case "11":
			clientState.SetWorkingDirectory()
		case "0":
			fmt.Println("Application shutdown.")
			return
		default:
			fmt.Println("Unknown command, please try again")
		}
	}
}

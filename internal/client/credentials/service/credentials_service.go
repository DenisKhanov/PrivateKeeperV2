package service

import (
	"bufio"
	"context"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/lib"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/state"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

// CredentialsService defines the interface for interacting with credentials-related operations.
// It includes methods to save, load, and retrieve information about user credentials.
type CredentialsService interface {
	SaveCredentials(ctx context.Context, token string, cred model.CredentialsPostRequest) (model.Credentials, error)
	LoadCredentialsData(ctx context.Context, token string, dataID string) (model.Credentials, error)
	LoadAllCredentialsDataInfo(ctx context.Context, token string) ([]model.DataInfo, error)
}

// CredentialsProvider provides methods for managing user credentials.
// It holds a reference to a CredentialsService and maintains the client's state.
type CredentialsProvider struct {
	credentialsService CredentialsService // Service to handle credentials operations
	state              *state.ClientState // Client's state, including authorization and directory information
}

// NewCredentialsService initializes a new CredentialsProvider with the given CredentialsService
// and ClientState. It returns a pointer to the newly created CredentialsProvider.
func NewCredentialsService(u CredentialsService, state *state.ClientState) *CredentialsProvider {
	return &CredentialsProvider{
		credentialsService: u,
		state:              state,
	}
}

// Save prompts the user for credentials data (login, password, metadata) and saves it
// using the credentialsService. It checks if the user is authorized before proceeding.
func (p *CredentialsProvider) Save(ctx context.Context) {
	red := color.New(color.FgRed).SprintFunc()

	if !p.state.IsAuthorized() {
		fmt.Println(red("You are not authorized, please use 'login' or 'register'"))
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	cyanBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	req := model.CredentialsPostRequest{}
	fmt.Println(cyanBold("Input credentials data 'login password metadata':"))

	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("Input login as %s: ", yellow("'text'"))
	scanner.Scan()
	req.Login = scanner.Text()

	fmt.Printf("Input password as %s: ", yellow("'text'"))
	scanner.Scan()
	req.Password = scanner.Text()

	fmt.Printf("Input metadata as %s: ", yellow("'text'"))
	scanner.Scan()
	req.MetaData = scanner.Text()

	_, err := p.credentialsService.SaveCredentials(ctx, p.state.GetToken(), req)
	if err != nil {
		lib.UnpackGRPCError(err)
		return
	}

	fmt.Println(color.New(color.FgGreen).SprintFunc()("Credentials successfully saved"))
}

// LoadAllInfo retrieves and displays information about all saved credentials.
// It checks for user authorization and a valid working directory before loading the data.
func (p *CredentialsProvider) LoadAllInfo(ctx context.Context) {
	red := color.New(color.FgRed).SprintFunc()

	if !p.state.IsAuthorized() {
		fmt.Println(red("You are not authorized, please use 'login' or 'register'"))
		return
	}

	if p.state.GetDirPath() == "" {
		fmt.Println(red("To proceed you must set working directory"))
		return
	}

	credentialsDataInfo, err := p.credentialsService.LoadAllCredentialsDataInfo(ctx, p.state.GetToken())
	if err != nil {
		logrus.WithError(err).Error("All user data info load failed")
		fmt.Println(red("All credentials data info load failed"), "please try again")
		lib.UnpackGRPCError(err)
		return
	}

	green := color.New(color.FgGreen).SprintFunc()
	fmt.Println(green("-------------------------------------"))

	var sb strings.Builder
	for _, dataInfo := range credentialsDataInfo {
		sb.WriteString("Data ID: " + dataInfo.ID + "\n")
		sb.WriteString("Data type: " + dataInfo.DataType + "\n")
		sb.WriteString("Metadata : " + dataInfo.MetaData + "\n")
		sb.WriteString("Created at: : " + dataInfo.CreatedAt + "\n")
		sb.WriteString("-------------------------------------" + "\n")
	}
	fmt.Println(sb.String())
}

// LoadData retrieves specific credentials data based on the provided ID and displays it.
// It allows the user to print the information or save it to a file.
func (p *CredentialsProvider) LoadData(ctx context.Context) {
	red := color.New(color.FgRed).SprintFunc()

	if !p.state.IsAuthorized() {
		fmt.Println(red("You are not authorized, please use 'login' or 'register'"))
		return
	}

	if p.state.GetDirPath() == "" {
		fmt.Println(red("To proceed you must set working directory"))
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	cyanBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	req := model.BinaryDataLoadRequest{}
	fmt.Println(cyanBold("Input data ID to load credentials data:"))

	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("Input data ID as %s: ", yellow("'text'"))
	scanner.Scan()
	req.ID = scanner.Text()

	credentialsData, err := p.credentialsService.LoadCredentialsData(ctx, p.state.GetToken(), req.ID)
	if err != nil {
		lib.UnpackGRPCError(err)
		return
	}

	green := color.New(color.FgGreen).SprintFunc()
	fmt.Println(green("-------------------------------------"))

	var sb strings.Builder
	sb.WriteString("Credential login: " + credentialsData.Login + "\n")
	sb.WriteString("Credential password: " + credentialsData.Password + "\n")
	sb.WriteString("Credential metadata: " + credentialsData.MetaData + "\n")
	sb.WriteString("-------------------------------------" + "\n")
	fmt.Println(sb.String())
	fmt.Print(green("Write info to file or print (leave empty or write to file): "))
	scanner.Scan()
	path := scanner.Text()

	if len(path) == 0 {
		fmt.Print(sb.String())
		return
	}

	if p.state.GetDirPath() != "" {
		path = p.state.GetDirPath() + "/" + path
	}

	err = lib.SaveToFile(path, sb.String())
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Error writing to file with path %s, please try again\n", red(path))
		return
	}
}

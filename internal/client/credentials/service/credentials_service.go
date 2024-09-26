package service

import (
	"bufio"
	"context"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/lib"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/state"
	"github.com/fatih/color"
	"os"
	"strings"
)

type CredentialsService interface {
	SaveCredentials(ctx context.Context, token string, cred model.CredentialsPostRequest) (model.Credentials, error)
	LoadCredentials(ctx context.Context, token string, cred model.CredentialsLoadRequest) ([]model.Credentials, error)
}

type CredentialsProvider struct {
	credentialsService CredentialsService
	state              *state.ClientState
}

func NewCredentialsService(u CredentialsService, state *state.ClientState) *CredentialsProvider {
	return &CredentialsProvider{
		credentialsService: u,
		state:              state,
	}
}

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

func (p *CredentialsProvider) Load(ctx context.Context) {
	red := color.New(color.FgRed).SprintFunc()

	if !p.state.IsAuthorized() {
		fmt.Println(red("You are not authorized, please use 'login' or 'register'"))
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	cyanBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	req := model.CredentialsLoadRequest{}
	fmt.Println(cyanBold("Input filter data to load credentials 'login password metadata':"))

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

	creds, err := p.credentialsService.LoadCredentials(ctx, p.state.GetToken(), req)
	if err != nil {
		lib.UnpackGRPCError(err)
		return
	}

	fmt.Println("-------------------------------------")

	green := color.New(color.FgGreen).SprintFunc()

	var sb strings.Builder
	for _, cred := range creds {
		sb.WriteString("Credential login: " + cred.Login + "\n")
		sb.WriteString("Credential password: " + cred.Password + "\n")
		sb.WriteString("Credential metadata: " + cred.MetaData + "\n")
		sb.WriteString("-------------------------------------" + "\n")
	}

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

	fmt.Printf("Data successfully written to file %s\n", green(path))
}

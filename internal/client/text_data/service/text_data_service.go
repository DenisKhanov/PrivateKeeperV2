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

// TextDataService is an interface that defines methods for saving and loading text data.
// It allows interaction with the underlying text data service.
type TextDataService interface {
	SaveTextData(ctx context.Context, token string, text model.TextDataPostRequest) (model.TextData, error)
	LoadTextData(ctx context.Context, token string, dataID string) (model.TextData, error)
	LoadAllTextDataInfo(ctx context.Context, token string) ([]model.DataInfo, error)
}

// TextDataProvider implements the TextDataService interface and holds the state for user sessions.
type TextDataProvider struct {
	textDataService TextDataService    // Service to handle text data operations
	state           *state.ClientState // Holds the client state, including authorization and directory path
}

// NewTextDataService initializes a new TextDataProvider with the provided text data service and client state.
// It returns a pointer to the initialized TextDataProvider.
func NewTextDataService(u TextDataService, state *state.ClientState) *TextDataProvider {
	return &TextDataProvider{
		textDataService: u,
		state:           state,
	}
}

// Save prompts the user for text data and metadata, then saves it using the text data service.
// It checks for authorization and directory path before proceeding.
func (p *TextDataProvider) Save(ctx context.Context) {
	red := color.New(color.FgRed).SprintFunc()

	if !p.state.IsAuthorized() {
		fmt.Println(red("You are not authorized, please use 'login' or 'register'"))
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	cyanBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	req := model.TextDataPostRequest{}
	fmt.Println(cyanBold("Input text data 'text metadata':"))

	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("Input text data as %s: ", yellow("'text'"))
	scanner.Scan()
	data := scanner.Text()
	req.Text = data

	fmt.Printf("Input metadata as %s: ", yellow("'text'"))
	scanner.Scan()
	data = scanner.Text()
	req.MetaData = data

	_, err := p.textDataService.SaveTextData(ctx, p.state.GetToken(), req)
	if err != nil {
		lib.UnpackGRPCError(err)
		return
	}

	fmt.Println(color.New(color.FgGreen).SprintFunc()("Text data successfully saved"))
}

// LoadAllInfo retrieves and displays information about all text data stored by the user.
// It checks for authorization and the working directory before proceeding.
func (p *TextDataProvider) LoadAllInfo(ctx context.Context) {
	red := color.New(color.FgRed).SprintFunc()

	if !p.state.IsAuthorized() {
		fmt.Println(red("You are not authorized, please use 'login' or 'register'"))
		return
	}

	if p.state.GetDirPath() == "" {
		fmt.Println(red("To proceed you must set working directory"))
		return
	}

	textDataInfo, err := p.textDataService.LoadAllTextDataInfo(ctx, p.state.GetToken())
	if err != nil {
		logrus.WithError(err).Error("All user data info load failed")
		fmt.Println(red("All text data info load failed"), "please try again")
		lib.UnpackGRPCError(err)
		return
	}

	green := color.New(color.FgGreen).SprintFunc()
	fmt.Println(green("-------------------------------------"))

	var sb strings.Builder
	for _, dataInfo := range textDataInfo {
		sb.WriteString("Data ID: " + dataInfo.ID + "\n")
		sb.WriteString("Data type: " + dataInfo.DataType + "\n")
		sb.WriteString("Metadata : " + dataInfo.MetaData + "\n")
		sb.WriteString("Created at: : " + dataInfo.CreatedAt + "\n")
		sb.WriteString("-------------------------------------" + "\n")
	}
	fmt.Println(sb.String())
}

// LoadData retrieves and displays a specific text data entry based on the provided data ID.
// It prompts the user for the data ID and checks for authorization and directory path before proceeding.
func (p *TextDataProvider) LoadData(ctx context.Context) {
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
	fmt.Println(cyanBold("Input data ID to load text data:"))

	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("Input data ID as %s: ", yellow("'text'"))
	scanner.Scan()
	req.ID = scanner.Text()

	textData, err := p.textDataService.LoadTextData(ctx, p.state.GetToken(), req.ID)
	if err != nil {
		lib.UnpackGRPCError(err)
		return
	}

	green := color.New(color.FgGreen).SprintFunc()
	fmt.Println(green("-------------------------------------"))

	var sb strings.Builder
	sb.WriteString("Credential login: " + textData.Text + "\n")
	sb.WriteString("Credential metadata: " + textData.MetaData + "\n")
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

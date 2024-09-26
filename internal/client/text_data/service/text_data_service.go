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

type TextDataService interface {
	SaveTextData(ctx context.Context, token string, text model.TextDataPostRequest) (model.TextData, error)
	LoadTextData(ctx context.Context, token string, textData model.TextDataLoadRequest) ([]model.TextData, error)
}

type TextDataProvider struct {
	textDataService TextDataService
	state           *state.ClientState
}

func NewTextDataService(u TextDataService, state *state.ClientState) *TextDataProvider {
	return &TextDataProvider{
		textDataService: u,
		state:           state,
	}
}

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

func (p *TextDataProvider) Load(ctx context.Context) {
	red := color.New(color.FgRed).SprintFunc()

	if !p.state.IsAuthorized() {
		fmt.Println(red("You are not authorized, please use 'login' or 'register'"))
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	yellowBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	req := model.TextDataLoadRequest{}
	fmt.Println(yellowBold("Input filter data to load text data 'text metadata':"))

	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("Input text data as %s: ", yellow("'text'"))
	scanner.Scan()
	data := scanner.Text()
	req.Text = data

	fmt.Printf("Input metadata as %s: ", yellow("'text'"))
	scanner.Scan()
	data = scanner.Text()
	req.MetaData = data

	texts, err := p.textDataService.LoadTextData(ctx, p.state.GetToken(), req)
	if err != nil {
		lib.UnpackGRPCError(err)
		return
	}

	fmt.Println("-------------------------------------")

	green := color.New(color.FgGreen).SprintFunc()

	var sb strings.Builder
	for _, txt := range texts {
		sb.WriteString("Text data: " + txt.Text + "\n")
		sb.WriteString("Text metadata: " + txt.MetaData + "\n")
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

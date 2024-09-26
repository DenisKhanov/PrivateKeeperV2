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
	"path/filepath"
)

type BinaryDataService interface {
	SaveBinaryData(ctx context.Context, token string, bData model.BinaryDataPostRequest) (model.BinaryData, error)
	LoadBinaryData(ctx context.Context, token string, bData model.BinaryDataLoadRequest) ([]model.BinaryData, error)
}

type BinaryDataProvider struct {
	binaryDataService BinaryDataService
	state             *state.ClientState
}

func NewBinaryDataService(u BinaryDataService, state *state.ClientState) *BinaryDataProvider {
	return &BinaryDataProvider{
		binaryDataService: u,
		state:             state,
	}
}

func (p *BinaryDataProvider) Save(ctx context.Context) {
	red := color.New(color.FgRed).SprintFunc()

	if !p.state.IsAuthorized() {
		fmt.Println(red("You are not authorized, please use 'login' or 'register'"))
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	cyanBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	req := model.BinaryDataPostRequest{}
	fmt.Println(cyanBold("Input binary data to save 'path name extension metadata':"))

	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("Input path as %s: ", yellow("'text'"))
	scanner.Scan()
	path := scanner.Text()

	data, err := lib.LoadFromFile(path)
	if err != nil {
		fmt.Println("Error loading file please try again")
		return
	}

	req.Data = data

	fmt.Printf("Input name as %s: ", yellow("'text'"))
	scanner.Scan()
	req.Name = scanner.Text()

	fmt.Printf("Input extension as %s: ", yellow("'text'"))
	scanner.Scan()
	req.Extension = scanner.Text()

	fmt.Printf("Input metadata as %s: ", yellow("'text'"))
	scanner.Scan()
	req.MetaData = scanner.Text()

	_, err = p.binaryDataService.SaveBinaryData(ctx, p.state.GetToken(), req)
	if err != nil {
		lib.UnpackGRPCError(err)
		return
	}

	fmt.Println(color.New(color.FgGreen).SprintFunc()("Binary data successfully saved"))
}

func (p *BinaryDataProvider) Load(ctx context.Context) {
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
	fmt.Println(cyanBold("Input filter data to load binary data 'name metadata':"))

	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("Input name as %s: ", yellow("'text'"))
	req.Name = scanner.Text()

	fmt.Printf("Input metadata as %s: ", yellow("'text'"))
	scanner.Scan()
	req.MetaData = scanner.Text()

	bData, err := p.binaryDataService.LoadBinaryData(ctx, p.state.GetToken(), req)
	if err != nil {
		lib.UnpackGRPCError(err)
		return
	}

	for _, data := range bData {
		path := filepath.Join(p.state.GetDirPath(), "/", data.Name+"."+data.Extension)
		err = lib.SaveBinaryToFile(path, data.Data)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("Error writing to file with path %s, please try again\n", red(path))
			return
		}
	}

	fmt.Printf("Data successfully written to your working dir %s\n", p.state.GetDirPath())
}

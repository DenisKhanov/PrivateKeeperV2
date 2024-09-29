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
	"path/filepath"
	"strings"
)

// BinaryDataService defines methods for saving and loading binary data.
type BinaryDataService interface {
	SaveBinaryData(ctx context.Context, token string, bData model.BinaryDataPostRequest) (model.BinaryData, error)
	LoadBinaryData(ctx context.Context, token string, dataID string) (model.BinaryData, error)
	LoadAllBinaryDataInfo(ctx context.Context, token string) ([]model.DataInfo, error)
}

// BinaryDataProvider implements the BinaryDataService interface and manages client-side operations.
type BinaryDataProvider struct {
	binaryDataService BinaryDataService
	state             *state.ClientState
}

// NewBinaryDataService returns a new instance of BinaryDataProvider.
func NewBinaryDataService(u BinaryDataService, state *state.ClientState) *BinaryDataProvider {
	return &BinaryDataProvider{
		binaryDataService: u,
		state:             state,
	}
}

// Save prompts the user to input binary data details, including path, name, extension, and metadata,
// and saves the binary data using the BinaryDataService. It ensures the user is authorized.
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

// LoadAllInfo retrieves and displays metadata for all saved binary data. The user must be authorized,
// and a working directory must be set.
func (p *BinaryDataProvider) LoadAllInfo(ctx context.Context) {
	red := color.New(color.FgRed).SprintFunc()

	if !p.state.IsAuthorized() {
		fmt.Println(red("You are not authorized, please use 'login' or 'register'"))
		return
	}

	if p.state.GetDirPath() == "" {
		fmt.Println(red("To proceed you must set working directory"))
		return
	}

	binariesDataInfo, err := p.binaryDataService.LoadAllBinaryDataInfo(ctx, p.state.GetToken())
	if err != nil {
		logrus.WithError(err).Error("All user data info load failed")
		fmt.Println(red("All binaries data info load failed"), "please try again")
		lib.UnpackGRPCError(err)
		return
	}

	green := color.New(color.FgGreen).SprintFunc()
	fmt.Println(green("-------------------------------------"))

	var sb strings.Builder
	for _, dataInfo := range binariesDataInfo {
		sb.WriteString("Data ID: " + dataInfo.ID + "\n")
		sb.WriteString("Data type: " + dataInfo.DataType + "\n")
		sb.WriteString("Metadata : " + dataInfo.MetaData + "\n")
		sb.WriteString("Created at: : " + dataInfo.CreatedAt + "\n")
		sb.WriteString("-------------------------------------" + "\n")
	}
	fmt.Println(sb.String())
}

// LoadData loads specific binary data by its ID, saves it to the client's working directory,
// and ensures the user is authorized and has set a working directory.
func (p *BinaryDataProvider) LoadData(ctx context.Context) {
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
	fmt.Println(cyanBold("Input data ID to load binary data:"))

	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("Input data ID as %s: ", yellow("'text'"))
	scanner.Scan()
	req.ID = scanner.Text()

	bData, err := p.binaryDataService.LoadBinaryData(ctx, p.state.GetToken(), req.ID)
	if err != nil {
		lib.UnpackGRPCError(err)
		return
	}

	path := filepath.Join(p.state.GetDirPath(), "/", bData.Name+"."+bData.Extension)
	err = lib.SaveBinaryToFile(path, bData.Data)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Error writing to file with path %s, please try again\n", red(path))
		return
	}

	fmt.Printf("Data successfully written to your working dir %s\n", p.state.GetDirPath())
}

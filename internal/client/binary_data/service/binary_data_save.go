package service

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/lib"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
)

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

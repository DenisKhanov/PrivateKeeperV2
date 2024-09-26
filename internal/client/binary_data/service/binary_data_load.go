package service

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/lib"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
)

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

package service

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/lib"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
)

func (p *TextDataProvider) Save() {
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

	_, err := p.textDataService.SaveTextData(p.state.GetToken(), req)
	if err != nil {
		lib.UnpackGRPCError(err)
		return
	}

	fmt.Println(color.New(color.FgGreen).SprintFunc()("Text data successfully saved"))
}

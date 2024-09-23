package service

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/lib"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
)

func (p *CreditCardProvider) Save() {
	red := color.New(color.FgRed).SprintFunc()

	if !p.state.IsAuthorized() {
		fmt.Println(red("You are not authorized, please use 'login' or 'register'"))
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	cyanBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	req := model.CreditCardPostRequest{}
	fmt.Println(cyanBold("Input credit card data 'number owner expires cvv pin metadata':"))

	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("Input number in format %s: ", yellow("'dddd dddd dddd dddd'"))
	scanner.Scan()
	data := scanner.Text()
	req.Number = data

	fmt.Printf("Input owner in format %s: ", yellow("'name surname'"))
	scanner.Scan()
	data = scanner.Text()
	req.OwnerName = data

	fmt.Printf("Input expiry date in format %s: ", yellow("'dd-mm-yyyy'"))
	scanner.Scan()
	data = scanner.Text()
	req.ExpiresAt = data

	fmt.Printf("Input pin in format %s: ", yellow("'dddd'"))
	scanner.Scan()
	data = scanner.Text()
	req.PinCode = data

	fmt.Printf("Input cvv in format %s: ", yellow("'ddd'"))
	scanner.Scan()
	data = scanner.Text()
	req.CVV = data

	fmt.Printf("Input metadata as %s: ", yellow("'text'"))
	scanner.Scan()
	data = scanner.Text()
	req.MetaData = data

	_, err := p.creditCardService.SaveCreditCard(p.state.GetToken(), req)
	if err != nil {
		lib.UnpackGRPCError(err)
		return
	}

	fmt.Println(color.New(color.FgGreen).SprintFunc()("Card successfully saved"))
}

package service

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/lib"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
)

func (p *CreditCardProvider) Load(ctx context.Context) {
	red := color.New(color.FgRed).SprintFunc()

	if !p.state.IsAuthorized() {
		fmt.Println(red("You are not authorized, please use 'login' or 'register'"))
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	cyanBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	req := model.CreditCardLoadRequest{}
	fmt.Println(cyanBold("Input filter data to load credit cards 'number owner cvv pin metadata expires_after expires_before':"))

	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("Input number in format %s: ", yellow("'dddd dddd dddd dddd'"))
	scanner.Scan()
	data := scanner.Text()
	req.Number = data

	fmt.Printf("Input owner in format %s: ", yellow("'name surname'"))
	scanner.Scan()
	data = scanner.Text()
	req.Owner = data

	fmt.Printf("Input cvv in format %s: ", yellow("'ddd'"))
	scanner.Scan()
	data = scanner.Text()
	req.CvvCode = data

	fmt.Printf("Input pin in format %s: ", yellow("'dddd'"))
	scanner.Scan()
	data = scanner.Text()
	req.PinCode = data

	fmt.Printf("Input metadata as %s: ", yellow("'text'"))
	scanner.Scan()
	data = scanner.Text()
	req.Metadata = data

	fmt.Printf("Input expires_after in format %s: ", yellow("'dd-mm-yyyy'"))
	scanner.Scan()
	data = scanner.Text()
	req.ExpiresAfter = data

	fmt.Printf("Input expires_before in format %s: ", yellow("'dd-mm-yyyy'"))
	scanner.Scan()
	data = scanner.Text()
	req.ExpiresBefore = data

	cards, err := p.creditCardService.LoadCreditCard(p.state.GetToken(), req)
	if err != nil {
		lib.UnpackGRPCError(err)
		return
	}

	fmt.Println("-------------------------------------")

	green := color.New(color.FgGreen).SprintFunc()

	var sb strings.Builder
	for _, card := range cards {
		sb.WriteString("Card number: " + card.Number + "\n")
		sb.WriteString("Card owner: " + card.OwnerName + "\n")
		sb.WriteString("Card expires at: " + card.ExpiresAt + "\n")
		sb.WriteString("Card cvv: " + card.CVV + "\n")
		sb.WriteString("Card in code: " + card.PinCode + "\n")
		sb.WriteString("Card metadata: " + card.MetaData + "\n")
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

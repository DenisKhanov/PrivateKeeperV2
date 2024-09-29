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

// CreditCardService defines the interface for credit card operations.
// It includes methods to save, load, and retrieve information about credit cards.
type CreditCardService interface {
	SaveCreditCard(ctx context.Context, token string, card model.CreditCardPostRequest) (model.CreditCard, error)
	LoadCreditCardData(ctx context.Context, token string, dataID string) (model.CreditCard, error)
	LoadAllCreditCardDataInfo(ctx context.Context, token string) ([]model.DataInfo, error)
}

// CreditCardProvider provides methods for interacting with credit card services.
// It holds a reference to a CreditCardService and client state.
type CreditCardProvider struct {
	creditCardService CreditCardService  // The service responsible for credit card operations
	state             *state.ClientState // Client state to manage user session
}

// NewUserService creates a new CreditCardProvider instance.
// It takes a CreditCardService and ClientState as parameters and returns the provider instance.
func NewUserService(u CreditCardService, state *state.ClientState) *CreditCardProvider {
	return &CreditCardProvider{
		creditCardService: u,
		state:             state,
	}
}

// Save prompts the user for credit card details and saves them using the credit card service.
// It checks for user authorization and validates input before saving the card information.
func (p *CreditCardProvider) Save(ctx context.Context) {
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

	fmt.Printf("Input data description as %s: ", yellow("'text'"))
	scanner.Scan()
	data = scanner.Text()
	req.MetaData = data

	_, err := p.creditCardService.SaveCreditCard(ctx, p.state.GetToken(), req)
	if err != nil {
		lib.UnpackGRPCError(err)
		return
	}

	fmt.Println(color.New(color.FgGreen).SprintFunc()("Card successfully saved"))
}

// LoadAllInfo retrieves and displays information about all saved credit cards.
// It checks for user authorization and working directory setup before loading the data.
func (p *CreditCardProvider) LoadAllInfo(ctx context.Context) {
	red := color.New(color.FgRed).SprintFunc()

	if !p.state.IsAuthorized() {
		fmt.Println(red("You are not authorized, please use 'login' or 'register'"))
		return
	}

	if p.state.GetDirPath() == "" {
		fmt.Println(red("To proceed you must set working directory"))
		return
	}

	cardDataInfo, err := p.creditCardService.LoadAllCreditCardDataInfo(ctx, p.state.GetToken())
	if err != nil {
		logrus.WithError(err).Error("All user data info load failed")
		fmt.Println(red("All credit card data info load failed"), "please try again")
		lib.UnpackGRPCError(err)
		return
	}

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Println(green("-------------------------------------"))

	var sb strings.Builder
	for _, dataInfo := range cardDataInfo {
		sb.WriteString("Data ID: " + dataInfo.ID + "\n")
		sb.WriteString("Data type: " + dataInfo.DataType + "\n")
		sb.WriteString("Metadata : " + dataInfo.MetaData + "\n")
		sb.WriteString("Created at: : " + dataInfo.CreatedAt + "\n")
		sb.WriteString(green("-------------------------------------") + "\n")
	}
	if len(cardDataInfo) > 0 {
		fmt.Println(sb.String())
	} else {
		fmt.Println(yellow("Your haven't saved any data or data load filed, please try again"))
	}
}

// LoadData retrieves and displays the details of a specific credit card based on its ID.
// It checks for user authorization and working directory setup before loading the card data.
func (p *CreditCardProvider) LoadData(ctx context.Context) {
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
	fmt.Printf("Input data ID as %s: ", yellow("'example (b7fa5761-7e83-11ef-a610-0242ac140004)'"))
	scanner.Scan()
	req.ID = scanner.Text()

	cardData, err := p.creditCardService.LoadCreditCardData(ctx, p.state.GetToken(), req.ID)
	if err != nil {
		lib.UnpackGRPCError(err)
		return
	}

	green := color.New(color.FgGreen).SprintFunc()

	var sb strings.Builder
	sb.WriteString(red("-------------------------------------") + "\n")
	sb.WriteString("Card number: " + cardData.Number + "\n")
	sb.WriteString("Card owner: " + cardData.OwnerName + "\n")
	sb.WriteString("Card expires at: " + cardData.ExpiresAt + "\n")
	sb.WriteString("Card cvv: " + cardData.CVV + "\n")
	sb.WriteString("Card in code: " + cardData.PinCode + "\n")
	sb.WriteString("Card metadata: " + cardData.MetaData + "\n")
	sb.WriteString(red("-------------------------------------") + "\n")
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

	fmt.Printf("Data successfully written to file %s\n", green(path))
}

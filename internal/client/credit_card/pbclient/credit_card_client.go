package pbclient

import (
	"context"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credit_card"
	"google.golang.org/grpc/metadata"
)

// CreditCardPBClient is a client for interacting with the credit card service over gRPC.
// It provides methods to save and load credit card data.
type CreditCardPBClient struct {
	creditCardService pb.CreditCardServiceClient
}

// NewCreditCardPBClient creates a new instance of CreditCardPBClient.
// It takes a gRPC client for the credit card service as an argument and returns a pointer to the new client.
func NewCreditCardPBClient(u pb.CreditCardServiceClient) *CreditCardPBClient {
	return &CreditCardPBClient{
		creditCardService: u,
	}
}

// SaveCreditCard saves a new credit card using the credit card service.
// It takes a context, an authentication token, and a CreditCardPostRequest containing the card details.
// It returns the saved CreditCard and any error encountered during the process.
func (u *CreditCardPBClient) SaveCreditCard(ctx context.Context, token string, card model.CreditCardPostRequest) (model.CreditCard, error) {
	req := &pb.PostCreditCardRequest{
		Number:    card.Number,
		OwnerName: card.OwnerName,
		ExpiresAt: card.ExpiresAt,
		CvvCode:   card.CVV,
		PinCode:   card.PinCode,
		Metadata:  card.MetaData,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := u.creditCardService.PostSaveCreditCard(ctx, req)
	if err != nil {
		return model.CreditCard{}, err
	}

	creditCard := model.CreditCard{
		Number:    resp.Number,
		OwnerName: resp.OwnerName,
		ExpiresAt: resp.ExpiresAt,
		CVV:       resp.CvvCode,
		PinCode:   resp.PinCode,
		MetaData:  resp.Metadata,
	}

	return creditCard, nil
}

// LoadAllCreditCardDataInfo retrieves information about all saved credit cards.
// It takes a context and an authentication token as arguments.
// It returns a slice of DataInfo containing the details of each credit card and any error encountered.
func (u *CreditCardPBClient) LoadAllCreditCardDataInfo(ctx context.Context, token string) ([]model.DataInfo, error) {
	req := &pb.GetAllCreditCardInfoRequest{}

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := u.creditCardService.GetLoadAllCreditCardDataInfo(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("load credit card data: %w", err)
	}

	cards := make([]model.DataInfo, 0, len(resp.Cards))
	for _, data := range resp.Cards {
		cards = append(cards, model.DataInfo{
			ID:        data.Id,
			DataType:  data.DataType,
			MetaData:  data.Metadata,
			CreatedAt: data.CreatedAt,
		})
	}

	return cards, nil
}

// LoadCreditCardData retrieves the details of a specific credit card using its ID.
// It takes a context, an authentication token, and the credit card ID as arguments.
// It returns the CreditCard associated with the given ID and any error encountered during the process.
func (u *CreditCardPBClient) LoadCreditCardData(ctx context.Context, token string, dataID string) (model.CreditCard, error) {
	req := &pb.GetCreditCardRequest{
		Id: dataID,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := u.creditCardService.GetLoadCreditCard(ctx, req)
	if err != nil {
		return model.CreditCard{}, fmt.Errorf("load credit card data: %w", err)
	}
	data := resp.CardData
	creditCard := model.CreditCard{
		Number:    data.Number,
		OwnerName: data.OwnerName,
		ExpiresAt: data.ExpiresAt,
		CVV:       data.CvvCode,
		PinCode:   data.PinCode,
		MetaData:  data.Metadata,
	}

	return creditCard, nil
}

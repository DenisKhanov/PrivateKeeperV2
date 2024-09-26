package pbclient

import (
	"context"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credit_card"
	"google.golang.org/grpc/metadata"
)

type CreditCardPBClient struct {
	creditCardService pb.CreditCardServiceClient
}

func NewCreditCardPBClient(u pb.CreditCardServiceClient) *CreditCardPBClient {
	return &CreditCardPBClient{
		creditCardService: u,
	}
}

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

func (u *CreditCardPBClient) LoadCreditCard(ctx context.Context, token string, card model.CreditCardLoadRequest) ([]model.CreditCard, error) {
	req := &pb.GetCreditCardRequest{
		Number:        card.Number,
		Owner:         card.Owner,
		CvvCode:       card.CvvCode,
		PinCode:       card.PinCode,
		Metadata:      card.Metadata,
		ExpiresAfter:  card.ExpiresAfter,
		ExpiresBefore: card.ExpiresBefore,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := u.creditCardService.GetLoadCreditCard(ctx, req)
	if err != nil {
		return nil, err
	}

	cards := make([]model.CreditCard, 0, len(resp.Cards))
	for _, card := range resp.Cards {
		cards = append(cards, model.CreditCard{
			Number:    card.Number,
			OwnerName: card.OwnerName,
			ExpiresAt: card.ExpiresAt,
			CVV:       card.CvvCode,
			PinCode:   card.PinCode,
			MetaData:  card.Metadata,
		})
	}

	return cards, nil
}

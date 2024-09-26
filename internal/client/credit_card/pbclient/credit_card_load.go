package pbclient

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credit_card"
)

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

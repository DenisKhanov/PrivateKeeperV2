package pbclient

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credit_card"
)

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

package grpchandlers

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credit_card"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/credit_card/specification"
)

func (h *CreditCardHandler) GetLoadCreditCard(ctx context.Context, in *pb.GetCreditCardRequest) (*pb.GetCreditCardResponse, error) {
	spec, err := specification.NewCreditCardSpecification(in)
	if err != nil {
		logrus.WithError(err).Error("Error while creating credit card specification: ")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	cards, err := h.creditCardService.LoadAllCreditCard(ctx, spec)
	if err != nil {
		logrus.WithError(err).Error("Error while loading credit cards: ")
		return nil, status.Error(codes.Internal, "internal error")
	}

	creditCards := make([]*pb.CreditCard, 0, len(cards))
	for _, v := range cards {
		creditCards = append(creditCards, &pb.CreditCard{
			Id:        v.ID,
			OwnerId:   v.OwnerID,
			Number:    v.Number,
			OwnerName: v.OwnerName,
			ExpiresAt: v.ExpiresAt,
			CvvCode:   v.CVV,
			PinCode:   v.PinCode,
			Metadata:  v.MetaData,
			CreatedAt: v.CreatedAt.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.GetCreditCardResponse{Cards: creditCards}, nil
}

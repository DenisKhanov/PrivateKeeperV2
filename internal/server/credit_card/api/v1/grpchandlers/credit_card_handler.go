package grpchandlers

import (
	"context"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/lib"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credit_card"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/credit_card/specification"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

type CreditCardService interface {
	SaveCreditCard(ctx context.Context, req model.CreditCardPostRequest) (model.CreditCard, error)
	LoadAllCreditCard(ctx context.Context, spec specification.CreditCardSpecification) ([]model.CreditCard, error)
}

type Validator interface {
	ValidatePostRequest(req *model.CreditCardPostRequest) (map[string]string, bool)
}

type CreditCardHandler struct {
	creditCardService CreditCardService
	pb.UnimplementedCreditCardServiceServer
	validator Validator
}

func New(creditCardService CreditCardService, validator Validator) *CreditCardHandler {
	return &CreditCardHandler{
		creditCardService: creditCardService,
		validator:         validator,
	}
}

func (h *CreditCardHandler) PostSaveCreditCard(ctx context.Context, in *pb.PostCreditCardRequest) (*pb.PostCreditCardResponse, error) {
	req := model.CreditCardPostRequest{
		Number:    in.Number,
		OwnerName: in.OwnerName,
		ExpiresAt: in.ExpiresAt,
		CVV:       in.CvvCode,
		PinCode:   in.PinCode,
		MetaData:  in.Metadata,
	}

	report, ok := h.validator.ValidatePostRequest(&req)
	if !ok {
		logrus.Info("Unable to register user: invalid user request")
		logrus.Infof("violated_fields %v", report)
		return nil, lib.ProcessValidationError("invalid credit_card post request", report)
	}

	creditCard, err := h.creditCardService.SaveCreditCard(ctx, req)
	if err != nil {
		logrus.WithError(err).Error("Unable to save credit_card")
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.PostCreditCardResponse{
		Id:        creditCard.ID,
		OwnerId:   creditCard.OwnerID,
		Number:    creditCard.Number,
		OwnerName: creditCard.OwnerName,
		ExpiresAt: creditCard.ExpiresAt,
		CvvCode:   creditCard.CVV,
		PinCode:   creditCard.PinCode,
		Metadata:  creditCard.MetaData,
		CreatedAt: creditCard.CreatedAt.Format(time.RFC3339),
	}, nil
}

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
		})
	}

	return &pb.GetCreditCardResponse{Cards: creditCards}, nil
}

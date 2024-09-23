package grpchandlers

import (
	"context"

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

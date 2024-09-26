package service

import (
	"context"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/state"
)

type CreditCardService interface {
	SaveCreditCard(ctx context.Context, token string, card model.CreditCardPostRequest) (model.CreditCard, error)
	LoadCreditCard(ctx context.Context, token string, card model.CreditCardLoadRequest) ([]model.CreditCard, error)
}

type CreditCardProvider struct {
	creditCardService CreditCardService
	state             *state.ClientState
}

func NewUserService(u CreditCardService, state *state.ClientState) *CreditCardProvider {
	return &CreditCardProvider{
		creditCardService: u,
		state:             state,
	}
}

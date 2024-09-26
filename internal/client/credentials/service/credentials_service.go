package service

import (
	"context"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/state"
)

type CredentialsService interface {
	SaveCredentials(ctx context.Context, token string, cred model.CredentialsPostRequest) (model.Credentials, error)
	LoadCredentials(ctx context.Context, token string, cred model.CredentialsLoadRequest) ([]model.Credentials, error)
}

type CredentialsProvider struct {
	credentialsService CredentialsService
	state              *state.ClientState
}

func NewCredentialsService(u CredentialsService, state *state.ClientState) *CredentialsProvider {
	return &CredentialsProvider{
		credentialsService: u,
		state:              state,
	}
}

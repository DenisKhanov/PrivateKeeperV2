package service

import (
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/state"
)

type TextDataService interface {
	SaveTextData(token string, text model.TextDataPostRequest) (model.TextData, error)
	LoadTextData(token string, textData model.TextDataLoadRequest) ([]model.TextData, error)
}

type TextDataProvider struct {
	textDataService TextDataService
	state           *state.ClientState
}

func NewTextDataService(u TextDataService, state *state.ClientState) *TextDataProvider {
	return &TextDataProvider{
		textDataService: u,
		state:           state,
	}
}

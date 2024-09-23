package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/auth"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/keyextraction"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

func (s *CredentialsService) SaveCredentials(ctx context.Context, req model.CredentialsPostRequest) (model.Credentials, error) {
	userID, ok := ctx.Value(auth.UserIDContextKey("userID")).(string)
	if !ok {
		return model.Credentials{}, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(keyextraction.UserKeyContextKey("userKey")).([]byte)
	if !ok {
		return model.Credentials{}, fmt.Errorf("failed to get userKey from context")
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return model.Credentials{}, fmt.Errorf("new uuid: %w", err)
	}

	card := model.CredentialsCryptData{
		Login:    req.Login,
		Password: req.Password,
	}

	data, err := json.Marshal(card)
	if err != nil {
		return model.Credentials{}, fmt.Errorf("marshal: %w", err)
	}

	cryptData, err := s.crypt.Encrypt(userKey, data)
	if err != nil {
		return model.Credentials{}, fmt.Errorf("encrypt data: %w", err)
	}

	dataToSave := model.Data{
		ID:       id.String(),
		OwnerID:  userID,
		Type:     s.dataType,
		Data:     cryptData,
		MetaData: req.MetaData,
	}

	savedCredentials, err := s.repository.Insert(ctx, dataToSave)
	if err != nil {
		return model.Credentials{}, fmt.Errorf("insert credentials: %w", err)
	}

	return model.Credentials{
		ID:        savedCredentials.ID,
		OwnerID:   savedCredentials.OwnerID,
		Login:     req.Login,
		Password:  req.Password,
		MetaData:  savedCredentials.MetaData,
		CreatedAt: savedCredentials.CreatedAt,
		UpdatedAt: savedCredentials.UpdatedAt,
	}, nil
}

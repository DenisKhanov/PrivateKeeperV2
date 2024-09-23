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

func (s *TextDataService) SaveTextData(ctx context.Context, req model.TextDataPostRequest) (model.TextData, error) {
	userID, ok := ctx.Value(auth.UserIDContextKey("userID")).(string)
	if !ok {
		return model.TextData{}, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(keyextraction.UserKeyContextKey("userKey")).([]byte)
	if !ok {
		return model.TextData{}, fmt.Errorf("failed to get userKey from context")
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return model.TextData{}, fmt.Errorf("new uuid: %w", err)
	}

	text := model.TextCryptData{
		Text: req.Text,
	}

	data, err := json.Marshal(text)
	if err != nil {
		return model.TextData{}, fmt.Errorf("marshal: %w", err)
	}

	cryptData, err := s.crypt.Encrypt(userKey, data)
	if err != nil {
		return model.TextData{}, fmt.Errorf("encrypt data: %w", err)
	}

	dataToSave := model.Data{
		ID:       id.String(),
		OwnerID:  userID,
		Type:     s.dataType,
		Data:     cryptData,
		MetaData: req.MetaData,
	}

	savedTextData, err := s.repository.Insert(ctx, dataToSave)
	if err != nil {
		return model.TextData{}, fmt.Errorf("insert credit card: %w", err)
	}

	return model.TextData{
		ID:        savedTextData.ID,
		OwnerID:   savedTextData.OwnerID,
		Text:      req.Text,
		MetaData:  savedTextData.MetaData,
		CreatedAt: savedTextData.CreatedAt,
		UpdatedAt: savedTextData.UpdatedAt,
	}, nil
}

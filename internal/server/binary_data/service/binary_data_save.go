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

func (s *BinaryDataService) SaveBinaryData(ctx context.Context, req model.BinaryDataPostRequest) (model.BinaryData, error) {
	userID, ok := ctx.Value(auth.UserIDContextKey("userID")).(string)
	if !ok {
		return model.BinaryData{}, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(keyextraction.UserKeyContextKey("userKey")).([]byte)
	if !ok {
		return model.BinaryData{}, fmt.Errorf("failed to get userKey from context")
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return model.BinaryData{}, fmt.Errorf("new uuid: %w", err)
	}

	binary := model.BinaryCryptData{
		Name:      req.Name,
		Extension: req.Extension,
		Data:      req.Data,
	}

	data, err := json.Marshal(binary)
	if err != nil {
		return model.BinaryData{}, fmt.Errorf("marshal: %w", err)
	}

	cryptData, err := s.crypt.Encrypt(userKey, data)
	if err != nil {
		return model.BinaryData{}, fmt.Errorf("encrypt data: %w", err)
	}

	dataToSave := model.Data{
		ID:       id.String(),
		OwnerID:  userID,
		Type:     s.dataType,
		Data:     cryptData,
		MetaData: req.MetaData,
	}

	savedBinaryData, err := s.repository.Insert(ctx, dataToSave)
	if err != nil {
		return model.BinaryData{}, fmt.Errorf("insert credit card: %w", err)
	}

	return model.BinaryData{
		ID:        savedBinaryData.ID,
		OwnerID:   savedBinaryData.OwnerID,
		Name:      req.Name,
		Extension: req.Extension,
		Data:      req.Data,
		MetaData:  savedBinaryData.MetaData,
		CreatedAt: savedBinaryData.CreatedAt,
		UpdatedAt: savedBinaryData.UpdatedAt,
	}, nil
}

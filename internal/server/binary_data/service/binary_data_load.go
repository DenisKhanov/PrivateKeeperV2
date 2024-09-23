package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/binary_data/specification"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/auth"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/keyextraction"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

func (s *BinaryDataService) LoadAllBinaryData(ctx context.Context, spec specification.BinaryDataSpecification) ([]model.BinaryData, error) {
	userID, ok := ctx.Value(auth.UserIDContextKey("userID")).(string)
	if !ok {
		return nil, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(keyextraction.UserKeyContextKey("userKey")).([]byte)
	if !ok {
		return nil, fmt.Errorf("failed to get userKey from context")
	}

	encryptedBinaryData, err := s.repository.SelectAll(ctx, userID, s.dataType)
	if err != nil {
		return nil, fmt.Errorf("select all binary_data: %w", err)
	}

	texts := make([]model.BinaryData, 0, len(encryptedBinaryData))
	for _, encryptedBinary := range encryptedBinaryData {
		decryptedData, err := s.crypt.Decrypt(userKey, encryptedBinary.Data)
		if err != nil {
			return nil, fmt.Errorf("decrypt data: %w", err)
		}

		var data model.BinaryCryptData
		err = json.Unmarshal(decryptedData, &data)
		if err != nil {
			return nil, fmt.Errorf("unmarshal data: %w", err)
		}

		texts = append(texts, model.BinaryData{
			ID:        encryptedBinary.ID,
			OwnerID:   encryptedBinary.OwnerID,
			Name:      data.Name,
			Extension: data.Extension,
			Data:      data.Data,
			MetaData:  encryptedBinary.MetaData,
			CreatedAt: encryptedBinary.CreatedAt,
			UpdatedAt: encryptedBinary.UpdatedAt,
		})
	}

	predicates := spec.MakeFilterPredicates()
	var filteredBinaryData []model.BinaryData
	for _, binary := range texts {
		take := true
		for _, filteredBinaryDataWithSpec := range predicates {
			if !filteredBinaryDataWithSpec(spec, binary) {
				take = false
				break
			}
		}
		if take {
			filteredBinaryData = append(filteredBinaryData, binary)
		}
	}

	return filteredBinaryData, nil
}

package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/auth"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/keyextraction"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/text_data/specification"
)

func (s *TextDataService) LoadAllTextData(ctx context.Context, spec specification.TextDataSpecification) ([]model.TextData, error) {
	userID, ok := ctx.Value(auth.UserIDContextKey("userID")).(string)
	if !ok {
		return nil, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(keyextraction.UserKeyContextKey("userKey")).([]byte)
	if !ok {
		return nil, fmt.Errorf("failed to get userKey from context")
	}

	encryptedTextData, err := s.repository.SelectAll(ctx, userID, s.dataType)
	if err != nil {
		return nil, fmt.Errorf("select all text_data: %w", err)
	}

	texts := make([]model.TextData, 0, len(encryptedTextData))
	for _, encryptedText := range encryptedTextData {
		decryptedData, err := s.crypt.Decrypt(userKey, encryptedText.Data)
		if err != nil {
			return nil, fmt.Errorf("decrypt data: %w", err)
		}

		var data model.TextCryptData
		err = json.Unmarshal(decryptedData, &data)
		if err != nil {
			return nil, fmt.Errorf("unmarshal data: %w", err)
		}

		texts = append(texts, model.TextData{
			ID:        encryptedText.ID,
			OwnerID:   encryptedText.OwnerID,
			Text:      data.Text,
			MetaData:  encryptedText.MetaData,
			CreatedAt: encryptedText.CreatedAt,
			UpdatedAt: encryptedText.UpdatedAt,
		})
	}

	predicates := spec.MakeFilterPredicates()
	var filteredTextData []model.TextData
	for _, text := range texts {
		take := true
		for _, filteredTextDataWithSpec := range predicates {
			if !filteredTextDataWithSpec(spec, text) {
				take = false
				break
			}
		}
		if take {
			filteredTextData = append(filteredTextData, text)
		}
	}

	return filteredTextData, nil
}

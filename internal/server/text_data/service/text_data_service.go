package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/auth"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/keyextraction"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/text_data/specification"
	"github.com/google/uuid"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/DenisKhanov/PrivateKeeperV2/pkg/jwtmanager"
)

const textData = "text_data"

type DataRepository interface {
	Insert(ctx context.Context, data model.Data) (model.Data, error)
	SelectAll(ctx context.Context, userID, dataType string) ([]model.Data, error)
}

type CryptService interface {
	Encrypt(key, data []byte) ([]byte, error)
	Decrypt(key, data []byte) ([]byte, error)
	GenerateKey() ([]byte, error)
}

type TextDataService struct {
	repository DataRepository
	crypt      CryptService
	jwtManager *jwtmanager.JWTManager
	dataType   string
}

func New(repository DataRepository, crypt CryptService, jwtManager *jwtmanager.JWTManager) *TextDataService {
	return &TextDataService{
		repository: repository,
		crypt:      crypt,
		jwtManager: jwtManager,
		dataType:   textData,
	}
}

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

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/DenisKhanov/PrivateKeeperV2/pkg/jwtmanager"
)

const textData = "text_data" // Define a constant for the data type

// DataRepository interface defines methods for data persistence
type DataRepository interface {
	Insert(ctx context.Context, data model.Data) (model.Data, error)
	SelectAll(ctx context.Context, userID, dataType string) ([]model.Data, error)
	SelectByID(ctx context.Context, userID, dataType, dataID string) (model.Data, error)
}

// CryptService interface defines methods for encryption and decryption
type CryptService interface {
	Encrypt(key, data []byte) ([]byte, error)
	Decrypt(key, data []byte) ([]byte, error)
	GenerateKey() ([]byte, error)
}

// TextDataService provides methods to handle text data operations
type TextDataService struct {
	repository DataRepository         // Repository for data operations
	crypt      CryptService           // Service for encryption and decryption
	jwtManager *jwtmanager.JWTManager // JWT management
	dataType   string                 // Type of data this service handles
}

// New initializes a new TextDataService instance
func New(repository DataRepository, crypt CryptService, jwtManager *jwtmanager.JWTManager) *TextDataService {
	return &TextDataService{
		repository: repository,
		crypt:      crypt,
		jwtManager: jwtManager,
		dataType:   textData,
	}
}

// SaveTextData saves the provided text data to the repository
func (s *TextDataService) SaveTextData(ctx context.Context, req model.TextDataPostRequest) (model.TextData, error) {
	userID, ok := ctx.Value(model.UserIDKey).(string)
	if !ok {
		return model.TextData{}, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(model.UserKey).([]byte)
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
	}, nil
}

// LoadAllTextInfo retrieves all text data information for the user
func (s *TextDataService) LoadAllTextInfo(ctx context.Context) ([]model.DataInfo, error) {
	userID, ok := ctx.Value(model.UserIDKey).(string)
	if !ok {
		return nil, fmt.Errorf("failed to get userID from context")
	}

	encryptedTextData, err := s.repository.SelectAll(ctx, userID, s.dataType)
	if err != nil {
		return nil, fmt.Errorf("select all text_data: %w", err)
	}

	textDataInfo := make([]model.DataInfo, 0, len(encryptedTextData))
	for _, encryptedText := range encryptedTextData {
		textDataInfo = append(textDataInfo, model.DataInfo{
			ID:        encryptedText.ID,
			DataType:  s.dataType,
			MetaData:  encryptedText.MetaData,
			CreatedAt: encryptedText.CreatedAt,
		})
	}
	return textDataInfo, nil
}

// LoadTextData retrieves and decrypts text data by its ID
func (s *TextDataService) LoadTextData(ctx context.Context, dataID string) (model.TextData, error) {
	userID, ok := ctx.Value(model.UserIDKey).(string)
	if !ok {
		return model.TextData{}, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(model.UserKey).([]byte)
	if !ok {
		return model.TextData{}, fmt.Errorf("failed to get userKey from context")
	}

	encryptedTextData, err := s.repository.SelectByID(ctx, userID, s.dataType, dataID)
	if err != nil {
		return model.TextData{}, fmt.Errorf("select all text_data: %w", err)
	}
	decryptedData, err := s.crypt.Decrypt(userKey, encryptedTextData.Data)
	if err != nil {
		return model.TextData{}, fmt.Errorf("decrypt text: %w", err)
	}

	var decryptedTextData model.TextCryptData
	err = json.Unmarshal(decryptedData, &decryptedTextData)
	if err != nil {
		return model.TextData{}, fmt.Errorf("unmarshal text: %w", err)
	}

	text := model.TextData{
		ID:        encryptedTextData.ID,
		OwnerID:   encryptedTextData.OwnerID,
		Text:      decryptedTextData.Text,
		MetaData:  encryptedTextData.MetaData,
		CreatedAt: encryptedTextData.CreatedAt,
	}

	return text, nil
}

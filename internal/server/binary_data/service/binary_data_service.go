package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"

	"github.com/DenisKhanov/PrivateKeeperV2/pkg/jwtmanager"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

const binaryData = "binary_data" // Define a constant for binary data type

// DataRepository defines methods for interacting with the data storage layer.
type DataRepository interface {
	Insert(ctx context.Context, data model.Data) (model.Data, error)
	SelectAll(ctx context.Context, userID, dataType string) ([]model.Data, error)
	SelectByID(ctx context.Context, userID, dataType, dataID string) (model.Data, error)
}

// CryptService defines methods for cryptographic operations.
type CryptService interface {
	Encrypt(key, data []byte) ([]byte, error)
	Decrypt(key, data []byte) ([]byte, error)
	GenerateKey() ([]byte, error)
}

// BinaryDataService handles operations related to binary data.
type BinaryDataService struct {
	repository DataRepository         // The repository to store/retrieve data
	crypt      CryptService           // The cryptographic service for encrypting/decrypting data
	jwtManager *jwtmanager.JWTManager // JWT manager for handling authentication
	dataType   string                 // The type of data being handled (binary_data)
}

// New creates a new BinaryDataService instance.
func New(repository DataRepository, crypt CryptService, jwtManager *jwtmanager.JWTManager) *BinaryDataService {
	return &BinaryDataService{
		repository: repository,
		crypt:      crypt,
		jwtManager: jwtManager,
		dataType:   binaryData,
	}
}

// SaveBinaryData saves a new binary data entry after encrypting it.
func (s *BinaryDataService) SaveBinaryData(ctx context.Context, req model.BinaryDataPostRequest) (model.BinaryData, error) {
	userID, ok := ctx.Value(model.UserIDKey).(string)
	if !ok {
		return model.BinaryData{}, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(model.UserKey).([]byte)
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
	}, nil
}

// LoadAllBinaryInfo retrieves metadata for all binary data entries for the user.
func (s *BinaryDataService) LoadAllBinaryInfo(ctx context.Context) ([]model.DataInfo, error) {
	userID, ok := ctx.Value(model.UserIDKey).(string)
	if !ok {
		return nil, fmt.Errorf("failed to get userID from context")
	}

	encryptedBinaryData, err := s.repository.SelectAll(ctx, userID, s.dataType)
	if err != nil {
		return nil, fmt.Errorf("select all binary_data: %w", err)
	}

	binaryDataInfo := make([]model.DataInfo, 0, len(encryptedBinaryData))
	for _, encryptedBinary := range encryptedBinaryData {
		binaryDataInfo = append(binaryDataInfo, model.DataInfo{
			ID:        encryptedBinary.ID,
			DataType:  s.dataType,
			MetaData:  encryptedBinary.MetaData,
			CreatedAt: encryptedBinary.CreatedAt,
		})
	}
	return binaryDataInfo, nil
}

// LoadBinaryData retrieves and decrypts a binary data entry by its ID.
func (s *BinaryDataService) LoadBinaryData(ctx context.Context, dataID string) (model.BinaryData, error) {
	userID, ok := ctx.Value(model.UserIDKey).(string)
	if !ok {
		return model.BinaryData{}, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(model.UserKey).([]byte)
	if !ok {
		return model.BinaryData{}, fmt.Errorf("failed to get userKey from context")
	}

	encryptedBinaryData, err := s.repository.SelectByID(ctx, userID, s.dataType, dataID)
	if err != nil {
		return model.BinaryData{}, fmt.Errorf("select all binary_data: %w", err)
	}
	decryptedData, err := s.crypt.Decrypt(userKey, encryptedBinaryData.Data)
	if err != nil {
		return model.BinaryData{}, fmt.Errorf("decrypt binary: %w", err)
	}

	var decryptedBinaryData model.BinaryCryptData
	err = json.Unmarshal(decryptedData, &decryptedBinaryData)
	if err != nil {
		return model.BinaryData{}, fmt.Errorf("unmarshal binary: %w", err)
	}

	binary := model.BinaryData{
		ID:        encryptedBinaryData.ID,
		OwnerID:   encryptedBinaryData.OwnerID,
		Name:      decryptedBinaryData.Name,
		Extension: decryptedBinaryData.Extension,
		Data:      decryptedBinaryData.Data,
		MetaData:  encryptedBinaryData.MetaData,
		CreatedAt: encryptedBinaryData.CreatedAt,
	}

	return binary, nil
}

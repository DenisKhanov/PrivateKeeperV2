package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/DenisKhanov/PrivateKeeperV2/pkg/jwtmanager"
)

const credentials = "credentials" // Constant to define the data type for credentials

// DataRepository defines the methods for data operations on the repository level.
type DataRepository interface {
	Insert(ctx context.Context, data model.Data) (model.Data, error)
	SelectAll(ctx context.Context, userID, dataType string) ([]model.Data, error)
	SelectByID(ctx context.Context, userID, dataType, dataID string) (model.Data, error)
}

// CryptService defines methods for encryption and decryption operations.
type CryptService interface {
	Encrypt(key, data []byte) ([]byte, error)
	Decrypt(key, data []byte) ([]byte, error)
	GenerateKey() ([]byte, error)
}

// CredentialsService is responsible for handling operations related to user credentials.
type CredentialsService struct {
	repository DataRepository         // The repository for data operations
	crypt      CryptService           // The service for cryptographic operations
	jwtManager *jwtmanager.JWTManager // The JWT manager for token operations
	dataType   string                 // The type of data this service handles
}

// New creates a new instance of CredentialsService with the provided dependencies.
func New(repository DataRepository, crypt CryptService, jwtManager *jwtmanager.JWTManager) *CredentialsService {
	return &CredentialsService{
		repository: repository,
		crypt:      crypt,
		jwtManager: jwtManager,
		dataType:   credentials,
	}
}

// SaveCredentials saves the provided credentials to the repository after encrypting them.
// It returns the saved Credentials object or an error if the operation fails.
func (s *CredentialsService) SaveCredentials(ctx context.Context, req model.CredentialsPostRequest) (model.Credentials, error) {
	userID, ok := ctx.Value(model.UserIDKey).(string)
	if !ok {
		return model.Credentials{}, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(model.UserKey).([]byte)
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
	}, nil
}

// LoadAllCredentialsDataInfo retrieves all credentials data information for the user.
func (s *CredentialsService) LoadAllCredentialsDataInfo(ctx context.Context) ([]model.DataInfo, error) {
	userID, ok := ctx.Value(model.UserIDKey).(string)
	logrus.Info("UserID", userID)
	if !ok {
		return nil, fmt.Errorf("failed to get userID from context")
	}

	encryptedCredentialsData, err := s.repository.SelectAll(ctx, userID, s.dataType)
	if err != nil {
		return nil, fmt.Errorf("select all credentials_data: %w", err)
	}

	credentialsDataInfo := make([]model.DataInfo, 0, len(encryptedCredentialsData))
	for _, encryptedCredentials := range encryptedCredentialsData {
		credentialsDataInfo = append(credentialsDataInfo, model.DataInfo{
			ID:        encryptedCredentials.ID,
			DataType:  s.dataType,
			MetaData:  encryptedCredentials.MetaData,
			CreatedAt: encryptedCredentials.CreatedAt,
		})
	}
	return credentialsDataInfo, nil
}

// LoadCredentialsData retrieves and decrypts the credentials data for a given dataID.
func (s *CredentialsService) LoadCredentialsData(ctx context.Context, dataID string) (model.Credentials, error) {
	userID, ok := ctx.Value(model.UserIDKey).(string)
	if !ok {
		return model.Credentials{}, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(model.UserKey).([]byte)
	if !ok {
		return model.Credentials{}, fmt.Errorf("failed to get userKey from context")
	}

	encryptedBinaryData, err := s.repository.SelectByID(ctx, userID, s.dataType, dataID)
	if err != nil {
		return model.Credentials{}, fmt.Errorf("select all credentials_data: %w", err)
	}
	decryptedData, err := s.crypt.Decrypt(userKey, encryptedBinaryData.Data)
	if err != nil {
		return model.Credentials{}, fmt.Errorf("decrypt credentials: %w", err)
	}

	var decryptedCredentialsData model.CredentialsCryptData
	err = json.Unmarshal(decryptedData, &decryptedCredentialsData)
	if err != nil {
		return model.Credentials{}, fmt.Errorf("unmarshal cred: %w", err)
	}

	cred := model.Credentials{
		ID:        encryptedBinaryData.ID,
		OwnerID:   encryptedBinaryData.OwnerID,
		Login:     decryptedCredentialsData.Login,
		Password:  decryptedCredentialsData.Password,
		MetaData:  encryptedBinaryData.MetaData,
		CreatedAt: encryptedBinaryData.CreatedAt,
	}

	return cred, nil
}

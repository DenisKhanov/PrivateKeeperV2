package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/credentials/specification"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/auth"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/keyextraction"
	"github.com/google/uuid"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/DenisKhanov/PrivateKeeperV2/pkg/jwtmanager"
)

const credentials = "credentials"

type DataRepository interface {
	Insert(ctx context.Context, data model.Data) (model.Data, error)
	SelectAll(ctx context.Context, userID, dataType string) ([]model.Data, error)
}

type CryptService interface {
	Encrypt(key, data []byte) ([]byte, error)
	Decrypt(key, data []byte) ([]byte, error)
	GenerateKey() ([]byte, error)
}

type CredentialsService struct {
	repository DataRepository
	crypt      CryptService
	jwtManager *jwtmanager.JWTManager
	dataType   string
}

func New(repository DataRepository, crypt CryptService, jwtManager *jwtmanager.JWTManager) *CredentialsService {
	return &CredentialsService{
		repository: repository,
		crypt:      crypt,
		jwtManager: jwtManager,
		dataType:   credentials,
	}
}

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

func (s *CredentialsService) LoadAllCredentials(ctx context.Context, spec specification.CredentialsSpecification) ([]model.Credentials, error) {
	userID, ok := ctx.Value(auth.UserIDContextKey("userID")).(string)
	if !ok {
		return nil, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(keyextraction.UserKeyContextKey("userKey")).([]byte)
	if !ok {
		return nil, fmt.Errorf("failed to get userKey from context")
	}

	encryptedCredentials, err := s.repository.SelectAll(ctx, userID, s.dataType)
	if err != nil {
		return nil, fmt.Errorf("select all credentials: %w", err)
	}

	creds := make([]model.Credentials, 0, len(encryptedCredentials))
	for _, encryptedCred := range encryptedCredentials {
		decryptedData, err := s.crypt.Decrypt(userKey, encryptedCred.Data)
		if err != nil {
			return nil, fmt.Errorf("decrypt data: %w", err)
		}

		var cred model.CredentialsCryptData
		err = json.Unmarshal(decryptedData, &cred)
		if err != nil {
			return nil, fmt.Errorf("unmarshal data: %w", err)
		}

		creds = append(creds, model.Credentials{
			ID:        encryptedCred.ID,
			OwnerID:   encryptedCred.OwnerID,
			Login:     cred.Login,
			Password:  cred.Password,
			MetaData:  encryptedCred.MetaData,
			CreatedAt: encryptedCred.CreatedAt,
			UpdatedAt: encryptedCred.UpdatedAt,
		})
	}

	predicates := spec.MakeFilterPredicates()
	var filteredCreds []model.Credentials
	for _, cred := range creds {
		take := true
		for _, filterCardWithSpec := range predicates {
			if !filterCardWithSpec(spec, cred) {
				take = false
				break
			}
		}
		if take {
			filteredCreds = append(filteredCreds, cred)
		}
	}

	return filteredCreds, nil
}

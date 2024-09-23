package service

import (
	"context"

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

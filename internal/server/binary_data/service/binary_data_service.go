package service

import (
	"context"

	"github.com/DenisKhanov/PrivateKeeperV2/pkg/jwtmanager"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

const binaryData = "binary_data"

type DataRepository interface {
	Insert(ctx context.Context, data model.Data) (model.Data, error)
	SelectAll(ctx context.Context, userID, dataType string) ([]model.Data, error)
}

type CryptService interface {
	Encrypt(key, data []byte) ([]byte, error)
	Decrypt(key, data []byte) ([]byte, error)
	GenerateKey() ([]byte, error)
}

type BinaryDataService struct {
	repository DataRepository
	crypt      CryptService
	jwtManager *jwtmanager.JWTManager
	dataType   string
}

func New(repository DataRepository, crypt CryptService, jwtManager *jwtmanager.JWTManager) *BinaryDataService {
	return &BinaryDataService{
		repository: repository,
		crypt:      crypt,
		jwtManager: jwtManager,
		dataType:   binaryData,
	}
}

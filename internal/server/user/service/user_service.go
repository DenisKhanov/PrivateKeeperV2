package service

import (
	"context"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/cache"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/DenisKhanov/PrivateKeeperV2/pkg/jwtmanager"
)

type UserRepository interface {
	Insert(ctx context.Context, user model.User) (model.User, error)
	SelectByLogin(ctx context.Context, login string) (model.User, error)
}

type CryptService interface {
	EncryptWithMasterKey(data []byte) ([]byte, error)
	DecryptWithMasterKey(data []byte) ([]byte, error)
	GenerateKey() ([]byte, error)
}

type UserService struct {
	repository UserRepository
	crypt      CryptService
	jwtManager *jwtmanager.JWTManager
	redis      *cache.Redis
}

func New(repository UserRepository, crypt CryptService, jwtManager *jwtmanager.JWTManager, redis *cache.Redis) *UserService {
	return &UserService{
		repository: repository,
		crypt:      crypt,
		jwtManager: jwtManager,
		redis:      redis,
	}
}

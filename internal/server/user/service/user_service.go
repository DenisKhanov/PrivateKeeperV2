package service

import (
	"context"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/user/cerrors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"

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

func (u *UserService) Login(ctx context.Context, req model.UserLoginRequest) (string, error) {
	user, err := u.repository.SelectByLogin(ctx, req.Login)
	if err != nil {
		return "", fmt.Errorf("login SelectByLogin: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password))
	if err != nil {
		return "", cerrors.ErrInvalidPassword
	}

	token, err := u.jwtManager.BuildJWTString(user.ID)
	if err != nil {
		return "", fmt.Errorf("login build jwt: %w", err)
	}

	userKey, err := u.crypt.DecryptWithMasterKey(user.CryptKey)
	if err != nil {
		return "", fmt.Errorf("decryptWithMasterKey: %w", err)
	}

	st := u.redis.Client.Set(ctx, user.ID, userKey, 24*time.Hour)
	if st.Err() != nil {
		return "", fmt.Errorf("login redis set: %w", st.Err())
	}

	return token, nil
}

func (u *UserService) Register(ctx context.Context, req model.UserRegisterRequest) (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("new uuid: %w", err)
	}

	userKey, err := u.crypt.GenerateKey()
	if err != nil {
		return "", fmt.Errorf("genereate key: %w", err)
	}

	cryptUserKey, err := u.crypt.EncryptWithMasterKey(userKey)
	if err != nil {
		return "", fmt.Errorf("encrypt with master key: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate hash from password: %w", err)
	}

	userToSave := model.User{
		ID:       id.String(),
		Login:    req.Login,
		Password: passwordHash,
		CryptKey: cryptUserKey,
	}

	user, err := u.repository.Insert(ctx, userToSave)
	if err != nil {
		return "", fmt.Errorf("register user: %w", err)
	}

	st := u.redis.Client.Set(ctx, user.ID, userKey, 24*time.Hour)
	if st.Err() != nil {
		return "", fmt.Errorf("redis set: %w", st.Err())
	}

	token, err := u.jwtManager.BuildJWTString(user.ID)
	if err != nil {
		return "", fmt.Errorf("build jwt: %w", err)
	}

	return token, nil
}

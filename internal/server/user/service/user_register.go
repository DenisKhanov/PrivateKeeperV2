package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

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

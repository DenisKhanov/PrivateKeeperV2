package keyextraction

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/cache"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/auth"
)

var userKeyExtractorMandatoryMethods = map[string]struct{}{
	"/proto.CreditCardService/PostSaveCreditCard":   {},
	"/proto.CreditCardService/GetLoadCreditCard":    {},
	"/proto.TextDataService/PostSaveTextData":       {},
	"/proto.TextDataService/GetLoadTextData":        {},
	"/proto.BinaryDataService/PostSaveBinaryData":   {},
	"/proto.BinaryDataService/GetLoadBinaryData":    {},
	"/proto.CredentialsService/PostSaveCredentials": {},
	"/proto.CredentialsService/GetLoadCredentials":  {},
}

type CryptService interface {
	DecryptWithMasterKey(data []byte) ([]byte, error)
}

type UserRepository interface {
	SelectKeyByID(ctx context.Context, userID string) ([]byte, error)
}

type UserKeyContextKey string

type UserKeyExtraction struct {
	cryptService CryptService
	userRepo     UserRepository
	redis        *cache.Redis
}

func New(service CryptService, repository UserRepository, redis *cache.Redis) *UserKeyExtraction {
	return &UserKeyExtraction{
		cryptService: service,
		userRepo:     repository,
		redis:        redis,
	}
}

func (j *UserKeyExtraction) ExtractUserKey(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if _, ok := userKeyExtractorMandatoryMethods[info.FullMethod]; !ok {
		return handler(ctx, req)
	}

	userID, ok := ctx.Value(auth.UserIDContextKey("userID")).(string)
	if !ok {
		logrus.Error("Unable to extract user key: failed to get user id from context")
		return nil, status.Error(codes.Internal, "internal error")
	}

	var key []byte
	key, err := j.redis.Client.Get(ctx, userID).Bytes()
	if err != nil {
		logrus.WithError(err).Error("Unable to extract user key: failed to get user key from redis")
		cryptKey, err := j.userRepo.SelectKeyByID(ctx, userID)
		if err != nil {
			logrus.WithError(err).Error("Unable to extract user key: failed to get user_id from db")
			return nil, status.Error(codes.Internal, "internal error")
		}

		key, err = j.cryptService.DecryptWithMasterKey(cryptKey)
		if err != nil {
			logrus.WithError(err).Error("Unable to extract user key: failed to decrypt user key")
			return nil, status.Error(codes.Internal, "internal error")
		}
		st := j.redis.Client.Set(ctx, userID, key, 24*time.Hour)
		if st.Err() != nil {
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	ctx = context.WithValue(ctx, UserKeyContextKey("userKey"), key)
	return handler(ctx, req)
}

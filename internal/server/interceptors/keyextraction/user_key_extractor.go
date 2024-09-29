package keyextraction

import (
	"context"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/sirupsen/logrus"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/cache"
)

// Define a map of methods that require user key extraction userKeyExtractorMandatoryMethods.
var userKeyExtractorMandatoryMethods = map[string]struct{}{
	"/proto.CreditCardService/PostSaveCreditCard":             {},
	"/proto.CreditCardService/GetLoadCreditCard":              {},
	"/proto.CreditCardService/GetLoadAllCreditCardDataInfo":   {},
	"/proto.TextDataService/PostSaveTextData":                 {},
	"/proto.TextDataService/GetLoadTextData":                  {},
	"/proto.TextDataService/GetLoadAllTextDataInfo":           {},
	"/proto.BinaryDataService/PostSaveBinaryData":             {},
	"/proto.BinaryDataService/GetLoadBinaryData":              {},
	"/proto.BinaryDataService/GetLoadAllBinaryDataInfo":       {},
	"/proto.CredentialsService/PostSaveCredentials":           {},
	"/proto.CredentialsService/GetLoadCredentials":            {},
	"/proto.CredentialsService/GetLoadAllCredentialsDataInfo": {},
}

// CryptService interface defines the method for decrypting data with a master key.
type CryptService interface {
	DecryptWithMasterKey(data []byte) ([]byte, error)
}

// UserRepository interface defines the method for fetching user keys from a repository.
type UserRepository interface {
	SelectKeyByID(ctx context.Context, userID string) ([]byte, error)
}

// UserKeyExtraction struct handles the extraction of user keys.
type UserKeyExtraction struct {
	cryptService CryptService   // Instance of CryptService for decryption
	userRepo     UserRepository // Instance of UserRepository for fetching user keys
	redis        *cache.Redis   // Instance of Redis for caching user keys
}

// New creates a new instance of UserKeyExtraction.
func New(service CryptService, repository UserRepository, redis *cache.Redis) *UserKeyExtraction {
	return &UserKeyExtraction{
		cryptService: service,
		userRepo:     repository,
		redis:        redis,
	}
}

// ExtractUserKey checks if user key extraction is needed and retrieves the user key.
// It handles the logic for getting the key from Redis, database, and decrypting if necessary.
func (j *UserKeyExtraction) ExtractUserKey(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if _, ok := userKeyExtractorMandatoryMethods[info.FullMethod]; !ok {
		return handler(ctx, req)
	}

	userID, ok := ctx.Value(model.UserIDKey).(string)
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

	ctx = context.WithValue(ctx, model.UserKey, key)
	return handler(ctx, req)
}

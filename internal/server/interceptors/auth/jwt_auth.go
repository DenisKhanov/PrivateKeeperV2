package auth

import (
	"context"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/sirupsen/logrus"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/DenisKhanov/PrivateKeeperV2/pkg/jwtmanager"
)

// Define a map of methods that require authentication authMandatoryMethods.
var authMandatoryMethods = map[string]struct{}{
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

// JWTAuth struct holds the JWT manager for authentication.
type JWTAuth struct {
	jwtManager *jwtmanager.JWTManager // Instance of JWTManager for token handling
}

// New creates a new instance of JWTAuth with the provided JWT manager.
func New(jwtManager *jwtmanager.JWTManager) *JWTAuth {
	return &JWTAuth{jwtManager: jwtManager}
}

// GRPCJWTAuth checks token from gRPC metadata and sets userID in the context.
// If authentication fails, it returns an error with the corresponding status code.
func (j *JWTAuth) GRPCJWTAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if _, ok := authMandatoryMethods[info.FullMethod]; !ok {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logrus.Info("Authentication failed: missing metadata")
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	c := md.Get(j.jwtManager.TokenName)
	if len(c) < 1 {
		logrus.Info("Authentication failed: token not found")
		return nil, status.Errorf(codes.Unauthenticated, "token not found")
	}

	userID, err := j.jwtManager.GetUserID(c[0])
	if err != nil {
		logrus.Info("Authentication failed: unable to get userID from token", slog.String("error", err.Error()))
		return nil, status.Errorf(codes.Unauthenticated, "authentification by UserID failed")
	}
	logrus.Info("Authentication succeeded UserId is: ", userID)
	ctx = context.WithValue(ctx, model.UserIDKey, userID)
	return handler(ctx, req)
}

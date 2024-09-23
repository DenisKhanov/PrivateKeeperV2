package grpchandlers

import (
	"context"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/user"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

type UserService interface {
	Register(ctx context.Context, user model.UserRegisterRequest) (string, error)
	Login(ctx context.Context, user model.UserLoginRequest) (string, error)
}

type Validator interface {
	ValidateLoginRequest(req *model.UserLoginRequest) (map[string]string, bool)
	ValidateRegisterRequest(req *model.UserRegisterRequest) (map[string]string, bool)
}

type UserHandler struct {
	userService UserService
	pb.UnimplementedUserServiceServer
	validator Validator
}

func New(userService UserService, validator Validator) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
	}
}

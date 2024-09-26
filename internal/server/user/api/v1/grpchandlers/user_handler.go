package grpchandlers

import (
	"context"
	"errors"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/lib"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/user/cerrors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func (h *UserHandler) PostRegisterUser(ctx context.Context, in *pb.PostUserRegisterRequest) (*pb.PostUserRegisterResponse, error) {
	req := model.UserRegisterRequest{
		Login:    in.Login,
		Password: in.Password,
	}

	report, ok := h.validator.ValidateRegisterRequest(&req)
	if !ok {
		logrus.Info("Unable to register user: invalid user request")
		logrus.Info("user_login", req.Login)
		logrus.Infof("violated_fields %v", report)
		return nil, lib.ProcessValidationError("invalid user request", report)
	}

	token, err := h.userService.Register(ctx, req)
	if errors.Is(err, cerrors.ErrUserAlreadyExists) {
		logrus.Info("Unable to register user: user already exists")
		logrus.Infof("user_login %v", req.Login)
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}

	if err != nil {
		logrus.WithError(err).Error("Unable to register user")
		logrus.Infof("user_login %v", req.Login)

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.PostUserRegisterResponse{Token: token}, nil
}

func (h *UserHandler) PostLoginUser(ctx context.Context, in *pb.PostUserLoginRequest) (*pb.PostUserLoginResponse, error) {
	req := model.UserLoginRequest{
		Login:    in.Login,
		Password: in.Password,
	}

	report, ok := h.validator.ValidateLoginRequest(&req)
	if !ok {
		logrus.Info("Unable to login user: invalid user request")
		logrus.Infof("user_login %v", req.Login)
		logrus.Infof("violated_fields %v", report)
		return nil, lib.ProcessValidationError("invalid user request", report)
	}

	token, err := h.userService.Login(ctx, req)
	if errors.Is(err, cerrors.ErrUserNotFound) {
		logrus.Info("Unable to login user: user not found")
		logrus.Infof("user_login %v", req.Login)
		return nil, status.Error(codes.NotFound, "user with this login not found")
	}

	if errors.Is(err, cerrors.ErrInvalidPassword) {
		logrus.Info("Unable to login user: invalid password")
		logrus.Infof("user_login %v", req.Login)
		return nil, status.Error(codes.Unauthenticated, "incorrect password")
	}

	if err != nil {
		logrus.WithError(err).Error("Unable to login user: invalid user request")
		logrus.Infof("user_login %v", req.Login)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.PostUserLoginResponse{Token: token}, nil
}

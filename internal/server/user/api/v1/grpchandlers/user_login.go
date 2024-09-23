package grpchandlers

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/user"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/lib"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/user/cerrors"
)

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

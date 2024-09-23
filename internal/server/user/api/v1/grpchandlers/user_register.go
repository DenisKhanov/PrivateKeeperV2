package grpchandlers

import (
	"context"
	"errors"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/lib"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/user"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/user/cerrors"
)

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

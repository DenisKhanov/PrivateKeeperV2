package grpchandlers

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credentials"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/lib"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

func (h *CredentialsHandler) PostSaveCredentials(ctx context.Context, in *pb.PostCredentialsRequest) (*pb.PostCredentialsResponse, error) {
	req := model.CredentialsPostRequest{
		Login:    in.Login,
		Password: in.Password,
		MetaData: in.Metadata,
	}

	report, ok := h.validator.ValidatePostRequest(&req)
	if !ok {
		logrus.Info("Unable to register user: invalid credentials request")
		logrus.Infof("violated_fields %v", report)
		return nil, lib.ProcessValidationError("invalid credentials post request", report)
	}

	cred, err := h.credentialsService.SaveCredentials(ctx, req)
	if err != nil {
		logrus.WithError(err).Error("Unable to save credentials")
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.PostCredentialsResponse{
		Id:        cred.ID,
		Login:     cred.Login,
		Password:  cred.Password,
		Metadata:  cred.MetaData,
		CreatedAt: cred.CreatedAt.Format(time.RFC3339),
		UpdatedAt: cred.UpdatedAt.Format(time.RFC3339),
	}, nil
}

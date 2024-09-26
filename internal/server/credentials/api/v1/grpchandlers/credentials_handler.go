package grpchandlers

import (
	"context"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/lib"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credentials"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/credentials/specification"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

type CredentialsService interface {
	SaveCredentials(ctx context.Context, req model.CredentialsPostRequest) (model.Credentials, error)
	LoadAllCredentials(ctx context.Context, spec specification.CredentialsSpecification) ([]model.Credentials, error)
}

type Validator interface {
	ValidatePostRequest(req *model.CredentialsPostRequest) (map[string]string, bool)
}

type CredentialsHandler struct {
	credentialsService CredentialsService
	pb.UnimplementedCredentialsServiceServer
	validator Validator
}

func New(textDataService CredentialsService, validator Validator) *CredentialsHandler {
	return &CredentialsHandler{
		credentialsService: textDataService,
		validator:          validator,
	}
}

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

func (h *CredentialsHandler) GetLoadCredentials(ctx context.Context, in *pb.GetCredentialsRequest) (*pb.GetCredentialsResponse, error) {
	spec, err := specification.NewCredentialsSpecification(in)
	if err != nil {
		logrus.WithError(err).Error("Error while creating credentials specification: ")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	creds, err := h.credentialsService.LoadAllCredentials(ctx, spec)
	if err != nil {
		logrus.WithError(err).Error("Error while loading credentials: ")
		return nil, status.Error(codes.Internal, "internal error")
	}

	credsData := make([]*pb.Credentials, 0, len(creds))
	for _, v := range creds {
		credsData = append(credsData, &pb.Credentials{
			Id:        v.ID,
			OwnerId:   v.OwnerID,
			Login:     v.Login,
			Password:  v.Password,
			Metadata:  v.MetaData,
			CreatedAt: v.CreatedAt.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.GetCredentialsResponse{Creds: credsData}, nil
}

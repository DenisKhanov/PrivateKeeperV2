package grpchandlers

import (
	"context"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/lib"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credentials"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

// CredentialsService is an interface that defines the methods for handling credentials.
// It encapsulates the logic for saving and loading credentials data.
type CredentialsService interface {
	SaveCredentials(ctx context.Context, req model.CredentialsPostRequest) (model.Credentials, error)
	LoadCredentialsData(ctx context.Context, dataID string) (model.Credentials, error)
	LoadAllCredentialsDataInfo(ctx context.Context) ([]model.DataInfo, error)
}

// Validator is an interface for validating incoming requests.
type Validator interface {
	ValidatePostRequest(req *model.CredentialsPostRequest) (map[string]string, bool)
}

// CredentialsHandler implements the gRPC server for credentials-related operations.
// It holds references to the credentials service and validator.
type CredentialsHandler struct {
	credentialsService                       CredentialsService // Service for handling credentials
	pb.UnimplementedCredentialsServiceServer                    // Embed the unimplemented server for compliance with the gRPC server interface
	validator                                Validator          // Validator for incoming requests
}

// New creates a new instance of CredentialsHandler with the provided services.
func New(credentialsService CredentialsService, validator Validator) *CredentialsHandler {
	return &CredentialsHandler{
		credentialsService: credentialsService,
		validator:          validator,
	}
}

// PostSaveCredentials handles the gRPC call to save credentials.
// It processes the incoming request, validates it, and invokes the service to save the data.
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
	}, nil
}

// GetLoadAllCredentialsDataInfo handles the gRPC call to load all credentials data information.
// It retrieves metadata for all credentials and constructs a response.
func (h *CredentialsHandler) GetLoadAllCredentialsDataInfo(ctx context.Context, _ *pb.GetAllCredentialsInfoRequest) (*pb.GetAllCredentialsInfoResponse, error) {

	credentialsInfo, err := h.credentialsService.LoadAllCredentialsDataInfo(ctx)
	if err != nil {
		logrus.WithError(err).Error("Error while loading credentials data: ")
		return nil, status.Error(codes.Internal, "internal error")
	}

	credentialsInfos := make([]*pb.CredentialsInfo, 0, len(credentialsInfo))
	for _, v := range credentialsInfo {
		credentialsInfos = append(credentialsInfos, &pb.CredentialsInfo{
			Id:        v.ID,
			DataType:  v.DataType,
			Metadata:  v.MetaData,
			CreatedAt: v.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.GetAllCredentialsInfoResponse{Creds: credentialsInfos}, nil
}

// GetLoadCredentials handles the gRPC call to load specific credentials data by ID.
// It retrieves and returns the credentials based on the provided data ID.
func (h *CredentialsHandler) GetLoadCredentials(ctx context.Context, in *pb.GetCredentialsRequest) (*pb.GetCredentialsResponse, error) {
	dataID := in.Id

	credentialsData, err := h.credentialsService.LoadCredentialsData(ctx, dataID)
	if err != nil {
		logrus.WithError(err).Error("Error while loading credentials data: ")
		return nil, status.Error(codes.Internal, "internal error")
	}

	bin := &pb.Credentials{
		Id:        credentialsData.ID,
		OwnerId:   credentialsData.OwnerID,
		Login:     credentialsData.Login,
		Password:  credentialsData.Password,
		Metadata:  credentialsData.MetaData,
		CreatedAt: credentialsData.CreatedAt.Format(time.RFC3339Nano),
	}
	return &pb.GetCredentialsResponse{CredentialsData: bin}, nil
}

package grpchandlers

import (
	"context"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/lib"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/binary_data"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

// BinaryDataService defines the methods for handling binary data operations.
type BinaryDataService interface {
	SaveBinaryData(ctx context.Context, req model.BinaryDataPostRequest) (model.BinaryData, error)
	LoadBinaryData(ctx context.Context, dataID string) (model.BinaryData, error)
	LoadAllBinaryInfo(ctx context.Context) ([]model.DataInfo, error)
}

// Validator defines the method for validating binary data post requests.
type Validator interface {
	ValidatePostRequest(req *model.BinaryDataPostRequest) (map[string]string, bool)
}

// BinaryDataHandler implements the gRPC server for handling binary data requests.
type BinaryDataHandler struct {
	binaryDataService                       BinaryDataService // Service for binary data operations
	pb.UnimplementedBinaryDataServiceServer                   // Embed the unimplemented server to provide backward compatibility
	validator                               Validator         // Validator for request validation
}

// New creates a new instance of BinaryDataHandler.
func New(binaryDataService BinaryDataService, validator Validator) *BinaryDataHandler {
	return &BinaryDataHandler{
		binaryDataService: binaryDataService,
		validator:         validator,
	}
}

// PostSaveBinaryData handles the gRPC request for saving binary data.
func (h *BinaryDataHandler) PostSaveBinaryData(ctx context.Context, in *pb.PostBinaryDataRequest) (*pb.PostBinaryDataResponse, error) {
	req := model.BinaryDataPostRequest{
		Name:      in.Name,
		Extension: in.Extension,
		Data:      in.Data,
		MetaData:  in.Metadata,
	}

	report, ok := h.validator.ValidatePostRequest(&req)
	if !ok {
		logrus.Info("Unable to register user: invalid user request")
		logrus.Infof("violated_fields %v", report)
		return nil, lib.ProcessValidationError("invalid text_data post request", report)
	}

	binary, err := h.binaryDataService.SaveBinaryData(ctx, req)
	if err != nil {
		logrus.WithError(err).Error("Unable to save binary_data")
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.PostBinaryDataResponse{
		Id:        binary.ID,
		Name:      binary.Name,
		Extension: binary.Extension,
		Metadata:  binary.MetaData,
		CreatedAt: binary.CreatedAt.Format(time.RFC3339Nano),
	}, nil
}

// GetLoadAllBinaryDataInfo handles the gRPC request for loading all binary data information.
func (h *BinaryDataHandler) GetLoadAllBinaryDataInfo(ctx context.Context, _ *pb.GetAllBinaryInfoRequest) (*pb.GetAllBinaryInfoResponse, error) {

	binariesInfo, err := h.binaryDataService.LoadAllBinaryInfo(ctx)
	if err != nil {
		logrus.WithError(err).Error("Error while loading binary data: ")
		return nil, status.Error(codes.Internal, "internal error")
	}

	binaryInfos := make([]*pb.BinaryInfo, 0, len(binariesInfo))
	for _, v := range binariesInfo {
		binaryInfos = append(binaryInfos, &pb.BinaryInfo{
			Id:        v.ID,
			DataType:  v.DataType,
			Metadata:  v.MetaData,
			CreatedAt: v.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.GetAllBinaryInfoResponse{Binaries: binaryInfos}, nil
}

// GetLoadBinaryData handles the gRPC request for loading specific binary data.
func (h *BinaryDataHandler) GetLoadBinaryData(ctx context.Context, in *pb.GetBinaryDataRequest) (*pb.GetBinaryDataResponse, error) {
	dataID := in.Id

	binaryData, err := h.binaryDataService.LoadBinaryData(ctx, dataID)
	if err != nil {
		logrus.WithError(err).Error("Error while loading binary data: ")
		return nil, status.Error(codes.Internal, "internal error")
	}

	bin := &pb.BinaryData{
		Id:        binaryData.ID,
		OwnerId:   binaryData.OwnerID,
		Data:      binaryData.Data,
		Name:      binaryData.Name,
		Extension: binaryData.Extension,
		Metadata:  binaryData.MetaData,
		CreatedAt: binaryData.CreatedAt.Format(time.RFC3339Nano),
	}
	return &pb.GetBinaryDataResponse{BinaryData: bin}, nil
}

package grpchandlers

import (
	"context"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/lib"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/binary_data"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/binary_data/specification"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

type BinaryDataService interface {
	SaveBinaryData(ctx context.Context, req model.BinaryDataPostRequest) (model.BinaryData, error)
	LoadAllBinaryData(ctx context.Context, spec specification.BinaryDataSpecification) ([]model.BinaryData, error)
}

type Validator interface {
	ValidatePostRequest(req *model.BinaryDataPostRequest) (map[string]string, bool)
}

type BinaryDataHandler struct {
	binaryDataService BinaryDataService
	pb.UnimplementedBinaryDataServiceServer
	validator Validator
}

func New(binaryDataService BinaryDataService, validator Validator) *BinaryDataHandler {
	return &BinaryDataHandler{
		binaryDataService: binaryDataService,
		validator:         validator,
	}
}

func (h *BinaryDataHandler) GetLoadBinaryData(ctx context.Context, in *pb.GetBinaryDataRequest) (*pb.GetBinaryDataResponse, error) {
	spec, err := specification.NewTextDataSpecification(in)
	if err != nil {
		logrus.WithError(err).Error("Error while creating text data specification: ")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	binaries, err := h.binaryDataService.LoadAllBinaryData(ctx, spec)
	if err != nil {
		logrus.WithError(err).Error("Error while loading binary data: ")
		return nil, status.Error(codes.Internal, "internal error")
	}

	binaryData := make([]*pb.BinaryData, 0, len(binaries))
	for _, v := range binaries {
		binaryData = append(binaryData, &pb.BinaryData{
			Id:        v.ID,
			OwnerId:   v.OwnerID,
			Data:      v.Data,
			Name:      v.Name,
			Extension: v.Extension,
			Metadata:  v.MetaData,
			CreatedAt: v.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.GetBinaryDataResponse{Binaries: binaryData}, nil
}

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

package grpchandlers

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/binary_data"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/lib"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

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
		UpdatedAt: binary.UpdatedAt.Format(time.RFC3339Nano),
	}, nil
}

package grpchandlers

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/binary_data"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/binary_data/specification"
)

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
			UpdatedAt: v.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.GetBinaryDataResponse{Binaries: binaryData}, nil
}

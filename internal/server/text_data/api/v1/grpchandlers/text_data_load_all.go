package grpchandlers

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/text_data"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/text_data/specification"
)

func (h *TextDataHandler) GetLoadTextData(ctx context.Context, in *pb.GetTextDataRequest) (*pb.GetTextDataResponse, error) {
	spec, err := specification.NewTextDataSpecification(in)
	if err != nil {
		logrus.WithError(err).Error("Error while creating text data specification: ")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	texts, err := h.textDataService.LoadAllTextData(ctx, spec)
	if err != nil {
		logrus.WithError(err).Error("Error while loading text data: ")
		return nil, status.Error(codes.Internal, "internal error")
	}

	textData := make([]*pb.TextData, 0, len(texts))
	for _, v := range texts {
		textData = append(textData, &pb.TextData{
			Id:        v.ID,
			OwnerId:   v.OwnerID,
			Text:      v.Text,
			Metadata:  v.MetaData,
			CreatedAt: v.CreatedAt.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.GetTextDataResponse{Text: textData}, nil
}

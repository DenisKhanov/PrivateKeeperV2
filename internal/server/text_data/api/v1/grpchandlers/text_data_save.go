package grpchandlers

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/text_data"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/lib"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

func (h *TextDataHandler) PostSaveTextData(ctx context.Context, in *pb.PostTextDataRequest) (*pb.PostTextDataResponse, error) {
	req := model.TextDataPostRequest{
		Text:     in.Text,
		MetaData: in.Metadata,
	}

	report, ok := h.validator.ValidatePostRequest(&req)
	if !ok {
		logrus.Info("Unable to register user: invalid user request")
		logrus.Infof("violated_fields %v", report)
		return nil, lib.ProcessValidationError("invalid text_data post request", report)
	}

	text, err := h.textDataService.SaveTextData(ctx, req)
	if err != nil {
		logrus.WithError(err).Errorf("failed to save text_data")
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.PostTextDataResponse{
		Id:        text.ID,
		Text:      text.Text,
		Metadata:  text.MetaData,
		CreatedAt: text.CreatedAt.Format(time.RFC3339),
		UpdatedAt: text.UpdatedAt.Format(time.RFC3339),
	}, nil
}

package grpchandlers

import (
	"context"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/lib"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/text_data"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/text_data/specification"
)

type TextDataService interface {
	SaveTextData(ctx context.Context, req model.TextDataPostRequest) (model.TextData, error)
	LoadAllTextData(ctx context.Context, spec specification.TextDataSpecification) ([]model.TextData, error)
}

type Validator interface {
	ValidatePostRequest(req *model.TextDataPostRequest) (map[string]string, bool)
}

type TextDataHandler struct {
	textDataService TextDataService
	pb.UnimplementedTextDataServiceServer
	validator Validator
}

func New(textDataService TextDataService, validator Validator) *TextDataHandler {
	return &TextDataHandler{
		textDataService: textDataService,
		validator:       validator,
	}
}

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
	}, nil
}

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
		})
	}

	return &pb.GetTextDataResponse{Text: textData}, nil
}

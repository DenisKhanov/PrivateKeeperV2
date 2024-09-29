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
)

// TextDataService interface defines the methods for text data management.
type TextDataService interface {
	SaveTextData(ctx context.Context, req model.TextDataPostRequest) (model.TextData, error)
	LoadTextData(ctx context.Context, dataID string) (model.TextData, error)
	LoadAllTextInfo(ctx context.Context) ([]model.DataInfo, error)
}

// Validator interface defines the method for validating requests.
type Validator interface {
	ValidatePostRequest(req *model.TextDataPostRequest) (map[string]string, bool)
}

// TextDataHandler struct implements the gRPC handler for text data operations.
type TextDataHandler struct {
	textDataService TextDataService
	pb.UnimplementedTextDataServiceServer
	validator Validator
}

// New creates a new instance of TextDataHandler.
func New(textDataService TextDataService, validator Validator) *TextDataHandler {
	return &TextDataHandler{
		textDataService: textDataService,
		validator:       validator,
	}
}

// PostSaveTextData handles the gRPC request to save text data.
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

// GetLoadAllTextDataInfo handles the gRPC request to load all text data information.
func (h *TextDataHandler) GetLoadAllTextDataInfo(ctx context.Context, _ *pb.GetAllTextInfoRequest) (*pb.GetAllTextInfoResponse, error) {

	textInfo, err := h.textDataService.LoadAllTextInfo(ctx)
	if err != nil {
		logrus.WithError(err).Error("Error while loading text data: ")
		return nil, status.Error(codes.Internal, "internal error")
	}

	textInfos := make([]*pb.TextInfo, 0, len(textInfo))
	for _, v := range textInfo {
		textInfos = append(textInfos, &pb.TextInfo{
			Id:        v.ID,
			DataType:  v.DataType,
			Metadata:  v.MetaData,
			CreatedAt: v.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.GetAllTextInfoResponse{Text: textInfos}, nil
}

// GetLoadTextData handles the gRPC request to load text data by ID.
func (h *TextDataHandler) GetLoadTextData(ctx context.Context, in *pb.GetTextDataRequest) (*pb.GetTextDataResponse, error) {
	dataID := in.Id

	textData, err := h.textDataService.LoadTextData(ctx, dataID)
	if err != nil {
		logrus.WithError(err).Error("Error while loading text data: ")
		return nil, status.Error(codes.Internal, "internal error")
	}

	text := &pb.TextData{
		Id:        textData.ID,
		OwnerId:   textData.OwnerID,
		Text:      textData.Text,
		Metadata:  textData.MetaData,
		CreatedAt: textData.CreatedAt.Format(time.RFC3339Nano),
	}
	return &pb.GetTextDataResponse{TextData: text}, nil
}

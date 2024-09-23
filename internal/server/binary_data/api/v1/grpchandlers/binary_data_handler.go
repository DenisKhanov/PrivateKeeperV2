package grpchandlers

import (
	"context"

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

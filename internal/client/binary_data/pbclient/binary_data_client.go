package pbclient

import (
	"context"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/binary_data"
	"google.golang.org/grpc/metadata"
)

// BinaryDataPBClient is a client wrapper for interacting with the BinaryDataService.
type BinaryDataPBClient struct {
	binaryDataService pb.BinaryDataServiceClient
}

// NewBinaryDataPBClient creates a new BinaryDataPBClient with the given BinaryDataServiceClient.
func NewBinaryDataPBClient(u pb.BinaryDataServiceClient) *BinaryDataPBClient {
	return &BinaryDataPBClient{
		binaryDataService: u,
	}
}

// SaveBinaryData sends a request to save binary data and returns the saved data or an error.
func (u *BinaryDataPBClient) SaveBinaryData(ctx context.Context, token string, bData model.BinaryDataPostRequest) (model.BinaryData, error) {
	req := &pb.PostBinaryDataRequest{
		Data:      bData.Data,
		Name:      bData.Name,
		Extension: bData.Extension,
		Metadata:  bData.MetaData,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := u.binaryDataService.PostSaveBinaryData(ctx, req)
	if err != nil {
		return model.BinaryData{}, fmt.Errorf("save binary data: %w", err)
	}

	binaryData := model.BinaryData{
		Name:      resp.Name,
		Extension: resp.Extension,
		MetaData:  resp.Metadata,
	}

	return binaryData, nil
}

// LoadAllBinaryDataInfo retrieves all binary data info and returns a list of DataInfo or an error.
func (u *BinaryDataPBClient) LoadAllBinaryDataInfo(ctx context.Context, token string) ([]model.DataInfo, error) {
	req := &pb.GetAllBinaryInfoRequest{}

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := u.binaryDataService.GetLoadAllBinaryDataInfo(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("load binary data: %w", err)
	}

	binaries := make([]model.DataInfo, 0, len(resp.Binaries))
	for _, data := range resp.Binaries {
		binaries = append(binaries, model.DataInfo{
			ID:        data.Id,
			DataType:  data.DataType,
			MetaData:  data.Metadata,
			CreatedAt: data.CreatedAt,
		})
	}

	return binaries, nil
}

// LoadBinaryData retrieves binary data by ID and returns the BinaryData or an error.
func (u *BinaryDataPBClient) LoadBinaryData(ctx context.Context, token string, dataID string) (model.BinaryData, error) {
	req := &pb.GetBinaryDataRequest{
		Id: dataID,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := u.binaryDataService.GetLoadBinaryData(ctx, req)
	if err != nil {
		return model.BinaryData{}, fmt.Errorf("load binary data: %w", err)
	}
	data := resp.BinaryData
	binaryData := model.BinaryData{
		Name:      data.Name,
		Extension: data.Extension,
		Data:      data.Data,
		MetaData:  data.Metadata,
	}

	return binaryData, nil
}

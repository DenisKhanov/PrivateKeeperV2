package pbclient

import (
	"context"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/binary_data"
	"google.golang.org/grpc/metadata"
)

type BinaryDataPBClient struct {
	binaryDataService pb.BinaryDataServiceClient
}

func NewBinaryDataPBClient(u pb.BinaryDataServiceClient) *BinaryDataPBClient {
	return &BinaryDataPBClient{
		binaryDataService: u,
	}
}

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

func (u *BinaryDataPBClient) LoadBinaryData(ctx context.Context, token string, bData model.BinaryDataLoadRequest) ([]model.BinaryData, error) {
	req := &pb.GetBinaryDataRequest{
		Name:     bData.Name,
		Metadata: bData.MetaData,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := u.binaryDataService.GetLoadBinaryData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("load binary data: %w", err)
	}

	binaries := make([]model.BinaryData, 0, len(resp.Binaries))
	for _, data := range resp.Binaries {
		binaries = append(binaries, model.BinaryData{
			Name:      data.Name,
			Extension: data.Extension,
			Data:      data.Data,
			MetaData:  data.Metadata,
		})
	}

	return binaries, nil
}

package pbclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/binary_data"
)

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

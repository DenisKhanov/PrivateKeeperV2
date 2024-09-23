package pbclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/binary_data"
)

func (u *BinaryDataPBClient) LoadBinaryData(token string, bData model.BinaryDataLoadRequest) ([]model.BinaryData, error) {
	req := &pb.GetBinaryDataRequest{
		Name:     bData.Name,
		Metadata: bData.MetaData,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

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

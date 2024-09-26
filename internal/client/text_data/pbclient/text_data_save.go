package pbclient

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/text_data"
)

func (u *TextDataPBClient) SaveTextData(ctx context.Context, token string, text model.TextDataPostRequest) (model.TextData, error) {
	req := &pb.PostTextDataRequest{
		Text:     text.Text,
		Metadata: text.MetaData,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := u.textDataService.PostSaveTextData(ctx, req)
	if err != nil {
		return model.TextData{}, err
	}

	txt := model.TextData{
		Text:     resp.Text,
		MetaData: resp.Metadata,
	}

	return txt, nil
}

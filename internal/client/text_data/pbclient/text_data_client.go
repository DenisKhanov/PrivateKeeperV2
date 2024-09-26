package pbclient

import (
	"context"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/text_data"
	"google.golang.org/grpc/metadata"
)

type TextDataPBClient struct {
	textDataService pb.TextDataServiceClient
}

func NewCreditCardPBClient(u pb.TextDataServiceClient) *TextDataPBClient {
	return &TextDataPBClient{
		textDataService: u,
	}
}

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

func (u *TextDataPBClient) LoadTextData(ctx context.Context, token string, textData model.TextDataLoadRequest) ([]model.TextData, error) {
	req := &pb.GetTextDataRequest{
		Text:     textData.Text,
		Metadata: textData.MetaData,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := u.textDataService.GetLoadTextData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("load text data: %w", err)
	}

	texts := make([]model.TextData, 0, len(resp.Text))
	for _, txt := range resp.Text {
		texts = append(texts, model.TextData{
			Text:     txt.Text,
			MetaData: txt.Metadata,
		})
	}

	return texts, nil
}

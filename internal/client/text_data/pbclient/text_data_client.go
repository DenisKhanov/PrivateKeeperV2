package pbclient

import (
	"context"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/text_data"
	"google.golang.org/grpc/metadata"
)

// TextDataPBClient is a client for interacting with the text data gRPC service.
// It provides methods to save and load text data.
type TextDataPBClient struct {
	textDataService pb.TextDataServiceClient
}

// NewTextDataPBClient initializes a new TextDataPBClient with the provided gRPC service client.
// It returns a pointer to the TextDataPBClient instance.
func NewTextDataPBClient(u pb.TextDataServiceClient) *TextDataPBClient {
	return &TextDataPBClient{
		textDataService: u,
	}
}

// SaveTextData saves a new piece of text data using the text data service.
// It takes a context, a token for authorization, and the text data to be saved.
// It returns the saved text data and any error encountered during the process.
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

// LoadAllTextDataInfo retrieves information about all text data stored in the service.
// It takes a context and a token for authorization.
// It returns a slice of DataInfo models and any error encountered.
func (u *TextDataPBClient) LoadAllTextDataInfo(ctx context.Context, token string) ([]model.DataInfo, error) {
	req := &pb.GetAllTextInfoRequest{}

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := u.textDataService.GetLoadAllTextDataInfo(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("load text data: %w", err)
	}

	textInfos := make([]model.DataInfo, 0, len(resp.Text))
	for _, data := range resp.Text {
		textInfos = append(textInfos, model.DataInfo{
			ID:        data.Id,
			DataType:  data.DataType,
			MetaData:  data.Metadata,
			CreatedAt: data.CreatedAt,
		})
	}

	return textInfos, nil
}

// LoadTextData retrieves a specific text data entry by its ID.
// It takes a context, a token for authorization, and the ID of the data to load.
// It returns the corresponding TextData model and any error encountered.
func (u *TextDataPBClient) LoadTextData(ctx context.Context, token string, dataID string) (model.TextData, error) {
	req := &pb.GetTextDataRequest{
		Id: dataID,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := u.textDataService.GetLoadTextData(ctx, req)
	if err != nil {
		return model.TextData{}, fmt.Errorf("load text data: %w", err)
	}
	data := resp.TextData
	binaryData := model.TextData{
		Text:     data.Text,
		MetaData: data.Metadata,
	}

	return binaryData, nil
}

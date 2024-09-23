package pbclient

import (
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/text_data"
)

type TextDataPBClient struct {
	textDataService pb.TextDataServiceClient
}

func NewCreditCardPBClient(u pb.TextDataServiceClient) *TextDataPBClient {
	return &TextDataPBClient{
		textDataService: u,
	}
}

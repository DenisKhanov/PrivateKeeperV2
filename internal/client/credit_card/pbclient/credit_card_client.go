package pbclient

import (
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credit_card"
)

type CreditCardPBClient struct {
	creditCardService pb.CreditCardServiceClient
}

func NewCreditCardPBClient(u pb.CreditCardServiceClient) *CreditCardPBClient {
	return &CreditCardPBClient{
		creditCardService: u,
	}
}

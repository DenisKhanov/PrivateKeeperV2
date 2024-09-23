package pbclient

import (
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credentials"
)

type CredentialsPBClient struct {
	credentialsService pb.CredentialsServiceClient
}

func NewCredentialsPBClient(u pb.CredentialsServiceClient) *CredentialsPBClient {
	return &CredentialsPBClient{
		credentialsService: u,
	}
}

package pbclient

import (
	"context"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credentials"
	"google.golang.org/grpc/metadata"
)

type CredentialsPBClient struct {
	credentialsService pb.CredentialsServiceClient
}

func NewCredentialsPBClient(u pb.CredentialsServiceClient) *CredentialsPBClient {
	return &CredentialsPBClient{
		credentialsService: u,
	}
}

func (u *CredentialsPBClient) SaveCredentials(ctx context.Context, token string, cred model.CredentialsPostRequest) (model.Credentials, error) {
	req := &pb.PostCredentialsRequest{
		Login:    cred.Login,
		Password: cred.Password,
		Metadata: cred.MetaData,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := u.credentialsService.PostSaveCredentials(ctx, req)
	if err != nil {
		return model.Credentials{}, err
	}

	credential := model.Credentials{
		Login:    resp.Login,
		Password: resp.Password,
		MetaData: resp.Metadata,
	}

	return credential, nil
}

func (u *CredentialsPBClient) LoadCredentials(ctx context.Context, token string, cred model.CredentialsLoadRequest) ([]model.Credentials, error) {
	req := &pb.GetCredentialsRequest{
		Login:    cred.Login,
		Password: cred.Password,
		Metadata: cred.MetaData,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := u.credentialsService.GetLoadCredentials(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("load credentials: %w", err)
	}

	creds := make([]model.Credentials, 0, len(resp.Creds))
	for _, cr := range resp.Creds {
		creds = append(creds, model.Credentials{
			Login:    cr.Login,
			Password: cr.Password,
			MetaData: cr.Metadata,
		})
	}

	return creds, nil
}

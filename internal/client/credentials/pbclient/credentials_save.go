package pbclient

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credentials"
)

func (u *CredentialsPBClient) SaveCredentials(token string, cred model.CredentialsPostRequest) (model.Credentials, error) {
	req := &pb.PostCredentialsRequest{
		Login:    cred.Login,
		Password: cred.Password,
		Metadata: cred.MetaData,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

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

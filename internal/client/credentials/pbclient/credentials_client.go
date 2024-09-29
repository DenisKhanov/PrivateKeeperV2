package pbclient

import (
	"context"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/client/model"
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credentials"
	"google.golang.org/grpc/metadata"
)

// CredentialsPBClient is a client wrapper around the gRPC CredentialsServiceClient,
// providing methods to interact with the credentials-related operations via gRPC.
type CredentialsPBClient struct {
	credentialsService pb.CredentialsServiceClient
}

// NewCredentialsPBClient initializes and returns a new instance of CredentialsPBClient
// which will use the provided gRPC CredentialsServiceClient.
func NewCredentialsPBClient(u pb.CredentialsServiceClient) *CredentialsPBClient {
	return &CredentialsPBClient{
		credentialsService: u,
	}
}

// SaveCredentials sends a request to save credentials to the gRPC credentials service.
// It accepts a context, an authentication token, and a model containing the login, password, and metadata.
// It returns the saved credentials or an error if the operation fails.
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

// LoadAllCredentialsDataInfo fetches all stored credentials metadata from the gRPC service.
// It accepts a context and an authentication token, and returns a list of credential info or an error.
func (u *CredentialsPBClient) LoadAllCredentialsDataInfo(ctx context.Context, token string) ([]model.DataInfo, error) {
	req := &pb.GetAllCredentialsInfoRequest{}

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := u.credentialsService.GetLoadAllCredentialsDataInfo(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("load credentials data: %w", err)
	}

	credentials := make([]model.DataInfo, 0, len(resp.Creds))
	for _, data := range resp.Creds {
		credentials = append(credentials, model.DataInfo{
			ID:        data.Id,
			DataType:  data.DataType,
			MetaData:  data.Metadata,
			CreatedAt: data.CreatedAt,
		})
	}

	return credentials, nil
}

// LoadCredentialsData retrieves specific credentials data by its ID from the gRPC service.
// It accepts a context, an authentication token, and the ID of the data to be loaded.
// It returns the requested credentials or an error if the operation fails.
func (u *CredentialsPBClient) LoadCredentialsData(ctx context.Context, token string, dataID string) (model.Credentials, error) {
	req := &pb.GetCredentialsRequest{
		Id: dataID,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := u.credentialsService.GetLoadCredentials(ctx, req)
	if err != nil {
		return model.Credentials{}, fmt.Errorf("load credentials data: %w", err)
	}
	data := resp.CredentialsData
	credentialsData := model.Credentials{
		Login:    data.Login,
		Password: data.Password,
		MetaData: data.Metadata,
	}

	return credentialsData, nil
}

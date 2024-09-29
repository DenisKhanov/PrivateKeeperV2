package pbclient

import (
	"context"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/user"
)

// UserPBClient is a client for interacting with the user service via gRPC.
// It holds a reference to the user service client interface.
type UserPBClient struct {
	userService pb.UserServiceClient
}

// NewUserPBClient initializes a new UserPBClient with the provided user service client.
// It returns a pointer to the newly created UserPBClient.
func NewUserPBClient(u pb.UserServiceClient) *UserPBClient {
	return &UserPBClient{
		userService: u,
	}
}

// LoginUser attempts to log in a user with the provided login credentials.
// It sends a login request to the user service and returns the user's token if successful.
func (u *UserPBClient) LoginUser(ctx context.Context, login, password string) (string, error) {
	req := &pb.PostUserLoginRequest{
		Login:    login,
		Password: password,
	}

	resp, err := u.userService.PostLoginUser(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Token, nil
}

// RegisterUser registers a new user with the provided login credentials.
// It sends a registration request to the user service and returns the user's token if successful.
func (u *UserPBClient) RegisterUser(ctx context.Context, login, password string) (string, error) {
	req := &pb.PostUserRegisterRequest{
		Login:    login,
		Password: password,
	}

	resp, err := u.userService.PostRegisterUser(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Token, nil
}

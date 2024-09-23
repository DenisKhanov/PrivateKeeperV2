package pbclient

import (
	"context"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/user"
)

type UserPBClient struct {
	userService pb.UserServiceClient
}

func NewUserPBClient(u pb.UserServiceClient) *UserPBClient {
	return &UserPBClient{
		userService: u,
	}
}

func (u *UserPBClient) LoginUser(login, password string) (string, error) {
	req := &pb.PostUserLoginRequest{
		Login:    login,
		Password: password,
	}

	resp, err := u.userService.PostLoginUser(context.Background(), req)
	if err != nil {
		return "", err
	}

	return resp.Token, nil
}

func (u *UserPBClient) RegisterUser(login, password string) (string, error) {
	req := &pb.PostUserRegisterRequest{
		Login:    login,
		Password: password,
	}

	resp, err := u.userService.PostRegisterUser(context.Background(), req)
	if err != nil {
		return "", err
	}

	return resp.Token, nil
}

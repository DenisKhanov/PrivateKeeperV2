package pbclient

import (
	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/binary_data"
)

type BinaryDataPBClient struct {
	binaryDataService pb.BinaryDataServiceClient
}

func NewBinaryDataPBClient(u pb.BinaryDataServiceClient) *BinaryDataPBClient {
	return &BinaryDataPBClient{
		binaryDataService: u,
	}
}

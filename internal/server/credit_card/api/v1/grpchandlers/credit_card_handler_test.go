package grpchandlers

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credit_card"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/config"
	creditCardValidation "github.com/DenisKhanov/PrivateKeeperV2/internal/server/credit_card/api/v1/validation"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/auth"
	mocks "github.com/DenisKhanov/PrivateKeeperV2/internal/server/mocks/credit_card"
	"github.com/DenisKhanov/PrivateKeeperV2/pkg/jwtmanager"
)

var cfgMock = &config.Config{
	GRPCServer:    ":3300",
	TokenName:     "token",
	TokenSecret:   "secret",
	TokenExpHours: 24,
}

type CreditCardHandlerTestSuite struct {
	suite.Suite
	creditCardService *mocks.MockCreditCardService
	dialer            func(ctx context.Context, address string) (net.Conn, error)
	jwtManager        *jwtmanager.JWTManager
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(CreditCardHandlerTestSuite))
}

func (c *CreditCardHandlerTestSuite) SetupSuite() {
	ctrl := gomock.NewController(c.T())
	c.creditCardService = mocks.NewMockCreditCardService(ctrl)
	c.jwtManager = jwtmanager.New(cfgMock.TokenName, cfgMock.TokenSecret, cfgMock.TokenExpHours)
	authentication := auth.New(c.jwtManager).GRPCJWTAuth
	validate := validator.New()
	creditCardValidator, err := creditCardValidation.New(validate)
	require.NoError(c.T(), err)

	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(authentication))
	pb.RegisterCreditCardServiceServer(server, New(c.creditCardService, creditCardValidator))

	c.dialer = func(ctx context.Context, address string) (net.Conn, error) {
		return lis.Dial()
	}

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}

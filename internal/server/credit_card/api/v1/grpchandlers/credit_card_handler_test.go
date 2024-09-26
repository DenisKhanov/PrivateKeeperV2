package grpchandlers

import (
	"context"
	"errors"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"testing"
	"time"

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

func (c *CreditCardHandlerTestSuite) Test_PostSaveCreditCard() {
	token, err := c.jwtManager.BuildJWTString("050a289a-d10a-417b-ab89-3acfca0f6529")
	require.NoError(c.T(), err)

	testCases := []struct {
		name                         string
		token                        string
		body                         *pb.PostCreditCardRequest
		expectedCode                 codes.Code
		expectedStatusMessage        string
		expectedViolationField       string
		expectedViolationDescription string
		prepare                      func()
	}{
		{
			name:  "BadRequest - invalid number",
			token: token,
			body: &pb.PostCreditCardRequest{
				Number:    "4368 0811 1360",
				OwnerName: "user name",
				ExpiresAt: "20-06-2024",
				CvvCode:   "111",
				PinCode:   "2222",
				Metadata:  "some user metadata",
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid credit_card post request",
			expectedViolationField:       "Number",
			expectedViolationDescription: "must be valid card_number",
			prepare: func() {
				c.creditCardService.EXPECT().SaveCreditCard(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "BadRequest - invalid owner",
			token: token,
			body: &pb.PostCreditCardRequest{
				Number:    "4368 0811 1360 1890",
				OwnerName: "user_name",
				ExpiresAt: "20-06-2024",
				CvvCode:   "111",
				PinCode:   "2222",
				Metadata:  "some user metadata",
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid credit_card post request",
			expectedViolationField:       "OwnerName",
			expectedViolationDescription: "must be valid owner: first_name second_name",
			prepare: func() {
				c.creditCardService.EXPECT().SaveCreditCard(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "BadRequest - invalid expire date",
			token: token,
			body: &pb.PostCreditCardRequest{
				Number:    "4368 0811 1360 1890",
				OwnerName: "user name",
				ExpiresAt: "06-20-2024",
				CvvCode:   "111",
				PinCode:   "2222",
				Metadata:  "some user metadata",
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid credit_card post request",
			expectedViolationField:       "ExpiresAt",
			expectedViolationDescription: "expires_at must be in DD-MM-YYYY format",
			prepare: func() {
				c.creditCardService.EXPECT().SaveCreditCard(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "BadRequest - invalid cvv",
			token: token,
			body: &pb.PostCreditCardRequest{
				Number:    "4368 0811 1360 1890",
				OwnerName: "user name",
				ExpiresAt: "20-06-2024",
				CvvCode:   "wrong cvv",
				PinCode:   "2222",
				Metadata:  "some user metadata",
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid credit_card post request",
			expectedViolationField:       "CVV",
			expectedViolationDescription: "must be valid cvv",
			prepare: func() {
				c.creditCardService.EXPECT().SaveCreditCard(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "BadRequest - invalid PIN",
			token: token,
			body: &pb.PostCreditCardRequest{
				Number:    "4368 0811 1360 1890",
				OwnerName: "user name",
				ExpiresAt: "20-06-2024",
				CvvCode:   "111",
				PinCode:   "22",
				Metadata:  "some user metadata",
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid credit_card post request",
			expectedViolationField:       "PinCode",
			expectedViolationDescription: "must be valid pin",
			prepare: func() {
				c.creditCardService.EXPECT().SaveCreditCard(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "BadRequest - name is required",
			token: token,
			body: &pb.PostCreditCardRequest{
				Number:    "",
				OwnerName: "user name",
				ExpiresAt: "20-06-2024",
				CvvCode:   "111",
				PinCode:   "2222",
				Metadata:  "some user metadata",
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid credit_card post request",
			expectedViolationField:       "Number",
			expectedViolationDescription: "is required",
			prepare: func() {
				c.creditCardService.EXPECT().SaveCreditCard(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Unauthorized - token not found",
			token: "",
			body: &pb.PostCreditCardRequest{
				Number:    "4368 0811 1360 1890",
				OwnerName: "user name",
				ExpiresAt: "20-06-2024",
				CvvCode:   "wrong cvv",
				PinCode:   "2222",
				Metadata:  "some user metadata",
			},
			expectedCode:          codes.Unauthenticated,
			expectedStatusMessage: "authentification by UserID failed",
			prepare: func() {
				c.creditCardService.EXPECT().SaveCreditCard(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Internal error - unable to save credit card",
			token: token,
			body: &pb.PostCreditCardRequest{
				Number:    "4368 0811 1360 1890",
				OwnerName: "user name",
				ExpiresAt: "20-06-2024",
				CvvCode:   "111",
				PinCode:   "2222",
				Metadata:  "some user metadata",
			},
			expectedCode:          codes.Internal,
			expectedStatusMessage: "internal error",
			prepare: func() {
				c.creditCardService.EXPECT().SaveCreditCard(gomock.Any(), gomock.Any()).Times(1).Return(model.CreditCard{}, errors.New("some error"))
			},
		},
	}
	for _, test := range testCases {
		c.T().Run(test.name, func(t *testing.T) {
			test.prepare()

			header := metadata.New(map[string]string{"token": test.token})
			ctx := metadata.NewOutgoingContext(context.Background(), header)
			conn, err := grpc.NewClient("passthrough:///bufnet",
				grpc.WithContextDialer(c.dialer),
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			require.NoError(t, err)
			defer conn.Close()

			client := pb.NewCreditCardServiceClient(conn)
			_, err = client.PostSaveCreditCard(ctx, test.body)
			st := status.Convert(err)
			assert.Equal(t, test.expectedCode, st.Code())
			assert.Equal(t, test.expectedStatusMessage, st.Message())
			for _, detail := range st.Details() {
				switch d := detail.(type) { //nolint:gocritic
				case *errdetails.BadRequest:
					for _, violation := range d.GetFieldViolations() {
						assert.Equal(t, test.expectedViolationField, violation.GetField())
						assert.Equal(t, test.expectedViolationDescription, violation.GetDescription())
					}
				}
			}
		})
	}
}

func (c *CreditCardHandlerTestSuite) Test_GetLoadCreditCard() {
	token, err := c.jwtManager.BuildJWTString("050a289a-d10a-417b-ab89-3acfca0f6529")
	require.NoError(c.T(), err)

	cards := []model.CreditCard{
		{
			ID:        "some id",
			OwnerID:   "some owner id",
			Number:    "some number",
			OwnerName: "some name",
			ExpiresAt: "20-06-2024",
			CVV:       "111",
			PinCode:   "2222",
			MetaData:  "some metadata",
			CreatedAt: time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:        "another id",
			OwnerID:   "another owner id",
			Number:    "another number",
			OwnerName: "another name",
			ExpiresAt: "20-06-2025",
			CVV:       "111",
			PinCode:   "2222",
			MetaData:  "another metadata",
			CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	creditCards := make([]*pb.CreditCard, 0, len(cards))
	for _, v := range cards {
		creditCards = append(creditCards, &pb.CreditCard{
			Id:        v.ID,
			OwnerId:   v.OwnerID,
			Number:    v.Number,
			OwnerName: v.OwnerName,
			ExpiresAt: v.ExpiresAt,
			CvvCode:   v.CVV,
			PinCode:   v.PinCode,
			Metadata:  v.MetaData,
			CreatedAt: v.CreatedAt.Format(time.RFC3339),
		})
	}

	testCases := []struct {
		name                         string
		token                        string
		body                         *pb.GetCreditCardRequest
		expectedCode                 codes.Code
		expectedStatusMessage        string
		expectedViolationField       string
		expectedViolationDescription string
		prepare                      func()
		expectedBody                 *pb.GetCreditCardResponse
	}{
		{
			name:  "BadRequest - invalid ExpiresAfter",
			token: token,
			body: &pb.GetCreditCardRequest{
				Number:        "4368 0811 1360 1890",
				Owner:         "user_name",
				CvvCode:       "111",
				PinCode:       "2222",
				Metadata:      "some user metadata",
				ExpiresAfter:  "06-24-2024",
				ExpiresBefore: "20-06-2024",
			},
			expectedCode:          codes.InvalidArgument,
			expectedStatusMessage: "expires after must be in format 'DD-MM-YYYY'",
			prepare: func() {
				c.creditCardService.EXPECT().LoadAllCreditCard(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "BadRequest - invalid ExpiresBefore",
			token: token,
			body: &pb.GetCreditCardRequest{
				Number:        "4368 0811 1360 1890",
				Owner:         "user_name",
				CvvCode:       "111",
				PinCode:       "2222",
				Metadata:      "some user metadata",
				ExpiresAfter:  "20-06-2024",
				ExpiresBefore: "06-24-2024",
			},
			expectedCode:          codes.InvalidArgument,
			expectedStatusMessage: "expires before must be in format 'DD-MM-YYYY'",
			prepare: func() {
				c.creditCardService.EXPECT().LoadAllCreditCard(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Unauthorized - token not found",
			token: "",
			body: &pb.GetCreditCardRequest{
				Number:        "4368 0811 1360 1890",
				Owner:         "user_name",
				CvvCode:       "111",
				PinCode:       "2222",
				Metadata:      "some user metadata",
				ExpiresAfter:  "20-06-2024",
				ExpiresBefore: "06-24-2024",
			},
			expectedCode:          codes.Unauthenticated,
			expectedStatusMessage: "authentification by UserID failed",
			prepare: func() {
				c.creditCardService.EXPECT().LoadAllCreditCard(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Internal error - unable to load credit card",
			token: token,
			body: &pb.GetCreditCardRequest{
				Number:        "4368 0811 1360 1890",
				Owner:         "user_name",
				CvvCode:       "111",
				PinCode:       "2222",
				Metadata:      "some user metadata",
				ExpiresAfter:  "20-06-2024",
				ExpiresBefore: "20-06-2024",
			},
			expectedCode:          codes.Internal,
			expectedStatusMessage: "internal error",
			prepare: func() {
				c.creditCardService.EXPECT().LoadAllCreditCard(gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.New("error"))
			},
		},
		{
			name:  "Success - load credit cards",
			token: token,
			body: &pb.GetCreditCardRequest{
				Number:        "4368 0811 1360 1890",
				Owner:         "user_name",
				CvvCode:       "111",
				PinCode:       "2222",
				Metadata:      "some user metadata",
				ExpiresAfter:  "20-06-2024",
				ExpiresBefore: "20-06-2024",
			},
			expectedCode:          codes.OK,
			expectedStatusMessage: "",
			prepare: func() {
				c.creditCardService.EXPECT().LoadAllCreditCard(gomock.Any(), gomock.Any()).Times(1).Return(cards, nil)
			},
			expectedBody: &pb.GetCreditCardResponse{Cards: creditCards},
		},
	}
	for _, test := range testCases {
		c.T().Run(test.name, func(t *testing.T) {
			test.prepare()

			header := metadata.New(map[string]string{"token": test.token})
			ctx := metadata.NewOutgoingContext(context.Background(), header)
			conn, err := grpc.NewClient("passthrough:///bufnet",
				grpc.WithContextDialer(c.dialer),
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			require.NoError(t, err)
			defer conn.Close()

			client := pb.NewCreditCardServiceClient(conn)
			resp, err := client.GetLoadCreditCard(ctx, test.body)
			st := status.Convert(err)
			assert.Equal(t, test.expectedCode, st.Code())
			assert.Equal(t, test.expectedStatusMessage, st.Message())
			for _, detail := range st.Details() {
				switch d := detail.(type) { //nolint:gocritic
				case *errdetails.BadRequest:
					for _, violation := range d.GetFieldViolations() {
						assert.Equal(t, test.expectedViolationField, violation.GetField())
						assert.Equal(t, test.expectedViolationDescription, violation.GetDescription())
					}
				}
			}
			if resp != nil {
				assert.Equal(t, test.expectedBody.GetCards(), resp.GetCards())
			}
		})
	}
}

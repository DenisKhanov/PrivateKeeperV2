package grpchandlers

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credit_card"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

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

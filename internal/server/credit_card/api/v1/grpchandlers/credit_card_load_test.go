package grpchandlers

import (
	"context"
	"errors"
	"testing"
	"time"

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
			UpdatedAt: time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
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
			UpdatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
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
			UpdatedAt: v.UpdatedAt.Format(time.RFC3339),
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

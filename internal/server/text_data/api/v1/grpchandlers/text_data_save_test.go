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

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/text_data"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

func (c *TextDataHandlerTestSuite) Test_PostSaveTextData() {
	token, err := c.jwtManager.BuildJWTString("050a289a-d10a-417b-ab89-3acfca0f6529")
	require.NoError(c.T(), err)

	testCases := []struct {
		name                         string
		token                        string
		body                         *pb.PostTextDataRequest
		expectedCode                 codes.Code
		expectedStatusMessage        string
		expectedViolationField       string
		expectedViolationDescription string
		prepare                      func()
	}{
		{
			name:  "BadRequest - text is required",
			token: token,
			body: &pb.PostTextDataRequest{
				Text:     "",
				Metadata: "some metadata",
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid text_data post request",
			expectedViolationField:       "Text",
			expectedViolationDescription: "is required",
			prepare: func() {
				c.textDataService.EXPECT().SaveTextData(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Unauthorized - token not found",
			token: "",
			body: &pb.PostTextDataRequest{
				Text:     "some valuable text",
				Metadata: "some metadata",
			},
			expectedCode:          codes.Unauthenticated,
			expectedStatusMessage: "authentification by UserID failed",
			prepare: func() {
				c.textDataService.EXPECT().SaveTextData(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Internal error - unable to save text data",
			token: token,
			body: &pb.PostTextDataRequest{
				Text:     "some valuable text",
				Metadata: "some metadata",
			},
			expectedCode:          codes.Internal,
			expectedStatusMessage: "internal error",
			prepare: func() {
				c.textDataService.EXPECT().SaveTextData(gomock.Any(), gomock.Any()).Times(1).Return(model.TextData{}, errors.New("some error"))
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

			client := pb.NewTextDataServiceClient(conn)
			_, err = client.PostSaveTextData(ctx, test.body)
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

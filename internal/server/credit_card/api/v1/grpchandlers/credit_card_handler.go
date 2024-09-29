package grpchandlers

import (
	"context"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/lib"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credit_card"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

// CreditCardService defines the methods for operations related to credit cards.
type CreditCardService interface {
	SaveCreditCard(ctx context.Context, req model.CreditCardPostRequest) (model.CreditCard, error)
	LoadCreditCardData(ctx context.Context, dataID string) (model.CreditCard, error)
	LoadAllCreditCardInfo(ctx context.Context) ([]model.DataInfo, error)
}

// Validator defines the method for validating credit card requests.
type Validator interface {
	ValidatePostRequest(req *model.CreditCardPostRequest) (map[string]string, bool)
}

// CreditCardHandler is the gRPC handler for credit card-related operations.
type CreditCardHandler struct {
	creditCardService                       CreditCardService // The service for credit card operations
	pb.UnimplementedCreditCardServiceServer                   // Embed the unimplemented server for compatibility
	validator                               Validator         // The validator for incoming requests
}

// New creates a new instance of CreditCardHandler with the provided dependencies.
func New(creditCardService CreditCardService, validator Validator) *CreditCardHandler {
	return &CreditCardHandler{
		creditCardService: creditCardService,
		validator:         validator,
	}
}

// PostSaveCreditCard handles the gRPC call for saving a new credit card.
// It validates the request and saves the credit card data.
func (h *CreditCardHandler) PostSaveCreditCard(ctx context.Context, in *pb.PostCreditCardRequest) (*pb.PostCreditCardResponse, error) {
	req := model.CreditCardPostRequest{
		Number:    in.Number,
		OwnerName: in.OwnerName,
		ExpiresAt: in.ExpiresAt,
		CVV:       in.CvvCode,
		PinCode:   in.PinCode,
		MetaData:  in.Metadata,
	}

	report, ok := h.validator.ValidatePostRequest(&req)
	if !ok {
		logrus.Info("Unable to register user: invalid user request")
		logrus.Infof("violated_fields %v", report)
		return nil, lib.ProcessValidationError("invalid credit_card post request", report)
	}

	creditCard, err := h.creditCardService.SaveCreditCard(ctx, req)
	if err != nil {
		logrus.WithError(err).Error("Unable to save credit_card")
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.PostCreditCardResponse{
		Id:        creditCard.ID,
		OwnerId:   creditCard.OwnerID,
		Number:    creditCard.Number,
		OwnerName: creditCard.OwnerName,
		ExpiresAt: creditCard.ExpiresAt,
		CvvCode:   creditCard.CVV,
		PinCode:   creditCard.PinCode,
		Metadata:  creditCard.MetaData,
		CreatedAt: creditCard.CreatedAt.Format(time.RFC3339),
	}, nil
}

// GetLoadAllCreditCardDataInfo handles the gRPC call for loading all credit card information.
func (h *CreditCardHandler) GetLoadAllCreditCardDataInfo(ctx context.Context, _ *pb.GetAllCreditCardInfoRequest) (*pb.GetAllCreditCardInfoResponse, error) {

	cardInfo, err := h.creditCardService.LoadAllCreditCardInfo(ctx)
	if err != nil {
		logrus.WithError(err).Error("Error while loading credit card data: ")
		return nil, status.Error(codes.Internal, "internal error")
	}

	cardInfos := make([]*pb.CreditCardInfo, 0, len(cardInfo))
	for _, v := range cardInfo {
		cardInfos = append(cardInfos, &pb.CreditCardInfo{
			Id:        v.ID,
			DataType:  v.DataType,
			Metadata:  v.MetaData,
			CreatedAt: v.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.GetAllCreditCardInfoResponse{Cards: cardInfos}, nil
}

// GetLoadCreditCard handles the gRPC call for loading a specific credit card's data.
func (h *CreditCardHandler) GetLoadCreditCard(ctx context.Context, in *pb.GetCreditCardRequest) (*pb.GetCreditCardResponse, error) {
	dataID := in.Id

	cardData, err := h.creditCardService.LoadCreditCardData(ctx, dataID)
	if err != nil {
		logrus.WithError(err).Error("Error while loading credit card data: ")
		return nil, status.Error(codes.Internal, "internal error")
	}

	card := &pb.CreditCard{
		Id:        cardData.ID,
		OwnerId:   cardData.OwnerID,
		Number:    cardData.Number,
		OwnerName: cardData.OwnerName,
		ExpiresAt: cardData.ExpiresAt,
		CvvCode:   cardData.CVV,
		PinCode:   cardData.PinCode,
		Metadata:  cardData.MetaData,
		CreatedAt: cardData.CreatedAt.Format(time.RFC3339Nano),
	}
	return &pb.GetCreditCardResponse{CardData: card}, nil
}

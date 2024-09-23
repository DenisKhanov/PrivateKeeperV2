package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/auth"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/keyextraction"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

func (s *CreditCardService) SaveCreditCard(ctx context.Context, req model.CreditCardPostRequest) (model.CreditCard, error) {
	userID, ok := ctx.Value(auth.UserIDContextKey("userID")).(string)
	if !ok {
		return model.CreditCard{}, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(keyextraction.UserKeyContextKey("userKey")).([]byte)
	if !ok {
		return model.CreditCard{}, fmt.Errorf("failed to get userKey from context")
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return model.CreditCard{}, fmt.Errorf("new uuid: %w", err)
	}

	card := model.CreditCardCryptData{
		Number:    req.Number,
		OwnerName: req.OwnerName,
		ExpiresAt: req.ExpiresAt,
		CVV:       req.CVV,
		PinCode:   req.PinCode,
	}

	data, err := json.Marshal(card)
	if err != nil {
		return model.CreditCard{}, fmt.Errorf("marshal: %w", err)
	}

	cryptData, err := s.crypt.Encrypt(userKey, data)
	if err != nil {
		return model.CreditCard{}, fmt.Errorf("encrypt data: %w", err)
	}

	dataToSave := model.Data{
		ID:       id.String(),
		OwnerID:  userID,
		Type:     s.dataType,
		Data:     cryptData,
		MetaData: req.MetaData,
	}

	savedCreditCard, err := s.repository.Insert(ctx, dataToSave)
	if err != nil {
		return model.CreditCard{}, fmt.Errorf("insert credit card: %w", err)
	}

	return model.CreditCard{
		ID:        savedCreditCard.ID,
		OwnerID:   savedCreditCard.OwnerID,
		Number:    req.Number,
		OwnerName: req.OwnerName,
		ExpiresAt: req.ExpiresAt,
		CVV:       req.CVV,
		PinCode:   req.PinCode,
		MetaData:  savedCreditCard.MetaData,
		CreatedAt: savedCreditCard.CreatedAt,
		UpdatedAt: savedCreditCard.UpdatedAt,
	}, nil
}

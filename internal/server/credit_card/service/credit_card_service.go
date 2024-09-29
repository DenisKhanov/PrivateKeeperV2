package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/DenisKhanov/PrivateKeeperV2/pkg/jwtmanager"
)

const creditCard = "credit_card" // Constant for the credit card data type

// DataRepository interface defines methods for data access.
type DataRepository interface {
	Insert(ctx context.Context, data model.Data) (model.Data, error)
	SelectAll(ctx context.Context, userID, dataType string) ([]model.Data, error)
	SelectByID(ctx context.Context, userID, dataType, dataID string) (model.Data, error)
}

// CryptService interface defines methods for encryption and decryption.
type CryptService interface {
	Encrypt(key, data []byte) ([]byte, error)
	Decrypt(key, data []byte) ([]byte, error)
	GenerateKey() ([]byte, error)
}

// CreditCardService struct manages credit card operations.
type CreditCardService struct {
	repository DataRepository         // Repository for data operations
	crypt      CryptService           // Service for cryptographic operations
	jwtManager *jwtmanager.JWTManager // JWT manager for authentication
	dataType   string                 // Type of data managed (credit card)
}

// New creates a new instance of CreditCardService.
func New(repository DataRepository, crypt CryptService, jwtManager *jwtmanager.JWTManager) *CreditCardService {
	return &CreditCardService{
		repository: repository,
		crypt:      crypt,
		jwtManager: jwtManager,
		dataType:   creditCard,
	}
}

// SaveCreditCard saves a new credit card to the repository.
func (s *CreditCardService) SaveCreditCard(ctx context.Context, req model.CreditCardPostRequest) (model.CreditCard, error) {
	userID, ok := ctx.Value(model.UserIDKey).(string)
	if !ok {
		return model.CreditCard{}, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(model.UserKey).([]byte)
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
	}, nil
}

// LoadAllCreditCardInfo retrieves all credit card information for the user.
func (s *CreditCardService) LoadAllCreditCardInfo(ctx context.Context) ([]model.DataInfo, error) {
	userID, ok := ctx.Value(model.UserIDKey).(string)
	if !ok {
		return nil, fmt.Errorf("failed to get userID from context")
	}

	encryptedCardData, err := s.repository.SelectAll(ctx, userID, s.dataType)
	if err != nil {
		return nil, fmt.Errorf("select all binary_data: %w", err)
	}

	credCardDataInfo := make([]model.DataInfo, 0, len(encryptedCardData))
	for _, encryptedBinary := range encryptedCardData {
		credCardDataInfo = append(credCardDataInfo, model.DataInfo{
			ID:        encryptedBinary.ID,
			DataType:  s.dataType,
			MetaData:  encryptedBinary.MetaData,
			CreatedAt: encryptedBinary.CreatedAt,
		})
	}
	return credCardDataInfo, nil
}

// LoadCreditCardData retrieves and decrypts a specific credit card's data.
func (s *CreditCardService) LoadCreditCardData(ctx context.Context, dataID string) (model.CreditCard, error) {
	userID, ok := ctx.Value(model.UserIDKey).(string)
	if !ok {
		return model.CreditCard{}, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(model.UserKey).([]byte)
	if !ok {
		return model.CreditCard{}, fmt.Errorf("failed to get userKey from context")
	}

	encryptedCardData, err := s.repository.SelectByID(ctx, userID, s.dataType, dataID)
	if err != nil {
		return model.CreditCard{}, fmt.Errorf("select all credit_card_data: %w", err)
	}
	decryptedData, err := s.crypt.Decrypt(userKey, encryptedCardData.Data)
	if err != nil {
		return model.CreditCard{}, fmt.Errorf("decrypt card: %w", err)
	}

	var decryptedCardData model.CreditCardCryptData
	err = json.Unmarshal(decryptedData, &decryptedCardData)
	if err != nil {
		return model.CreditCard{}, fmt.Errorf("unmarshal card: %w", err)
	}

	card := model.CreditCard{
		ID:        encryptedCardData.ID,
		OwnerID:   encryptedCardData.OwnerID,
		Number:    decryptedCardData.Number,
		OwnerName: decryptedCardData.OwnerName,
		ExpiresAt: decryptedCardData.ExpiresAt,
		CVV:       decryptedCardData.CVV,
		PinCode:   decryptedCardData.PinCode,
		MetaData:  encryptedCardData.MetaData,
		CreatedAt: encryptedCardData.CreatedAt,
	}

	return card, nil
}

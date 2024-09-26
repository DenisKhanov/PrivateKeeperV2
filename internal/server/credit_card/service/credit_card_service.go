package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/credit_card/specification"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/auth"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/keyextraction"
	"github.com/google/uuid"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
	"github.com/DenisKhanov/PrivateKeeperV2/pkg/jwtmanager"
)

const creditCard = "credit_card"

type DataRepository interface {
	Insert(ctx context.Context, data model.Data) (model.Data, error)
	SelectAll(ctx context.Context, userID, dataType string) ([]model.Data, error)
}

type CryptService interface {
	Encrypt(key, data []byte) ([]byte, error)
	Decrypt(key, data []byte) ([]byte, error)
	GenerateKey() ([]byte, error)
}

type CreditCardService struct {
	repository DataRepository
	crypt      CryptService
	jwtManager *jwtmanager.JWTManager
	dataType   string
}

func New(repository DataRepository, crypt CryptService, jwtManager *jwtmanager.JWTManager) *CreditCardService {
	return &CreditCardService{
		repository: repository,
		crypt:      crypt,
		jwtManager: jwtManager,
		dataType:   creditCard,
	}
}

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

func (s *CreditCardService) LoadAllCreditCard(ctx context.Context, spec specification.CreditCardSpecification) ([]model.CreditCard, error) {
	userID, ok := ctx.Value(auth.UserIDContextKey("userID")).(string)
	if !ok {
		return nil, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(keyextraction.UserKeyContextKey("userKey")).([]byte)
	if !ok {
		return nil, fmt.Errorf("failed to get userKey from context")
	}

	encryptedCards, err := s.repository.SelectAll(ctx, userID, s.dataType)
	if err != nil {
		return nil, fmt.Errorf("select all cards: %w", err)
	}

	cards := make([]model.CreditCard, 0, len(encryptedCards))
	for _, encryptedCard := range encryptedCards {
		decryptedData, err := s.crypt.Decrypt(userKey, encryptedCard.Data)
		if err != nil {
			return nil, fmt.Errorf("decrypt data: %w", err)
		}

		var card model.CreditCardCryptData
		err = json.Unmarshal(decryptedData, &card)
		if err != nil {
			return nil, fmt.Errorf("unmarshal data: %w", err)
		}

		cards = append(cards, model.CreditCard{
			ID:        encryptedCard.ID,
			OwnerID:   encryptedCard.OwnerID,
			Number:    card.Number,
			OwnerName: card.OwnerName,
			ExpiresAt: card.ExpiresAt,
			CVV:       card.CVV,
			PinCode:   card.PinCode,
			MetaData:  encryptedCard.MetaData,
			CreatedAt: encryptedCard.CreatedAt,
			UpdatedAt: encryptedCard.UpdatedAt,
		})
	}

	predicates := spec.MakeFilterPredicates()
	var filteredCards []model.CreditCard
	for _, card := range cards {
		take := true
		for _, filterCardWithSpec := range predicates {
			if !filterCardWithSpec(spec, card) {
				take = false
				break
			}
		}
		if take {
			filteredCards = append(filteredCards, card)
		}
	}

	return filteredCards, nil
}

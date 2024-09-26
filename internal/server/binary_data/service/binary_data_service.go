package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/binary_data/specification"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/auth"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/keyextraction"
	"github.com/google/uuid"

	"github.com/DenisKhanov/PrivateKeeperV2/pkg/jwtmanager"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

const binaryData = "binary_data"

type DataRepository interface {
	Insert(ctx context.Context, data model.Data) (model.Data, error)
	SelectAll(ctx context.Context, userID, dataType string) ([]model.Data, error)
}

type CryptService interface {
	Encrypt(key, data []byte) ([]byte, error)
	Decrypt(key, data []byte) ([]byte, error)
	GenerateKey() ([]byte, error)
}

type BinaryDataService struct {
	repository DataRepository
	crypt      CryptService
	jwtManager *jwtmanager.JWTManager
	dataType   string
}

func New(repository DataRepository, crypt CryptService, jwtManager *jwtmanager.JWTManager) *BinaryDataService {
	return &BinaryDataService{
		repository: repository,
		crypt:      crypt,
		jwtManager: jwtManager,
		dataType:   binaryData,
	}
}

func (s *BinaryDataService) SaveBinaryData(ctx context.Context, req model.BinaryDataPostRequest) (model.BinaryData, error) {
	userID, ok := ctx.Value(auth.UserIDContextKey("userID")).(string)
	if !ok {
		return model.BinaryData{}, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(keyextraction.UserKeyContextKey("userKey")).([]byte)
	if !ok {
		return model.BinaryData{}, fmt.Errorf("failed to get userKey from context")
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return model.BinaryData{}, fmt.Errorf("new uuid: %w", err)
	}

	binary := model.BinaryCryptData{
		Name:      req.Name,
		Extension: req.Extension,
		Data:      req.Data,
	}

	data, err := json.Marshal(binary)
	if err != nil {
		return model.BinaryData{}, fmt.Errorf("marshal: %w", err)
	}

	cryptData, err := s.crypt.Encrypt(userKey, data)
	if err != nil {
		return model.BinaryData{}, fmt.Errorf("encrypt data: %w", err)
	}

	dataToSave := model.Data{
		ID:       id.String(),
		OwnerID:  userID,
		Type:     s.dataType,
		Data:     cryptData,
		MetaData: req.MetaData,
	}

	savedBinaryData, err := s.repository.Insert(ctx, dataToSave)
	if err != nil {
		return model.BinaryData{}, fmt.Errorf("insert credit card: %w", err)
	}

	return model.BinaryData{
		ID:        savedBinaryData.ID,
		OwnerID:   savedBinaryData.OwnerID,
		Name:      req.Name,
		Extension: req.Extension,
		Data:      req.Data,
		MetaData:  savedBinaryData.MetaData,
		CreatedAt: savedBinaryData.CreatedAt,
	}, nil
}

func (s *BinaryDataService) LoadAllBinaryData(ctx context.Context, spec specification.BinaryDataSpecification) ([]model.BinaryData, error) {
	userID, ok := ctx.Value(auth.UserIDContextKey("userID")).(string)
	if !ok {
		return nil, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(keyextraction.UserKeyContextKey("userKey")).([]byte)
	if !ok {
		return nil, fmt.Errorf("failed to get userKey from context")
	}

	encryptedBinaryData, err := s.repository.SelectAll(ctx, userID, s.dataType)
	if err != nil {
		return nil, fmt.Errorf("select all binary_data: %w", err)
	}

	texts := make([]model.BinaryData, 0, len(encryptedBinaryData))
	for _, encryptedBinary := range encryptedBinaryData {
		decryptedData, err := s.crypt.Decrypt(userKey, encryptedBinary.Data)
		if err != nil {
			return nil, fmt.Errorf("decrypt data: %w", err)
		}

		var data model.BinaryCryptData
		err = json.Unmarshal(decryptedData, &data)
		if err != nil {
			return nil, fmt.Errorf("unmarshal data: %w", err)
		}

		texts = append(texts, model.BinaryData{
			ID:        encryptedBinary.ID,
			OwnerID:   encryptedBinary.OwnerID,
			Name:      data.Name,
			Extension: data.Extension,
			Data:      data.Data,
			MetaData:  encryptedBinary.MetaData,
			CreatedAt: encryptedBinary.CreatedAt,
		})
	}

	predicates := spec.MakeFilterPredicates()
	var filteredBinaryData []model.BinaryData
	for _, binary := range texts {
		take := true
		for _, filteredBinaryDataWithSpec := range predicates {
			if !filteredBinaryDataWithSpec(spec, binary) {
				take = false
				break
			}
		}
		if take {
			filteredBinaryData = append(filteredBinaryData, binary)
		}
	}

	return filteredBinaryData, nil
}

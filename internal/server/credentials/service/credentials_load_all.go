package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/credentials/specification"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/auth"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/interceptors/keyextraction"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

func (s *CredentialsService) LoadAllCredentials(ctx context.Context, spec specification.CredentialsSpecification) ([]model.Credentials, error) {
	userID, ok := ctx.Value(auth.UserIDContextKey("userID")).(string)
	if !ok {
		return nil, fmt.Errorf("failed to get userID from context")
	}

	userKey, ok := ctx.Value(keyextraction.UserKeyContextKey("userKey")).([]byte)
	if !ok {
		return nil, fmt.Errorf("failed to get userKey from context")
	}

	encryptedCredentials, err := s.repository.SelectAll(ctx, userID, s.dataType)
	if err != nil {
		return nil, fmt.Errorf("select all credentials: %w", err)
	}

	creds := make([]model.Credentials, 0, len(encryptedCredentials))
	for _, encryptedCred := range encryptedCredentials {
		decryptedData, err := s.crypt.Decrypt(userKey, encryptedCred.Data)
		if err != nil {
			return nil, fmt.Errorf("decrypt data: %w", err)
		}

		var cred model.CredentialsCryptData
		err = json.Unmarshal(decryptedData, &cred)
		if err != nil {
			return nil, fmt.Errorf("unmarshal data: %w", err)
		}

		creds = append(creds, model.Credentials{
			ID:        encryptedCred.ID,
			OwnerID:   encryptedCred.OwnerID,
			Login:     cred.Login,
			Password:  cred.Password,
			MetaData:  encryptedCred.MetaData,
			CreatedAt: encryptedCred.CreatedAt,
			UpdatedAt: encryptedCred.UpdatedAt,
		})
	}

	predicates := spec.MakeFilterPredicates()
	var filteredCreds []model.Credentials
	for _, cred := range creds {
		take := true
		for _, filterCardWithSpec := range predicates {
			if !filterCardWithSpec(spec, cred) {
				take = false
				break
			}
		}
		if take {
			filteredCreds = append(filteredCreds, cred)
		}
	}

	return filteredCreds, nil
}

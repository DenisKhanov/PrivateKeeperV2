package encryption

import (
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)

type Service struct {
	aead cipher.AEAD
}

func New(key []byte) (*Service, error) {
	hash := sha256.Sum256(key)
	aead, err := chacha20poly1305.NewX(hash[:])
	if err != nil {
		return nil, fmt.Errorf("chacha20poly1305.NewX: %w", err)
	}

	return &Service{aead: aead}, nil
}

func (e *Service) EncryptWithMasterKey(data []byte) ([]byte, error) {
	nonce := make([]byte, e.aead.NonceSize())
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("rand.Read: %w", err)
	}

	ciphertext := e.aead.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (e *Service) DecryptWithMasterKey(data []byte) ([]byte, error) {
	nonce, ciphertext := data[:chacha20poly1305.NonceSizeX], data[chacha20poly1305.NonceSizeX:]
	dec, err := e.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("aead.Open: %w", err)
	}

	return dec, nil
}

func (e *Service) GenerateKey() ([]byte, error) {
	key := make([]byte, chacha20poly1305.KeySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (e *Service) Encrypt(key, data []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("chacha20poly1305.NewX: %w", err)
	}

	nonce := make([]byte, chacha20poly1305.NonceSizeX)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("rand.Read: %w", err)
	}

	ciphertext := aead.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (e *Service) Decrypt(key, data []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("chacha20poly1305.NewX: %w", err)
	}

	nonce, ciphertext := data[:chacha20poly1305.NonceSizeX], data[chacha20poly1305.NonceSizeX:]
	dec, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("aead.Open: %w", err)
	}

	return dec, nil
}

package jwtmanager

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager struct holds configuration for managing JWT tokens
type JWTManager struct {
	TokenName string        // The name of the token
	secretKey string        // The secret key used for signing the token
	tokenExp  time.Duration // Duration before the token expires
}

// claims struct represents the claims stored in the JWT
type claims struct {
	jwt.RegisteredClaims        // Embedding RegisteredClaims from the jwt package
	UserID               string // Custom field to store the UserID
}

// New returns a new instance of JWTManager.
func New(tokenName string, secretKey string, hours int) *JWTManager {
	return &JWTManager{
		TokenName: tokenName,
		secretKey: secretKey,
		tokenExp:  time.Duration(hours * int(time.Hour)),
	}
}

// BuildJWTString creates a JWT token with the provided userID.
func (j *JWTManager) BuildJWTString(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("sign token %w", err)
	}

	return tokenString, nil
}

// GetUserID returns the userID from the provided JWT token.
func (j *JWTManager) GetUserID(tokenString string) (string, error) {
	jwtClaims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, jwtClaims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("getUserID %w", errors.New("unexpected signing method"))
			}
			return []byte(j.secretKey), nil
		})
	if err != nil {
		return "", fmt.Errorf("buildJWTString parse token %w", err)
	}

	if !token.Valid {
		logrus.Warnf("JWTManager token invalid")
		logrus.Infof("token %s", tokenString)
		return "", fmt.Errorf("buildJWTString signstring %w", errors.New("token is not valid"))
	}

	return jwtClaims.UserID, nil
}

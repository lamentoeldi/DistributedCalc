package authenticator

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

const (
	tokenType    = "type"
	tokenAccess  = "access"
	tokenRefresh = "refresh"
)

type Authenticator struct {
	accessPk   *ecdsa.PrivateKey
	refreshPk  *ecdsa.PrivateKey
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthenticator(accessPk, refreshPk *ecdsa.PrivateKey, accessTTL, refreshTTL time.Duration) *Authenticator {
	return &Authenticator{
		accessPk:   accessPk,
		refreshPk:  refreshPk,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (a *Authenticator) IssueTokens(userID string) (*jwt.Token, *jwt.Token) {
	access := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": "orchestrator",
		"sub": userID,
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(a.accessTTL)),
	})

	refreshID, _ := uuid.NewV7()
	refresh := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": "orchestrator",
		"sub": userID,
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(a.accessTTL)),
		"jti": refreshID.String(),
	})

	access.Header[tokenType] = tokenAccess
	refresh.Header[tokenType] = tokenRefresh

	return access, refresh
}

func (a *Authenticator) SignTokens(access, refresh *jwt.Token) (string, string, error) {
	accessString, err := access.SignedString(a.accessPk)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshString, err := refresh.SignedString(a.refreshPk)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return accessString, refreshString, nil
}

func (a *Authenticator) VerifyAndExtract(tokenString string) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, a.keyFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return token.Claims, nil
}

func (a *Authenticator) keyFunc(token *jwt.Token) (any, error) {
	errInvalidTokenType := fmt.Errorf("invalid token type")

	tokenType, ok := token.Header[tokenType].(string)
	if !ok {
		return nil, errInvalidTokenType
	}

	switch tokenType {
	case "access":
		return &a.accessPk.PublicKey, nil
	case "refresh":
		return &a.refreshPk.PublicKey, nil
	}

	return nil, errInvalidTokenType
}

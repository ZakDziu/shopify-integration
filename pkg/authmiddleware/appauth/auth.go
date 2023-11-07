package appauth

import (
	"crypto/ecdsa"
	"net/http"
	"strings"
	"time"

	"shopify-integration/pkg/authmiddleware"
	"shopify-integration/pkg/logger"
	"shopify-integration/pkg/model"

	"github.com/dgrijalva/jwt-go"
)

const StringsNumber = 2

type AuthMiddleware struct {
	loc   *time.Location
	atKey *ecdsa.PrivateKey
	rtKey *ecdsa.PrivateKey
}

func NewAuthMiddleware(atKey, rtKey *ecdsa.PrivateKey) *AuthMiddleware {
	loc, _ := time.LoadLocation("Europe/Moscow")
	var middleware = &AuthMiddleware{
		loc:   loc,
		atKey: atKey,
		rtKey: rtKey,
	}

	return middleware
}

func (m *AuthMiddleware) GetUserID(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &authmiddleware.AccessClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				logger.Errorf("GetUserId.unexpected signing method: %v", token.Header["alg"])

				return nil, model.ErrUnauthorized
			}

			return &m.atKey.PublicKey, nil
		})
	if err != nil {
		return "", model.ErrUnauthorized
	}

	claims, ok := token.Claims.(*authmiddleware.AccessClaims)
	if !ok {
		logger.Errorf("Refresh.invalid token claims: %v", token.Claims)

		return "", model.ErrUnauthorized
	}

	if !token.Valid {
		return "", model.ErrUnauthorized
	}

	return claims.BaseClaims.CustomerID, nil
}

func (m *AuthMiddleware) CreateTokens(customerID string) (*authmiddleware.Tokens, error) {
	accessClaims, refreshClaims := authmiddleware.GenerateClaims(customerID)

	at := jwt.NewWithClaims(jwt.SigningMethodES256, accessClaims)
	accessToken, err := at.SignedString(m.atKey)
	if err != nil {
		return nil, err
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaims)
	refreshToken, err := rt.SignedString(m.rtKey)
	if err != nil {
		return nil, err
	}

	return &authmiddleware.Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func (m *AuthMiddleware) Refresh(tokens authmiddleware.Tokens) (*authmiddleware.Tokens, error) {
	token, err := jwt.ParseWithClaims(tokens.Refresh, &authmiddleware.RefreshClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				logger.Errorf("Refresh.unexpected signing method: %v", token.Header["alg"])

				return nil, model.ErrUnauthorized
			}

			return &m.rtKey.PublicKey, nil
		})
	if err != nil {
		return nil, model.ErrUnauthorized
	}

	claims, ok := token.Claims.(*authmiddleware.RefreshClaims)
	if !ok {
		logger.Errorf("Refresh.invalid token claims: %v", token.Claims)

		return nil, model.ErrUnauthorized
	}

	if !token.Valid {
		return nil, model.ErrUnauthorized
	}

	return m.CreateTokens(claims.CustomerID)
}

func (m *AuthMiddleware) ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == StringsNumber {
		return strArr[1]
	}

	return ""
}

// Validate verifies token signature.
func (m *AuthMiddleware) Validate(raw string) (*authmiddleware.AccessClaims, error) {
	token, err := jwt.ParseWithClaims(raw, &authmiddleware.AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			logger.Errorf("Validate.unexpected signing method: %v", token.Header["alg"])

			return nil, model.ErrUnauthorized
		}

		return &m.atKey.PublicKey, nil
	})

	if err != nil {
		return nil, model.ErrUnauthorized
	}

	claims, ok := token.Claims.(*authmiddleware.AccessClaims)
	if !ok {
		logger.Errorf("Validate.invalid token claims: %v", token.Claims)

		return nil, model.ErrUnauthorized
	}

	if !token.Valid {
		return nil, model.ErrUnauthorized
	}

	return claims, nil
}

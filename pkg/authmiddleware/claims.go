package authmiddleware

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

const (
	AccessTokenTTL  = time.Hour * 8
	RefreshTokenTTL = time.Hour * 24 * 7
)

type BaseClaims struct {
	jwt.StandardClaims
	CustomerID string `json:"id"`
}

type AccessClaims struct {
	BaseClaims
	AccessUUID string `json:"access_uuid"`
}

type RefreshClaims struct {
	BaseClaims
	RefreshUUID string `json:"refresh_uuid"`
}

type Tokens struct {
	Access  string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}

func NewClaims(customerID string, ttl time.Duration) BaseClaims {
	return BaseClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			Id:        uuid.NewV4().String(),
			IssuedAt:  time.Now().Unix(),
		},
		CustomerID: customerID,
	}
}

func GenerateClaims(customerID string) (*AccessClaims, *RefreshClaims) {
	access := AccessClaims{
		BaseClaims: NewClaims(customerID, AccessTokenTTL),
	}

	refresh := RefreshClaims{
		BaseClaims: NewClaims(customerID, RefreshTokenTTL),
	}

	access.AccessUUID = refresh.Id
	refresh.RefreshUUID = access.Id

	return &access, &refresh
}

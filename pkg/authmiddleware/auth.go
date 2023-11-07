package authmiddleware

import (
	"net/http"
)

type AuthMiddleware interface {
	CreateTokens(customerID string) (*Tokens, error)
	Refresh(tokens Tokens) (*Tokens, error)
	ExtractToken(r *http.Request) string
	Validate(raw string) (*AccessClaims, error)
	GetUserID(accessToken string) (string, error)
}

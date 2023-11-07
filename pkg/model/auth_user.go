package model

import (
	"strings"
)

type AuthUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *AuthUser) IsValid() bool {
	email := strings.TrimSpace(a.Email)
	pass := strings.TrimSpace(a.Password)
	if email == "" || pass == "" {
		return false
	}

	a.Email = email
	a.Password = pass

	return true
}

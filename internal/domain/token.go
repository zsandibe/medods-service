package domain

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken []byte `json:"refresh_token"`
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	SessionId uuid.UUID `json:"session_id"`
	Guid      uuid.UUID `json:"guid"`
}

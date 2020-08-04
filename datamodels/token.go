package datamodels

import (
	"time"
)

const (
	ISSUER         = "demo"
	AUDIENCE       = "demo_client"
	EFFECTIVE_TIME = time.Hour * 2

	REFRESH_TOKEN_EFFECTIVE_TIME = time.Hour * 24 * 30
)

type RefreshToken struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

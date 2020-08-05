package datamodels

import (
	"time"
)

const (
	ISSUER         = "demo"
	AUDIENCE       = "demo_client"
	EFFECTIVE_TIME = int64(time.Hour * 2 / time.Second)

	REFRESH_TOKEN_EFFECTIVE_TIME = int64(time.Hour * 24 * 30 / time.Second)
)

type RefreshToken struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

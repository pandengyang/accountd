package datamodels

import (
	"time"
)

const (
	VERIFICATION_CODE_LEN           = 4
	VERIFICATION_CODE_SEND_INTERVAL = 150
	VerificationCodeExpire          = int64(15 * time.Minute / time.Second)
)

type VerificationCode struct {
	Phone            string `redis:"p,omitemtpy"`
	VerificationCode string `redis:"c,omitemtpy"`
	SentAt           int64  `redis:"s,omitemtpy"`
}

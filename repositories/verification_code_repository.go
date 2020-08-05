package repositories

import (
	"accountd/datamodels"
)

type VerificationCodeRepository interface {
	Insert(vc *datamodels.VerificationCode) (insertedPhone string, err error)
	SelectByPhone(phone string) (vc datamodels.VerificationCode, err error)
}

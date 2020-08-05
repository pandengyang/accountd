package services

import (
	"accountd/datamodels"
	"accountd/repositories"
)

type VerificationCodeService interface {
	Insert(vc *datamodels.VerificationCode) (insertedPhone string, err error)
	SelectByPhone(phone string) (vc datamodels.VerificationCode, err error)
}

type verificationCodeService struct {
	persistenceRepo repositories.VerificationCodeRepository
	cacheRepo       repositories.VerificationCodeRepository
}

func NewVerificationCodeService(persistenceRepo, cacheRepo repositories.VerificationCodeRepository) VerificationCodeService {
	return &verificationCodeService{
		persistenceRepo: persistenceRepo,
		cacheRepo:       cacheRepo,
	}
}

func (s *verificationCodeService) Insert(vc *datamodels.VerificationCode) (insertedPhone string, err error) {
	return s.cacheRepo.Insert(vc)
}

func (s *verificationCodeService) SelectByPhone(phone string) (vc datamodels.VerificationCode, err error) {
	return s.cacheRepo.SelectByPhone(phone)
}

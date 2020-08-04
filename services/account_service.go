package services

import (
	"accountd/datamodels"
	"accountd/repositories"
)

type AccountService interface {
	Insert(account *datamodels.Account) (insertedId int64, err error)
	Delete(insertedId int64) (rowsAffected int64, err error)

	UpdateNickname(insertedId int64, account *datamodels.Account) (rowsAffected int64, err error)
	UpdatePhone(insertedId int64, account *datamodels.Account) (rowsAffected int64, err error)
	UpdatePassword(insertedId int64, account *datamodels.Account) (rowsAffected int64, err error)

	SelectAllPerPage(page int64, pageSize int64) (accounts []interface{}, total int64, err error)
	SelectAll() (accounts []interface{}, total int64, err error)

	SelectAuthByNickname(nickname string) (account datamodels.Account, err error)
	SelectAuthByPhone(phone string) (account datamodels.Account, err error)
}

type accountService struct {
	persistenceRepo repositories.AccountRepository
	cacheRepo       repositories.AccountRepository
}

func NewAccountService(persistenceRepo, cacheRepo repositories.AccountRepository) AccountService {
	return &accountService{
		persistenceRepo: persistenceRepo,
		cacheRepo:       cacheRepo,
	}
}

func (s *accountService) Insert(account *datamodels.Account) (insertedId int64, err error) {
	return s.persistenceRepo.Insert(account)
}

func (s *accountService) Delete(id int64) (rowsAffected int64, err error) {
	return s.persistenceRepo.Delete(id)
}

func (s *accountService) UpdateNickname(id int64, account *datamodels.Account) (rowsAffected int64, err error) {
	return s.persistenceRepo.UpdateNickname(id, account)
}

func (s *accountService) UpdatePhone(id int64, account *datamodels.Account) (rowsAffected int64, err error) {
	return s.persistenceRepo.UpdatePhone(id, account)
}

func (s *accountService) UpdatePassword(id int64, account *datamodels.Account) (rowsAffected int64, err error) {
	return s.persistenceRepo.UpdatePassword(id, account)
}

func (s *accountService) SelectAll() (accounts []interface{}, total int64, err error) {
	return s.persistenceRepo.SelectAll()
}

func (s *accountService) SelectAllPerPage(page int64, pageSize int64) (accounts []interface{}, total int64, err error) {
	return s.persistenceRepo.SelectAllPerPage(page, pageSize)
}

func (s *accountService) SelectAuthByNickname(nickname string) (account datamodels.Account, err error) {
	return s.persistenceRepo.SelectAuthByNickname(nickname)
}

func (s *accountService) SelectAuthByPhone(phone string) (account datamodels.Account, err error) {
	return s.persistenceRepo.SelectAuthByPhone(phone)
}

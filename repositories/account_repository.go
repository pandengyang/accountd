package repositories

import (
	"accountd/datamodels"
)

type AccountRepository interface {
	Insert(account *datamodels.Account) (insertedId int64, err error)
	Delete(id int64) (rowsAffected int64, err error)

	UpdateNickname(id int64, account *datamodels.Account) (rowsAffected int64, err error)
	UpdatePhone(id int64, account *datamodels.Account) (rowsAffected int64, err error)
	UpdatePassword(id int64, account *datamodels.Account) (rowsAffected int64, err error)

	SelectAll() (accounts []interface{}, total int64, err error)
	SelectAllPerPage(page int64, pageSize int64) (accounts []interface{}, total int64, err error)

	SelectAuthByNickname(nickname string) (account datamodels.Account, err error)
	SelectAuthByPhone(phone string) (account datamodels.Account, err error)
}

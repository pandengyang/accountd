package repositories

import (
	"accountd/datamodels"
	"database/sql"
	"time"
)

type accountMySQLRepository struct {
	Db *sql.DB
}

func NewAccountMySQLRepository(db *sql.DB) AccountRepository {
	return &accountMySQLRepository{
		Db: db,
	}
}

func (r *accountMySQLRepository) Insert(account *datamodels.Account) (insertedId int64, err error) {
	var stmt *sql.Stmt
	var result sql.Result

	if stmt, err = r.Db.Prepare("INSERT INTO account (nickname, phone, password, salt, state, created_at) VALUES (?, ?, ?, ?, 'A', ?)"); err != nil {
		return insertedId, err
	}

	if result, err = stmt.Exec(account.Nickname, account.Phone, account.Password, account.Salt, time.Now().Unix()); err != nil {
		return insertedId, err
	}

	if insertedId, err = result.LastInsertId(); err != nil {
		return insertedId, err
	}

	return insertedId, err
}

func (r *accountMySQLRepository) Delete(id int64) (rowsAffected int64, err error) {
	var stmt *sql.Stmt
	var result sql.Result

	if stmt, err = r.Db.Prepare("DELETE FROM account WHERE id=?"); err != nil {
		return rowsAffected, err
	}

	if result, err = stmt.Exec(id); err != nil {
		return rowsAffected, err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return rowsAffected, err
	}

	return rowsAffected, err
}

func (r *accountMySQLRepository) UpdateNickname(id int64, account *datamodels.Account) (rowsAffected int64, err error) {
	var stmt *sql.Stmt
	var result sql.Result

	if stmt, err = r.Db.Prepare("UPDATE account SET nickname=? WHERE id=?"); err != nil {
		return rowsAffected, err
	}

	if result, err = stmt.Exec(account.Nickname, id); err != nil {
		return rowsAffected, err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return rowsAffected, err
	}

	return rowsAffected, err
}

func (r *accountMySQLRepository) UpdatePhone(id int64, account *datamodels.Account) (rowsAffected int64, err error) {
	var stmt *sql.Stmt
	var result sql.Result

	if stmt, err = r.Db.Prepare("UPDATE account SET phone=? WHERE id=?"); err != nil {
		return rowsAffected, err
	}

	if result, err = stmt.Exec(account.Phone, id); err != nil {
		return rowsAffected, err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return rowsAffected, err
	}

	return rowsAffected, err
}

func (r *accountMySQLRepository) UpdatePassword(id int64, account *datamodels.Account) (rowsAffected int64, err error) {
	var stmt *sql.Stmt
	var result sql.Result

	if stmt, err = r.Db.Prepare("UPDATE account SET password=? WHERE id=?"); err != nil {
		return rowsAffected, err
	}

	if result, err = stmt.Exec(account.Password, id); err != nil {
		return rowsAffected, err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		return rowsAffected, err
	}

	return rowsAffected, err
}

func (r *accountMySQLRepository) SelectAuthByNickname(nickname string) (account datamodels.Account, err error) {
	account = datamodels.Account{}

	sqlStatement := "SELECT id, nickname, phone, password, salt, state FROM account WHERE nickname=?"
	err = r.Db.QueryRow(sqlStatement, nickname).Scan(&account.Id, &account.Nickname, &account.Phone, &account.Password, &account.Salt, &account.State)

	return account, err
}

func (r *accountMySQLRepository) SelectAuthByPhone(phone string) (account datamodels.Account, err error) {
	account = datamodels.Account{}

	sqlStatement := "SELECT id, nickname, phone, password, salt, state FROM account WHERE phone=?"
	err = r.Db.QueryRow(sqlStatement, phone).Scan(&account.Id, &account.Nickname, &account.Phone, &account.Password, &account.Salt, &account.State)

	return account, err
}

func (r *accountMySQLRepository) SelectAll() (accounts []interface{}, total int64, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	sqlStatement := "SELECT id, nickname, phone, state, created_at FROM account ORDER BY created_at DESC"
	if stmt, err = r.Db.Prepare(sqlStatement); err != nil {
		return accounts, total, err
	}

	if rows, err = stmt.Query(); err != nil {
		return accounts, total, err
	}

	for rows.Next() {
		account := datamodels.Account{}

		if err = rows.Scan(&account.Id, &account.Nickname, &account.Phone, &account.State, &account.CreatedAt); err != nil {
			return accounts, total, err
		}

		accounts = append(accounts, account)
	}
	rows.Close()

	sqlStatement = "SELECT COUNT(id) AS total FROM account"
	if err = r.Db.QueryRow(sqlStatement).Scan(&total); err != nil {
		return accounts, total, err
	}

	return accounts, total, err
}

func (r *accountMySQLRepository) SelectAllPerPage(page int64, pageSize int64) (accounts []interface{}, total int64, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	sqlStatement := "SELECT id, nickname, phone, state, created_at FROM account ORDER BY created_at DESC LIMIT ? OFFSET ?"
	if stmt, err = r.Db.Prepare(sqlStatement); err != nil {
		return accounts, total, err
	}

	if rows, err = stmt.Query(pageSize, (page-1)*pageSize); err != nil {
		return accounts, total, err
	}

	for rows.Next() {
		account := datamodels.Account{}

		if err = rows.Scan(&account.Id, &account.Nickname, &account.Phone, &account.State, &account.CreatedAt); err != nil {
			return accounts, total, err
		}

		accounts = append(accounts, account)
	}
	rows.Close()

	sqlStatement = "SELECT COUNT(id) as total FROM account"
	if err = r.Db.QueryRow(sqlStatement).Scan(&total); err != nil {
		return accounts, total, err
	}

	return accounts, total, err
}

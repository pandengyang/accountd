package repositories

import (
	"accountd/datamodels"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type verificationCodeRedisRepository struct {
	Db *redis.Pool
}

func NewVerificationCodeRedisRepository(db *redis.Pool) VerificationCodeRepository {
	return &verificationCodeRedisRepository{
		Db: db,
	}
}

func (r *verificationCodeRedisRepository) Insert(vc *datamodels.VerificationCode) (insertedPhone string, err error) {
	conn := r.Db.Get()
	defer conn.Close()

	insertedPhone = vc.Phone

	vcKey := fmt.Sprintf("vc:%s", vc.Phone)
	if _, err = conn.Do("HMSET", vcKey, "c", vc.VerificationCode, "s", vc.SentAt); err != nil {
		return insertedPhone, err
	}

	if _, err = conn.Do("EXPIRE", vcKey, datamodels.VerificationCodeExpire); err != nil {
		return insertedPhone, err
	}

	return insertedPhone, err
}

func (r *verificationCodeRedisRepository) SelectByPhone(phone string) (vc datamodels.VerificationCode, err error) {
	conn := r.Db.Get()
	defer conn.Close()

	vcKey := fmt.Sprintf("vc:%s", phone)

	var values []interface{}
	if values, err = redis.Values(conn.Do("HMGET", vcKey, "c", "s")); err != nil {
		return vc, err
	}

	if _, err = redis.Scan(values, &vc.VerificationCode, &vc.SentAt); err != nil {
		return vc, err
	}
	vc.Phone = phone

	return vc, err
}

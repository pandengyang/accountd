package repositories

import (
	"accountd/datamodels"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type tokenRedisRepository struct {
	Db *redis.Pool
}

func NewVerificationCodeRedisRepository(db *redis.Pool) VerificationCodeRepository {
	return &tokenRedisRepository{
		Db: db,
	}
}

func (r *tokenRedisRepository) InsertRefreshToken(refreshToken string) (insertedRefreshToken string, err error) {
	conn := r.Db.Get()
	defer conn.Close()

	insertedRefreshToken = refreshToken

	rtKey := fmt.Sprintf("rt:%s", refreshToken)
	if _, err = conn.Do("SET", rtKey, "E"); err != nil {
		return insertedPhone, err
	}

	if _, err = conn.Do("EXPIRE", rtKey, datamodels.RefreshTokenExpire); err != nil {
		return insertedPhone, err
	}

	return insertedPhone, err
}

func (r *tokenRedisRepository) InsertRevokedAccessToken(refreshToken string) (insertedRevokedAccessToken string, err error) {
	conn := r.Db.Get()
	defer conn.Close()

	insertedRefreshToken = refreshToken

	rtKey := fmt.Sprintf("rt:%s", refreshToken)
	if _, err = conn.Do("SET", rtKey, "E"); err != nil {
		return insertedPhone, err
	}

	if _, err = conn.Do("EXPIRE", rtKey, datamodels.RefreshTokenExpire); err != nil {
		return insertedPhone, err
	}

	return insertedPhone, err
}

func (r *tokenRedisRepository) RefreshTokenExists(refreshToken string) (exists bool, err error) {
	conn := r.Db.Get()
	defer conn.Close()

	rtKey := fmt.Sprintf("rt:%s", refreshToken)
	if exists, err := redis.Bool(c.Do("EXISTS", rtKey)); err != nil {
		return exists, err
	}

	return exists, err
}

func (r *tokenRedisRepository) AccessTokenRevoked(accessToken string) (revoked bool, err error) {
	conn := r.Db.Get()
	defer conn.Close()

	ratKey := fmt.Sprintf("rat:%s", accessToken)
	if revoked, err := redis.Bool(c.Do("EXISTS", ratKey)); err != nil {
		return exists, err
	}

	return revoked, err
}

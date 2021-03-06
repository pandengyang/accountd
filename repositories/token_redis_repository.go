package repositories

import (
	"accountd/datamodels"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type tokenRedisRepository struct {
	Db *redis.Pool
}

func NewTokenRedisRepository(db *redis.Pool) TokenRepository {
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
		return insertedRefreshToken, err
	}

	if _, err = conn.Do("EXPIRE", rtKey, datamodels.REFRESH_TOKEN_EXPIRE); err != nil {
		return insertedRefreshToken, err
	}

	return insertedRefreshToken, err
}

func (r *tokenRedisRepository) InsertRevokedAccessToken(accessToken string) (insertedRevokedAccessToken string, err error) {
	conn := r.Db.Get()
	defer conn.Close()

	insertedRevokedAccessToken = accessToken

	ratKey := fmt.Sprintf("rat:%s", accessToken)
	if _, err = conn.Do("SET", ratKey, "E"); err != nil {
		return insertedRevokedAccessToken, err
	}

	if _, err = conn.Do("EXPIRE", ratKey, datamodels.EFFECTIVE_TIME); err != nil {
		return insertedRevokedAccessToken, err
	}

	return insertedRevokedAccessToken, err
}

func (r *tokenRedisRepository) RefreshTokenExists(refreshToken string) (exists bool, err error) {
	conn := r.Db.Get()
	defer conn.Close()

	rtKey := fmt.Sprintf("rt:%s", refreshToken)
	fmt.Println(rtKey)
	if exists, err = redis.Bool(conn.Do("EXISTS", rtKey)); err != nil {
		return exists, err
	}
	fmt.Println(exists)

	return exists, err
}

func (r *tokenRedisRepository) AccessTokenRevoked(accessToken string) (revoked bool, err error) {
	conn := r.Db.Get()
	defer conn.Close()

	ratKey := fmt.Sprintf("rat:%s", accessToken)
	if revoked, err = redis.Bool(conn.Do("EXISTS", ratKey)); err != nil {
		return revoked, err
	}

	return revoked, err
}

func (r *tokenRedisRepository) DeleteRefreshToken(refreshToken string) (rowsAffected int64, err error) {
	conn := r.Db.Get()
	defer conn.Close()

	rtKey := fmt.Sprintf("rt:%s", refreshToken)
	if _, err = conn.Do("DEL", rtKey); err != nil {
		return rowsAffected, err
	}

	return rowsAffected, err
}

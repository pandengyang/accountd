package services

import (
	"accountd/repositories"
)

type TokenService interface {
	InsertRefreshToken(refreshToken string) (insertedRefreshToken string, err error)
	DeleteRefreshToken(refreshToken string) (rowsAffected int64, err error)

	InsertRevokedAccessToken(refreshToken string) (insertedRevokedAccessToken string, err error)

	RefreshTokenExists(refreshToken string) (exists bool, err error)
	AccessTokenRevoked(accessToken string) (revoked bool, err error)
}

type tokenService struct {
	persistenceRepo repositories.TokenRepository
	cacheRepo       repositories.TokenRepository
}

func NewTokenService(persistenceRepo, cacheRepo repositories.TokenRepository) TokenService {
	return &tokenService{
		persistenceRepo: persistenceRepo,
		cacheRepo:       cacheRepo,
	}
}

func (s *tokenService) InsertRefreshToken(refreshToken string) (insertedRefreshToken string, err error) {
	return s.cacheRepo.InsertRefreshToken(refreshToken)
}

func (s *tokenService) InsertRevokedAccessToken(accessToken string) (insertedRevokedAccessToken string, err error) {
	return s.cacheRepo.InsertRevokedAccessToken(accessToken)
}

func (s *tokenService) RefreshTokenExists(refreshToken string) (exists bool, err error) {
	return s.cacheRepo.RefreshTokenExists(refreshToken)
}

func (s *tokenService) AccessTokenRevoked(accessToken string) (revoked bool, err error) {
	return s.cacheRepo.AccessTokenRevoked(accessToken)
}

func (s *tokenService) DeleteRefreshToken(refreshToken string) (rowsAffected int64, err error) {
	return s.cacheRepo.DeleteRefreshToken(refreshToken)
}

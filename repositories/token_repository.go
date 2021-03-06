package repositories

type TokenRepository interface {
	InsertRefreshToken(refreshToken string) (insertedRefreshToken string, err error)
	DeleteRefreshToken(refreshToken string) (rowsAffected int64, err error)

	InsertRevokedAccessToken(accessToken string) (insertedRevokedAccessToken string, err error)

	RefreshTokenExists(refreshToken string) (exists bool, err error)
	AccessTokenRevoked(accessToken string) (revoked bool, err error)
}

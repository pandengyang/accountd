package repositories

type TokenRepository interface {
	InsertRefreshToken(refreshToken string) (insertedRefreshToken string, err error)

	RefreshTokenExists(refreshToken string) (exists bool, err error)
	AccessTokenRevoked(accessToken string) (revoked bool, err error)
}

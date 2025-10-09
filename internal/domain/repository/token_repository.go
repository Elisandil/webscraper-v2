package repository

import "time"

type TokenRepository interface {
	RevokeToken(token string, expiresAt time.Time) error
	IsTokenRevoked(token string) (bool, error)
	CleanupExpiredTokens() error
}

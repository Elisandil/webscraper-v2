package persistence

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
	"webscraper-v2/internal/domain/repository"
	"webscraper-v2/internal/infrastructure/database"
)

type InMemoryTokenRepository struct {
	blacklist map[string]time.Time
	mu        sync.RWMutex
}

func NewInMemoryTokenRepository() repository.TokenRepository {
	return &InMemoryTokenRepository{
		blacklist: make(map[string]time.Time),
	}
}

func (r *InMemoryTokenRepository) RevokeToken(token string, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.blacklist[token] = expiresAt
	return nil
}

func (r *InMemoryTokenRepository) IsTokenRevoked(token string) (bool, error) {
	r.mu.RLock()
	expiry, exists := r.blacklist[token]
	r.mu.RUnlock()

	if !exists {
		return false, nil
	}

	if time.Now().Before(expiry) {
		return true, nil
	}

	r.mu.Lock()
	delete(r.blacklist, token)
	r.mu.Unlock()

	return false, nil
}

func (r *InMemoryTokenRepository) CleanupExpiredTokens() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for token, expiry := range r.blacklist {
		if now.After(expiry) {
			delete(r.blacklist, token)
		}
	}
	return nil
}

type SQLiteTokenRepository struct {
	db *database.SQLiteDB
}

func NewSQLiteTokenRepository(db *database.SQLiteDB) repository.TokenRepository {
	repo := &SQLiteTokenRepository{db: db}
	if err := repo.createTable(); err != nil {
		fmt.Printf("Warning: error creating token table: %v\n", err)
	}
	return repo
}

func (r *SQLiteTokenRepository) createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS revoked_tokens (
		token TEXT PRIMARY KEY,
		expires_at DATETIME NOT NULL,
		revoked_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_revoked_tokens_expires_at ON revoked_tokens(expires_at);
	`
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating revoked_tokens table: %w", err)
	}
	return nil
}

func (r *SQLiteTokenRepository) RevokeToken(token string, expiresAt time.Time) error {
	query := `INSERT OR REPLACE INTO revoked_tokens (token, expires_at) VALUES (?, ?)`
	_, err := r.db.Exec(query, token, expiresAt)
	if err != nil {
		return fmt.Errorf("error revoking token: %w", err)
	}
	return nil
}

func (r *SQLiteTokenRepository) IsTokenRevoked(token string) (bool, error) {
	query := `SELECT COUNT(*) FROM revoked_tokens WHERE token = ? AND expires_at > ?`
	var count int
	err := r.db.QueryRow(query, token, time.Now().UTC()).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error checking revoked token: %w", err)
	}
	return count > 0, nil
}

func (r *SQLiteTokenRepository) CleanupExpiredTokens() error {
	query := `DELETE FROM revoked_tokens WHERE expires_at < ?`
	_, err := r.db.Exec(query, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("error cleaning up tokens: %w", err)
	}
	return nil
}

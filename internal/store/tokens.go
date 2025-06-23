package store

import (
	"database/sql"
	"time"

	"github.com/DerekBelloni/fem_project/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgrestTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error)
}

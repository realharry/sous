package storage

import (
	"database/sql"
	"fmt"

	// it's a SQL db driver. This is how you do that.
	_ "github.com/lib/pq"
)

type (
	// The PostgresStateManager provides the StateManager interface by
	// reading/writing from a postgres database.
	PostgresStateManager struct {
		db *sql.DB
	}

	// A PostgresConfig describes how to connect to a postgres database
	PostgresConfig struct {
		DBName   string
		User     string
		Password string
		Host     string
		Port     string
	}
)

// NewPostgresStateManager creates a new PostgresStateManager.
func NewPostgresStateManager(cfg PostgresConfig) (*PostgresStateManager, error) {
	db, err := sql.Open("postgres", cfg.connStr())
	if err != nil {
		return nil, err
	}
	return &PostgresStateManager{db: db}, nil
}

func (c PostgresConfig) connStr() string {
	return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s", c.DBName, c.User, c.Password, c.Host, c.Port)
}

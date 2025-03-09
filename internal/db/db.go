package db

import (
	_ "database/sql"
	"fmt"

	"github.com/dmytro-vovk/shta/internal/types"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	db *sqlx.DB
}

func New(dsn string) (*DB, error) {
	conn, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &DB{
		db: conn,
	}, nil
}

func (db *DB) UpsertURL(url string) error {
	_, err := db.db.Exec(`
		INSERT INTO urls (url) VALUES ($1)
		ON CONFLICT (url) DO
		UPDATE SET last_seen = NOW()`,
		url,
	)

	return err
}

func (db *DB) FetchURLs(limit int) ([]*types.URLRecord, error) {
	var data []*types.URLRecord

	if err := db.db.Select(
		&data,
		`SELECT url, last_seen FROM urls ORDER BY last_seen DESC LIMIT $1`,
		limit,
	); err != nil {
		return nil, err
	}

	return data, nil
}

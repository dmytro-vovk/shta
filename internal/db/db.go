package db

import (
	_ "database/sql"
	"fmt"
	"time"

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
		INSERT INTO urls (
			url
		) VALUES (
			$1
		)
		ON CONFLICT (url) DO UPDATE 
			SET seen = urls.seen + 1`,
		url,
	)

	return err
}

type URLRecord struct {
	Created time.Time `db:"created_at" json:"created"`
	URL     string    `db:"url"        json:"url"`
	Seen    int64     `db:"seen"       json:"seen"`
}

func (db *DB) FetchURLs(sortBy, order string, limit int) ([]URLRecord, error) {
	var data []URLRecord

	if err := db.db.Select(
		&data,
		`SELECT url, seen, created_at FROM urls ORDER BY `+sortBy+` `+order+` LIMIT $1`,
		limit,
	); err != nil {
		return nil, fmt.Errorf("fetch urls: %w", err)
	}

	return data, nil
}

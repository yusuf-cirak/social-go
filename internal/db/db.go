package db

import (
	"context"
	"database/sql"
	"time"
)

const QueryTimeout = 5 * time.Second

// DB wrapper struct - sql.DB'yi wrap eder
type DB struct {
	db *sql.DB
}

func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	if maxIdleTime != "" {
		duration, err := time.ParseDuration(maxIdleTime)
		if err != nil {
			return nil, err
		}
		db.SetConnMaxIdleTime(duration)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (d *DB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	return d.db.QueryContext(ctx, query, args...)
}

func (d *DB) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *DB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	return d.db.ExecContext(ctx, query, args...)
}

func (d *DB) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	return d.db.PrepareContext(ctx, query)
}

func (d *DB) Close() error {
	return d.db.Close()
}

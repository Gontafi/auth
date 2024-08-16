package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func DBinit(ctx context.Context) (*pgx.Conn, error) {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS tokens (
			id UUID PRIMARY KEY,
			user_id UUID,
			refresh_token_hash VARCHAR(255),
			client_ip VARCHAR(20),
			used BOOLEAN,
			created_at TIMESTAMP
		);`)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

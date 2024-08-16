package repo

import (
	"auth/internal/models"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepo(db *pgx.Conn) *Repository {
	return &Repository{db: db}
}

func (r Repository) AddToken(ctx context.Context, token models.RefreshToken) error {
	query := `
		INSERT INTO tokens (id, user_id, refresh_token_hash, client_ip, used, created_at)
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err := r.db.Exec(ctx, query,
		uuid.New(),
		token.UserID,
		token.RefreshTokenHash,
		token.ClientIP,
		token.Used,
		time.Now(),
	)

	return err
}

func (r Repository) UpdateTokenUsedInfo(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE tokens
		SET used = true
		WHERE id = $1
	`, id)

	return err
}
func (r Repository) GetToken(ctx context.Context, userID string) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, refresh_token_hash, client_ip, used, created_at
		FROM tokens 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1;
	`

	row := r.db.QueryRow(ctx, query, userID)

	var token models.RefreshToken
	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.RefreshTokenHash,
		&token.ClientIP,
		&token.Used,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

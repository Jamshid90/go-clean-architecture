package refreshtoken

import (
	"fmt"
	"context"
	"github.com/Jamshid90/go-clean-architecture/pkg/errors"
	"github.com/jackc/pgx/v4"
	"github.com/Jamshid90/go-clean-architecture/pkg/entity"
	"github.com/jackc/pgx/v4/pgxpool"
)

type pgxRefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepositoryPgx(dbpool *pgxpool.Pool) entity.RefreshTokenRepository {
	return &pgxRefreshTokenRepository{db: dbpool}
}

func (p *pgxRefreshTokenRepository) Store(ctx context.Context, m *entity.RefreshToken) error {
	_, err := p.db.Exec(ctx,`INSERT INTO "refresh_token"(
		user_id, token, created_at)
		VALUES ($1, $2, $3);`,
		m.UserID,
		m.Token,
		m.CreatedAt,
	)

	if err != nil {
		return errors.ErrRepository{Err:fmt.Errorf("error during store to refresh token repository: %w", err)}
	}

	return nil
}

func (p *pgxRefreshTokenRepository) Delete(ctx context.Context, token string) error {
	if _, err := p.db.Exec(ctx,`DELETE FROM "refresh_token" WHERE token=$1`, token); err != nil {
		return errors.ErrRepository{Err:fmt.Errorf("error during delete to refresh token repository: %w", err)}
	}
	return nil
}

func (p *pgxRefreshTokenRepository) DeleteByUserId(ctx context.Context, id string) error {
	if _, err := p.db.Exec(ctx,`DELETE FROM "refresh_token" WHERE user_id=$1`, id); err != nil {
		return errors.ErrRepository{Err:fmt.Errorf("error during delete by user id to refresh token repository: %w", err)}
	}
	return nil
}

func (p *pgxRefreshTokenRepository) Find(ctx context.Context, token string) (*entity.RefreshToken, error) {
	refreshToken := entity.RefreshToken{}
	row := p.db.QueryRow(ctx,`SELECT user_id, token, created_at FROM "refresh_token" WHERE token=$1`, token)

	err := row.Scan(
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.NewErrNotFound("refresh token")
	}

	if err != nil {
		return nil, errors.ErrRepository{Err:fmt.Errorf("error during find to refresh token repository: %w", err)}
	}

	return &refreshToken, nil
}
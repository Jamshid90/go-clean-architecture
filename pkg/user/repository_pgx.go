package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
	"github.com/jackc/pgx/v4"
)

type pgxUserRepository struct {
	Conn *pgx.Conn
}

func NewPgxUserRepository(conn *pgx.Conn) domain.UserRepository {
	return &pgxUserRepository{Conn: conn}
}

func (p *pgxUserRepository) Store(ctx context.Context, m *domain.User) error {
	err := p.Conn.QueryRow(ctx,`INSERT INTO public."user"(
		status, email, first_name, last_name, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`,
		m.Status,
		m.Email,
		m.FirstName,
		m.LastName,
		m.Password,
		m.CreatedAt,
		m.UpdatedAt,
	).Scan(&m.ID)

	if err != nil {
		return domain.ErrRepository{Err:fmt.Errorf("error during store to user repository: %w", err)}
	}

	return nil
}

func (p *pgxUserRepository) Update(ctx context.Context, m *domain.User) error {
	_, err := p.Conn.Exec(ctx,`UPDATE "user" SET status=$1, email=$2, first_name=$3, last_name=$4, updated_at=$5 WHERE id=$6`,
		m.Status,
		m.Email,
		m.FirstName,
		m.LastName,
		m.UpdatedAt,
		m.ID,
	)

	if err != nil {
		return domain.ErrRepository{Err:fmt.Errorf("error during store to user repository: %w", err)}
	}

	return nil
}

func (p *pgxUserRepository) Delete(ctx context.Context, id int64) error {
	if _, err := p.Conn.Exec(ctx,`DELETE FROM "user" WHERE id=$1`, id); err != nil {
		return domain.ErrRepository{Err:fmt.Errorf("error during delete to user repository: %w", err)}
	}
	return nil
}

func (p *pgxUserRepository) Find(ctx context.Context, id int64) (*domain.User, error) {
	user := domain.User{}
	row := p.Conn.QueryRow(ctx,`SELECT id, status, email, first_name, last_name, password, created_at, updated_at FROM "user" WHERE id=$1`, id)

	err := row.Scan(
		&user.ID,
		&user.Status,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.NewErrNotFound("user")
	}

	if err != nil {
		return nil, domain.ErrRepository{Err:fmt.Errorf("error during find to user repository: %w", err)}
	}

	return &user, nil
}

func (p *pgxUserRepository) FindAll(ctx context.Context, limit, offset int, params map[string]interface{}) ([]*domain.User, error) {
	var items []*domain.User
	rows, err := p.Conn.Query(ctx, `SELECT id, status, email, first_name, last_name, created_at, updated_at FROM "user" LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return items, domain.ErrRepository{Err:fmt.Errorf("error during find all to user repository: %w", err)}
	}
	for rows.Next() {
		user := domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Status,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return items, domain.ErrRepository{Err:fmt.Errorf("error during find all to user repository: %w", err)}
		}
		items = append(items, &user)
	}
	return items, nil
}

func (p *pgxUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {

	user := domain.User{}
	row := p.Conn.QueryRow(ctx, `SELECT id, status, email, first_name, last_name, created_at, updated_at FROM "user" WHERE email=$1`, email)
	err := row.Scan(
		&user.ID,
		&user.Status,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.NewErrNotFound("user")
	}

	if err != nil {
		return nil, domain.ErrRepository{Err:fmt.Errorf("error during find by email to user repository: %w", err)}
	}

	return &user, nil
}
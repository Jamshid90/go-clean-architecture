package user

import (
	"fmt"
	"database/sql"
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
)

type postgresqlUserRepository struct {
	Conn *sql.DB
}

func NewPostgresqlUserRepository(conn *sql.DB) domain.UserRepository {
	return &postgresqlUserRepository{Conn: conn}
}

func (p *postgresqlUserRepository) Store(m *domain.User) error {
	err := p.Conn.QueryRow(`INSERT INTO public."user"(
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

func (p *postgresqlUserRepository) Update(m *domain.User) error {
	_, err := p.Conn.Exec(`UPDATE "user" SET status=$1, email=$2, first_name=$3, last_name=$4, updated_at=$5 WHERE id=$6`,
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

func (p *postgresqlUserRepository) Delete(id uint64) error {
	if _, err := p.Conn.Exec(`DELETE FROM "user" WHERE id=$1`, id); err != nil {
		return domain.ErrRepository{Err:fmt.Errorf("error during delete to user repository: %w", err)}
	}
	return nil
}

func (p *postgresqlUserRepository) Find(id uint64) (*domain.User, error) {
	user := domain.User{}
	row := p.Conn.QueryRow(`SELECT id, status, email, first_name, last_name, password, created_at, updated_at FROM "user" WHERE id=$1`, id)

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

func (p *postgresqlUserRepository) FindAll(limit, offset int, params map[string]interface{}) ([]*domain.User, error) {
	var items []*domain.User
	rows, err := p.Conn.Query(`SELECT id, status, email, first_name, last_name, created_at, updated_at FROM "user" LIMIT $1 OFFSET $2`, limit, offset)
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

func (p *postgresqlUserRepository) FindByEmail(email string) (*domain.User, error) {

	user := domain.User{}
	row := p.Conn.QueryRow(`SELECT id, status, email, first_name, last_name, created_at, updated_at FROM "user" WHERE email=$1`, email)
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
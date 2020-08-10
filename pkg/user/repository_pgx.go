package user

import (
	"fmt"
	"context"
	"github.com/Jamshid90/go-clean-architecture/pkg/errors"
	"github.com/jackc/pgx/v4"
	"github.com/Jamshid90/go-clean-architecture/pkg/entity"
	"github.com/jackc/pgx/v4/pgxpool"
)

type pgxUserRepository struct {
	db *pgxpool.Pool
}

func NewPgxUserRepository(dbpool *pgxpool.Pool) entity.UserRepository {
	return &pgxUserRepository{db: dbpool}
}

func (p *pgxUserRepository) Store(ctx context.Context, m *entity.User) error {
	_, err := p.db.Exec(ctx,`INSERT INTO "user"(
		id, status, email, phone, gender, first_name, last_name, password, birth_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		m.ID,
		m.Status,
		m.Email,
		m.Phone,
		m.Gender,
		m.FirstName,
		m.LastName,
		m.Password,
		m.BirthDate,
		m.CreatedAt,
		m.UpdatedAt,
	)

	if err != nil {
		return errors.ErrRepository{Err:fmt.Errorf("error during store to user repository: %w", err)}
	}

	return nil
}

func (p *pgxUserRepository) Update(ctx context.Context, m *entity.User) error {
	_, err := p.db.Exec(ctx,`UPDATE "user" 
	    SET status=$1, email=$2, phone=$3, gender=$4, first_name=$5, last_name=$6, birth_date=$7 updated_at=$8
	    WHERE id=$6`,
		m.Status,
		m.Email,
		m.Phone,
		m.Gender,
		m.FirstName,
		m.LastName,
		m.BirthDate,
		m.UpdatedAt,
		m.ID,
	)

	if err != nil {
		return errors.ErrRepository{Err:fmt.Errorf("error during store to user repository: %w", err)}
	}

	return nil
}

func (p *pgxUserRepository) Delete(ctx context.Context, id string) error {
	if _, err := p.db.Exec(ctx,`DELETE FROM "user" WHERE id=$1`, id); err != nil {
		return errors.ErrRepository{Err:fmt.Errorf("error during delete to user repository: %w", err)}
	}
	return nil
}

func (p *pgxUserRepository) Find(ctx context.Context, id string) (*entity.User, error) {
	user := entity.User{}
	row := p.db.QueryRow(ctx,`SELECT id, status, email, phone, gender, first_name, last_name, password, birth_date, created_at, updated_at 
                                   FROM "user" 
                                   WHERE id=$1`, id)

	err := row.Scan(
		&user.ID,
		&user.Status,
		&user.Email,
		&user.Phone,
		&user.Gender,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.BirthDate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.NewErrNotFound("user")
	}

	if err != nil {
		return nil, errors.ErrRepository{Err:fmt.Errorf("error during find to user repository: %w", err)}
	}

	return &user, nil
}

func (p *pgxUserRepository) FindAll(ctx context.Context, limit, offset int, params map[string]interface{}) ([]*entity.User, error) {
	var items []*entity.User
	rows, err := p.db.Query(ctx, `SELECT id, status, email, phone, gender, first_name, last_name, password, birth_date, created_at, updated_at
                                       FROM "user" 
 								       LIMIT $1 
									   OFFSET $2`, limit, offset)
	if err != nil {
		return items, errors.ErrRepository{Err:fmt.Errorf("error during find all to user repository: %w", err)}
	}
	for rows.Next() {
		user := entity.User{}
		err := rows.Scan(
			&user.ID,
			&user.Status,
			&user.Email,
			&user.Phone,
			&user.Gender,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.BirthDate,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return items, errors.ErrRepository{Err:fmt.Errorf("error during find all to user repository: %w", err)}
		}
		items = append(items, &user)
	}
	return items, nil
}

func (p *pgxUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := entity.User{}
	row := p.db.QueryRow(ctx, `SELECT id, status, email, phone, gender, first_name, last_name, password, birth_date, created_at, updated_at
 							        FROM "user"
  							        WHERE email=$1`, email)
	err := row.Scan(
		&user.ID,
		&user.Status,
		&user.Email,
		&user.Phone,
		&user.Gender,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.BirthDate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.NewErrNotFound("user")
	}

	if err != nil {
		return nil, errors.ErrRepository{Err:fmt.Errorf("error during find by email to user repository: %w", err)}
	}

	return &user, nil
}

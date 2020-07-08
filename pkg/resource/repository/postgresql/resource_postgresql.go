package postgresql

import (
	"database/sql"
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
)

type postgresqlResourceRepository struct {
	Conn *sql.DB
}

func NewPostgresqlResourceRepository(conn *sql.DB) domain.ResourceRepository {
	return &postgresqlResourceRepository{Conn: conn}
}

func (p *postgresqlResourceRepository) Store(r *domain.Resource) error {

	return nil
}

func (p *postgresqlResourceRepository) Update(r *domain.Resource) error {

	return nil
}

func (p *postgresqlResourceRepository) Delete(id uint64) error {

	return nil
}

func (p *postgresqlResourceRepository) Find(id uint64) (*domain.Resource, error) {

	resource := domain.Resource{}
	return &resource, nil
}

func (p *postgresqlResourceRepository) FindAll(limit, offset int, params map[string]interface{}) ([]*domain.Resource, error) {

	var resources []*domain.Resource
	return resources, nil
}
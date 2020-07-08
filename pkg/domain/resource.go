package domain

type Resource struct {
	ID             uint64      `json:"id,omitempty" db:"id"`
	ResourceID     string      `json:"resource_id,omitempty" db:"resource_id"`
	PractitionerID string      `json:"practitioner_id,omitempty" db:"practitioner_id"`
	ResourceType   string      `json:"resource_type,omitempty" db:"resource_type"`
	Clinic         string      `json:"clinic,omitempty" db:"clinic"`
	CreatedAt      string      `json:"created_at,omitempty" db:"created_at"`
	Data           interface{} `json:"data,omitempty" db:"data"`
}

type ResourceUsecase interface {
	Store(r *Resource) error
	Update(r *Resource) error
	Delete(id uint64) error
	Find (id uint64) (*Resource, error)
	FindAll(limit, offset int, params map[string]interface{}) ([]*Resource, error)
}

type ResourceRepository interface {
	Store(r *Resource) error
	Update(r *Resource) error
	Delete(id uint64) error
	Find(id uint64) (*Resource, error)
	FindAll(limit, offset int, params map[string]interface{}) ([]*Resource, error)
}
package auth

import "time"

type User struct {
	ID        string    `json:"id,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Gender    string    `json:"gender,omitempty"`
	Status    string    `json:"status,omitempty"`
	Email     string    `json:"email,omitempty"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	BirthDate time.Time `json:"birth_date,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

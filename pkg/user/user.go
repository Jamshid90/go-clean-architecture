package user

import "time"

type User struct {
	ID        int64 `json:"id,omitempty"`
	Status    string `json:"status,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (u *User) Sanitize() *User {
	u.Password = ""
	return u
}

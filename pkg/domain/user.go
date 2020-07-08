package domain

import "time"

type User struct {
	ID          uint64
	Email       string
	Status      string
	FirstName   string
	LastName    string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserUsecase interface {
	Store(*User) error
	Update(*User) error
	Delete(uint64) error
	Find (uint64) (*User, error)
	FindAll(limit, offset int, params map[string]interface{}) ([]*User, error)
}

type UserRepository interface {
	Store(*User) error
	Update(*User) error
	Delete(uint64) error
	Find(uint64) (*User, error)
	FindAll(limit, offset int, params map[string]interface{}) ([]*User, error)
	FindByEmail (string) (*User, error)
}

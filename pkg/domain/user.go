package domain

import (
	"context"
	"time"
)

const (
	USER_STATUS_ACTIVE   = "active"
	USER_STATUS_DEACTIVE = "deactive"
)

type User struct {
	ID          int64
	Email       string
	Status      string
	FirstName   string
	LastName    string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserUsecase interface {
	Store(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
	Find (ctx context.Context, id int64) (*User, error)
	FindAll(ctx context.Context, limit, offset int, params map[string]interface{}) ([]*User, error)
	FindByEmail (ctx context.Context, email string) (*User, error)
}

type UserRepository interface {
	Store(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
	Find(ctx context.Context, id int64) (*User, error)
	FindAll(ctx context.Context, limit, offset int, params map[string]interface{}) ([]*User, error)
	FindByEmail (ctx context.Context, email string) (*User, error)
}

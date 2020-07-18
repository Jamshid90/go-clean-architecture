package user

import (
	"context"
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
	"time"
)

type userUsecase struct {
	userRepo domain.UserRepository
	contextTimeout time.Duration
}
// New User Usecase
func NewUserUsecase(userRepo domain.UserRepository, timeout time.Duration) userUsecase  {
	return userUsecase{
		userRepo: userRepo,
		contextTimeout: timeout,
	}
}

// Before Store
func (u *userUsecase) BeforeStore(m *domain.User) {
	m.CreatedAt = time.Now().UTC()
	m.UpdatedAt = m.CreatedAt
}

// Store
func (u *userUsecase) Store(ctx context.Context, m *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.FindByEmail(ctx ,m.Email)
	if err != nil {
		return err
	}

	if user != nil {
		return domain.NewErrConflict("email")
	}

	u.BeforeStore(m)
	return u.userRepo.Store(ctx, m)
}

// Update
func (u *userUsecase) Update(ctx context.Context, m *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.Find(ctx, m.ID)
	if err != nil {
		return err
	}

	if userByEmail, _ := u.userRepo.FindByEmail(ctx, m.Email); userByEmail != nil && userByEmail.ID != user.ID {
		return  domain.NewErrConflict("email")
	}
	m.CreatedAt = user.CreatedAt
	m.UpdatedAt = time.Now().UTC()
	return u.userRepo.Update(ctx, m)
}

// Delete
func (u *userUsecase) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	existedUser, err := u.userRepo.Find(ctx, id)
	if err != nil {
		return err
	}

	if *existedUser == (domain.User{}) {
		return domain.NewErrNotFound("user")
	}

	return u.userRepo.Delete(ctx, id)
}

// Find
func (u *userUsecase) Find(ctx context.Context, id int64) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.userRepo.Find(ctx, id)
}

// FindAll
func (u *userUsecase) FindAll(ctx context.Context, limit, offset int, params map[string]interface{}) ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.userRepo.FindAll(ctx, limit, offset, params)
}
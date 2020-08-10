package user

import (
	"context"
	"github.com/Jamshid90/go-clean-architecture/pkg/entity"
	"github.com/Jamshid90/go-clean-architecture/pkg/errors"
	"github.com/Jamshid90/go-clean-architecture/pkg/hash"
	"github.com/Jamshid90/go-clean-architecture/pkg/rand"
	"time"
)

type userUsecase struct {
	userRepo entity.UserRepository
	contextTimeout time.Duration
}
// new user usecase
func NewUserUsecase(repo entity.UserRepository, timeout time.Duration) userUsecase  {
	return userUsecase{
		userRepo: repo,
		contextTimeout: timeout,
	}
}

// get new id
func (u *userUsecase) NewID(ctx context.Context) (string, error) {
	var id = rand.RandString(16)
	user, err := u.Find(ctx, id)

	if err != nil && err.Error() != errors.NewErrNotFound("user").Error() {
		return "", err
	}

	if user != nil {
		return u.NewID(ctx)
	}

	return id, nil
}

// before store
func (u *userUsecase) BeforeStore(ctx context.Context, m *entity.User) error {

	m.CreatedAt = time.Now().UTC()
	m.UpdatedAt = m.CreatedAt

	hashPassword, err := hash.HashPassword(m.Password)
	if err != nil {
		return err
	}
	m.Password = hashPassword

	m.ID, err = u.NewID(ctx)
	if err != nil {
		return err
	}

	return nil
}

// store
func (u *userUsecase) Store(ctx context.Context, m *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.FindByEmail(ctx, m.Email)

	if err != nil && err.Error() != errors.NewErrNotFound("user").Error() {
		return err
	}

	if user != nil {
		return errors.NewErrConflict("email")
	}

	if err := u.BeforeStore(ctx, m); err != nil {
		return err
	}

	return u.userRepo.Store(ctx, m)
}

// update
func (u *userUsecase) Update(ctx context.Context, m *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.Find(ctx, m.ID)
	if err != nil {
		return err
	}

	if userByEmail, _ := u.userRepo.FindByEmail(ctx, m.Email); userByEmail != nil && userByEmail.ID != user.ID {
		return  errors.NewErrConflict("email")
	}
	m.CreatedAt = user.CreatedAt
	m.UpdatedAt = time.Now().UTC()
	return u.userRepo.Update(ctx, m)
}

// delete
func (u *userUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	existedUser, err := u.userRepo.Find(ctx, id)
	if err != nil {
		return err
	}

	if *existedUser == (entity.User{}) {
		return errors.NewErrNotFound("user")
	}

	return u.userRepo.Delete(ctx, id)
}

// find
func (u *userUsecase) Find(ctx context.Context, id string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.userRepo.Find(ctx, id)
}

// find all
func (u *userUsecase) FindAll(ctx context.Context, limit, offset int, params map[string]interface{}) ([]*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.userRepo.FindAll(ctx, limit, offset, params)
}

// find by email
func (u *userUsecase) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.userRepo.FindByEmail(ctx, email)
}

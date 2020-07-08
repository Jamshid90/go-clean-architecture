package user

import (
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
	"time"
)

type userUsecase struct {
	userRepo domain.UserRepository
}
// New User Usecase
func NewUserUsecase(userRepo domain.UserRepository) userUsecase  {
	return userUsecase{
		userRepo: userRepo,
	}
}

// Before Store
func (ru *userUsecase) BeforeStore(m *domain.User) {
	m.CreatedAt = time.Now().UTC()
	m.UpdatedAt = m.CreatedAt
}

// Store
func (ru *userUsecase) Store(m *domain.User) error {
	if user, _ := ru.userRepo.FindByEmail(m.Email); user != nil {
		return domain.NewErrConflict("email")
	}
	ru.BeforeStore(m)
	return ru.userRepo.Store(m)
}

// Update
func (ru *userUsecase) Update(m *domain.User) error {
	user, err := ru.Find(m.ID)
	if err != nil {
		return err
	}

	if userByEmail, _ := ru.userRepo.FindByEmail(m.Email); userByEmail != nil && userByEmail.ID != user.ID {
		return  domain.NewErrConflict("email")
	}
	m.CreatedAt = user.CreatedAt
	m.UpdatedAt = time.Now().UTC()
	return ru.userRepo.Update(m)
}

// Delete
func (ru *userUsecase) Delete(id uint64) error {
	return ru.userRepo.Delete(id)
}

// Find
func (ru *userUsecase) Find(id uint64) (*domain.User, error) {
	return ru.userRepo.Find(id)
}

// FindAll
func (ru *userUsecase) FindAll(limit, offset int, params map[string]interface{}) ([]*domain.User, error) {
	return ru.userRepo.FindAll(limit, offset, params)
}
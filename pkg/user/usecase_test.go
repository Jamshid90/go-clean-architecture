package user

import (
	"context"
	"errors"
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
	"github.com/Jamshid90/go-clean-architecture/pkg/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestBeforeStore(t *testing.T) {
	mockUser := TestUser(t)
	userUse := new(userUsecase)
	userUse.BeforeStore(mockUser)

	assert := assert.New(t)
	assert.NotEmpty(mockUser.CreatedAt)
	assert.NotEmpty(mockUser.UpdatedAt)
	assert.Equal(mockUser.CreatedAt, mockUser.UpdatedAt)
}

func TestStore(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := TestUser(t)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, nil).Once()
		mockUserRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		userUse := NewUserUsecase(mockUserRepo, time.Second*2)
		err := userUse.Store(context.TODO(), mockUser)

		assert.NoError(t, err)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-email-already-exist", func(t *testing.T) {
		mockUserRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(TestUserEmpty(t), domain.NewErrConflict("email")).Once()

		userUse := NewUserUsecase(mockUserRepo, time.Second*2)
		err := userUse.Store(context.TODO(), mockUser)

		assert := assert.New(t)
		assert.Error(err)
		assert.Equal(err, domain.NewErrConflict("email"))

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-happens-in-db", func(t *testing.T) {
		errRepository := domain.NewErrRepository(errors.New("Unexpected error"))

		mockUserRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, nil).Once()
		mockUserRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errRepository).Once()

		userUse := NewUserUsecase(mockUserRepo, time.Second*2)
		err := userUse.Store(context.TODO(), mockUser)

		assert := assert.New(t)
		assert.Error(err)
		assert.Equal(err, errRepository)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := TestUser(t)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Find", mock.Anything, mock.AnythingOfType("int64")).Return(mockUser, nil).Once()
		mockUserRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(mockUser, nil).Once()
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		userUse := NewUserUsecase(mockUserRepo, time.Second*2)
		err := userUse.Update(context.TODO(), mockUser)

		assert.NoError(t, err)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-not-found", func(t *testing.T) {
		mockUserRepo.On("Find", mock.Anything, mock.AnythingOfType("int64")).Return(TestUserEmpty(t), domain.NewErrNotFound("user")).Once()

		userUse := NewUserUsecase(mockUserRepo, time.Second*2)
		err := userUse.Update(context.TODO(), mockUser)

		assert := assert.New(t)
		assert.Error(err)
		assert.Equal(err, domain.NewErrNotFound("user"))

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-email-already-exist", func(t *testing.T) {

		mockUserRepo.On("Find", mock.Anything, mock.AnythingOfType("int64")).Return(mockUser, nil).Once()
		mockUserRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(TestUserEmpty(t), domain.NewErrConflict("email")).Once()

		userUse := NewUserUsecase(mockUserRepo, time.Second*2)
		err := userUse.Update(context.TODO(), mockUser)

		assert := assert.New(t)
		assert.Error(err)
		assert.Equal(err, domain.NewErrConflict("email"))

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-happens-in-db", func(t *testing.T) {
		errRepository := domain.NewErrRepository(errors.New("Unexpected error"))

		mockUserRepo.On("Find", mock.Anything, mock.AnythingOfType("int64")).Return(mockUser, nil).Once()
		mockUserRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(mockUser, nil).Once()
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errRepository).Once()

		userUse := NewUserUsecase(mockUserRepo, time.Second*2)
		err := userUse.Update(context.TODO(), mockUser)

		assert := assert.New(t)
		assert.Error(err)
		assert.Equal(err, errRepository)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {

	mockUserRepo := new(mocks.UserRepository)
	mockUser := TestUser(t)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Find", mock.Anything, mock.AnythingOfType("int64")).Return(mockUser, nil).Once()
		mockUserRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		userUse := NewUserUsecase(mockUserRepo, time.Second*2)
		err := userUse.Delete(context.TODO(), mockUser.ID)

		assert.NoError(t, err)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-not-found", func(t *testing.T) {
		mockUserRepo.On("Find", mock.Anything, mock.AnythingOfType("int64")).Return(TestUserEmpty(t), domain.NewErrNotFound("user")).Once()

		userUse := NewUserUsecase(mockUserRepo, time.Second*2)
		err := userUse.Delete(context.TODO(), mockUser.ID)

		assert := assert.New(t)
		assert.Error(err)
		assert.Equal(err, domain.NewErrNotFound("user"))

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-happens-in-db", func(t *testing.T) {
		errRepository := domain.NewErrRepository(errors.New("Unexpected error"))
		mockUserRepo.On("Find", mock.Anything, mock.AnythingOfType("int64")).Return(TestUserEmpty(t), errRepository).Once()

		userUse := NewUserUsecase(mockUserRepo, time.Second*2)
		err := userUse.Delete(context.TODO(), mockUser.ID)

		assert := assert.New(t)
		assert.Error(err)
		assert.Equal(err, errRepository)

		mockUserRepo.AssertExpectations(t)
	})

}

func TestFind(t *testing.T) {

	mockUserRepo := new(mocks.UserRepository)
	mockUser := TestUser(t)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Find", mock.Anything, mock.AnythingOfType("int64")).Return(mockUser, nil).Once()
		userUse := NewUserUsecase(mockUserRepo, time.Second*2)
		user, err := userUse.Find(context.TODO(), mockUser.ID)

		assert.NoError(t, err)
		assert.NotNil(t, user)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockUserRepo.On("Find", mock.Anything, mock.AnythingOfType("int64")).Return(&domain.User{}, errors.New("Unexpected error")).Once()

		userUse := NewUserUsecase(mockUserRepo, time.Second*2)
		user, err := userUse.Find(context.TODO(), mockUser.ID)

		assert.Error(t, err)
		assert.Equal(t, &domain.User{}, user)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestFindAll(t *testing.T) {

	mockUserRepo := new(mocks.UserRepository)
	mockUser := TestUser(t)

	mockListUser := make([]*domain.User, 0)
	mockListUser = append(mockListUser, mockUser)

	t.Run("success", func(t *testing.T) {

		mockUserRepo.On("FindAll",
			mock.Anything,
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(mockListUser, nil).Once()

		userUse := NewUserUsecase(mockUserRepo, time.Second*2)
		list, err := userUse.FindAll(context.TODO(), 10, 0, make(map[string]interface{}))

		assert := assert.New(t)
		assert.NoError(err)
		assert.Len(list, len(mockListUser))

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-happens-in-db", func(t *testing.T) {
		errRepository := domain.NewErrRepository(errors.New("Unexpected error"))

		mockUserRepo.On("FindAll",
			mock.Anything,
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.Anything,
		).Return(mockListUser, errRepository).Once()

		userUse := NewUserUsecase(mockUserRepo, time.Second*2)
		_, err := userUse.FindAll(context.TODO(), 10, 0, make(map[string]interface{}))

		assert := assert.New(t)
		assert.Error(err)
		assert.Equal(err, errRepository)

		mockUserRepo.AssertExpectations(t)
	})
}
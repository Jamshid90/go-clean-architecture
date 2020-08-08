package user

import (
	"testing"
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
	"time"
)

func TestUser(t *testing.T) *domain.User {
	t.Helper()
	return &domain.User{
		ID : "123456789",
		Email :"user@inifo.com",
		Phone :"000000000000",
		Status : "active",
		FirstName :"User",
		LastName:"Qwerty",
		BirthDate: time.Now(),
	}
}

func TestUserEmpty(t *testing.T) *domain.User {
	t.Helper()
	return &domain.User{}
}



package user

import (
	"testing"
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
)

func TestUser(t *testing.T) *domain.User {
	t.Helper()
	return &domain.User{
		ID : 1,
		Email :"user@inifo.com",
		Status : "active",
		FirstName :"User",
		LastName:"Qwerty",
	}
}


func TestUserEmpty(t *testing.T) *domain.User {
	t.Helper()
	return &domain.User{}
}



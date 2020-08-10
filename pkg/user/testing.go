package user

import (
	"testing"
	"github.com/Jamshid90/go-clean-architecture/pkg/entity"
	"time"
)

func TestUser(t *testing.T) *entity.User {
	t.Helper()
	return &entity.User{
		ID : "123456789",
		Email :"user@inifo.com",
		Phone :"000000000000",
		Status : "active",
		FirstName :"User",
		LastName:"Qwerty",
		BirthDate: time.Now(),
	}
}

func TestUserEmpty(t *testing.T) *entity.User {
	t.Helper()
	return &entity.User{}
}



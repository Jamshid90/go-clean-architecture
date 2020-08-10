package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSanitize(t *testing.T) {
	user := User{
		Password: "123456789",
	}
	assert.Empty(t, user.Sanitize().Password)
}

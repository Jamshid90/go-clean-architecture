package validation

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
)

func TestValidator(t *testing.T)  {
	t.Run("success", func(t *testing.T) {
		err := Validator(TestValidationData(t))
		assert.NoError(t, err)
	})
	t.Run("error", func(t *testing.T) {
		data := TestValidationData(t)
		data.Email     = ""

		err := Validator(data)
		assert.Error(t, err)

		errValidation := err.(*errors.ErrValidation)
		assert.NotEmpty(t, errValidation.Errors["email"])
	})
}
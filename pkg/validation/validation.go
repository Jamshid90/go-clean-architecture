package validation

import (
	"errors"
	"reflect"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"
	ut "github.com/go-playground/universal-translator"
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
	field_tag "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate = validator.New()
	translator = en.New()
	uni = ut.New(translator, translator)
)

func Validator(s interface{}) error  {
	trans, found := uni.GetTranslator("en")
	if !found {
		return errors.New("Validator translator not found")
	}
	field_tag.RegisterDefaultTranslations(validate, trans)
	err := validate.Struct(s)

	if  err != nil {
		errValidation     := domain.NewErrValidation()
		errValidation.Err  = err
		for _, fieldError := range  err.(validator.ValidationErrors) {
			field_tag := fieldError.Field()
			field, ok := reflect.TypeOf(s).Elem().FieldByName(field_tag)
			if ok {
				field_tag = field.Tag.Get("json")
			}
			errValidation.Errors[field_tag] = fieldError.Translate(trans)
		}
		return errValidation
	}
	return nil
}

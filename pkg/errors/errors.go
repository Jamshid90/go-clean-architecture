package errors

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrUnauthorized           = errors.New(GetHTTPStatusText(http.StatusUnauthorized))
	ErrUnprocessableEntity    = errors.New(GetHTTPStatusText(http.StatusUnprocessableEntity))
	ErrInternalServerError    = errors.New(GetHTTPStatusText(http.StatusInternalServerError))
	BadRequest                = errors.New(GetHTTPStatusText(http.StatusBadRequest))
	ErrBadParamInput          = errors.New("Given param is not valid")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
)

// Get http status text
func GetHTTPStatusText(statusCode int) string {
	return strings.ToLower(http.StatusText(statusCode))
}

// Error Not Found
func NewErrNotFound(text string) error {
	return &ErrNotFound{text}
}

type ErrNotFound struct {
	name string
}

func (e *ErrNotFound) Error() string {
	return e.name + " not found"
}

// Error Conflict
func NewErrConflict(text string) error {
	return &ErrConflict{text}
}

type ErrConflict struct {
	name string
}

func (e *ErrConflict) Error() string {
	return e.name + " already exist"
}

// Error Repository
func NewErrRepository(err error) error {
	return &ErrRepository{Err: err}
}

type ErrRepository struct {
	Err error
}

func (e ErrRepository) Error() string {
	return e.Err.Error()
}

// Error validation
func NewErrValidation() *ErrValidation {
	return &ErrValidation{Errors: make(map[string]interface{})}
}

type ErrValidation struct {
	Err    error
	Errors map[string]interface{}
}

func (e ErrValidation) Error() string {
	return ErrUnprocessableEntity.Error()
}

// Error bad request
type ErrBadRequest struct {
	Err     error
	Message string
}

func (e ErrBadRequest) Error() string {
	return e.Message
}

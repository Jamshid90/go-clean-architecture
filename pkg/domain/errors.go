package domain

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrInternalServerError = errors.New(GetHTTPStatusText(http.StatusInternalServerError))
	BadRequest             = errors.New(GetHTTPStatusText(http.StatusBadRequest))
	ErrBadParamInput       = errors.New("Given Param is not valid")
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

// Error repository
type ErrRepository struct {
	Err error
}

func (e ErrRepository) Error() string  {
	return e.Err.Error()
}

// Error validation
func NewErrValidation() *ErrValidation {
	return &ErrValidation{Errors: make(map[string]interface{})}
}

type ErrValidation struct {
	Err    error
	Text   string
	Errors map[string]interface{}
}

func (e ErrValidation) Error() string {
	return e.Err.Error()
}
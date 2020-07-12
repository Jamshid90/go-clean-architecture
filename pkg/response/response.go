package response

import (
	"encoding/json"
	"errors"
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
	"github.com/Jamshid90/go-clean-architecture/pkg/request"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var (
	errValidation    *domain.ErrValidation
	errConflict      *domain.ErrConflict
	errNotFound      *domain.ErrNotFound
	errRepository    domain.ErrRepository
	errBadRequest    *request.ErrBadRequest
	validationErrors validator.ValidationErrors
)

type ResponseJSON struct {
	Status     string                 `json:"status"`
	Error      string                 `json:"error,omitempty"`
	Errors     map[string]interface{} `json:"errors,omitempty"`
	Data       interface{}            `json:"data,omitempty"`
}

func GetStatusCodeErr(err error) int {
	if err == nil {
		return http.StatusOK
	}

	//fmt.Println("GetStatusCodeErr", reflect.TypeOf(err))

	switch {
	// Error Bad Request
	case errors.As(err, &errBadRequest):
		return http.StatusBadRequest

	// Error Conflict
	case errors.As(err, &errConflict):
		return http.StatusBadRequest

	// Error Validation Errors
	case errors.As(err, &validationErrors):
		return http.StatusUnprocessableEntity

	// Error Not Found
	case errors.As(err, &errNotFound):
		return http.StatusNotFound

	// Error Validation
	case errors.As(err, &errValidation):
		return http.StatusUnprocessableEntity

	// Error Repository
	case errors.As(err, &errRepository):
		return http.StatusInternalServerError

	default:
		return http.StatusInternalServerError
	}

	return 0
}

func Json(w http.ResponseWriter, r *http.Request, status int, data interface{})  {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, r *http.Request, err error, status int)  {
	w.WriteHeader(status)
	response := ResponseJSON{
		Status:"error",
		Error: err.Error(),
	}

	switch {
		case errors.As(err, &errValidation):
			response = ResponseJSON{
				Status:"errors",
				Error:  errValidation.Error(),
				Errors: errValidation.Errors,
			}
	}

	json.NewEncoder(w).Encode(response)
}
package request

import (
	"encoding/json"
	"errors"
	"fmt"
	apperrors "github.com/Jamshid90/go-clean-architecture/pkg/errors"
	"io"
	"net/http"
)

var (
	syntaxError *json.SyntaxError
)

func DecodeJson(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &apperrors.ErrBadRequest{Err: err, Message: msg}
		case errors.Is(err, io.EOF):
			msg := "request body must not be empty"
			return &apperrors.ErrBadRequest{Err: err, Message: msg}
		default:
			return err
		}
	}
	return nil
}

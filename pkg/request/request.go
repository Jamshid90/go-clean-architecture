package request

import (
	"fmt"
	"io"
	"errors"
	"net/http"
	"encoding/json"
)

var (
	syntaxError *json.SyntaxError
)

type ErrBadRequest struct {
	Err     error
	Message string
}

func (e ErrBadRequest) Error() string  {
	return e.Message
}

func DecodeJson(r *http.Request, v interface{}) error{
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &ErrBadRequest{Err: err, Message: msg}
		case errors.Is(err, io.EOF):
			msg := "request body must not be empty"
			return &ErrBadRequest{Err: err, Message: msg}
		default:
			return err
		}
	}
	return nil
}
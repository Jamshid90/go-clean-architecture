package middleware

import (
	"context"
	"github.com/Jamshid90/go-clean-architecture/pkg/errors"
	"github.com/Jamshid90/go-clean-architecture/pkg/http/rest/response"
	"github.com/Jamshid90/go-clean-architecture/pkg/token"
	"net/http"
)

func Auth(jwtsecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := token.GetAuthUser(jwtsecret, r)
			if err != nil {
				response.Error(w, r, errors.ErrUnauthorized, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", user)))
		})
	}
}

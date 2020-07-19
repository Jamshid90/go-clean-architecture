package auth

import (
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
	"github.com/Jamshid90/go-clean-architecture/pkg/hash"
	"github.com/Jamshid90/go-clean-architecture/pkg/request"
	"github.com/Jamshid90/go-clean-architecture/pkg/response"
	"github.com/Jamshid90/go-clean-architecture/pkg/validation"
	"github.com/go-chi/chi"
	"net/http"
)

type AuthHandler struct {
	UserUsecase domain.UserUsecase
}

// New user handler
func NewAuthHandler(r chi.Router, userUsecase domain.UserUsecase)  {
	handler := AuthHandler{ UserUsecase:userUsecase }

	r.Post("/auth/login", handler.login())
	r.Post("/auth/signup", handler.signup())
}

// login
func (a *AuthHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var loginRequest LoginRequest
		if err := request.DecodeJson(r, &loginRequest); err != nil {
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		if err := validation.Validator(&loginRequest); err != nil {
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		ctx  := r.Context()
		user, err := a.UserUsecase.FindByEmail(ctx, loginRequest.Email)
		if err != nil {
			response.Error(w, r, domain.ErrInvalidEmailOrPassword, http.StatusUnauthorized)
			return
		}

		if hash.CheckPasswordHash(loginRequest.Password, user.Password) == false {
			response.Error(w, r, domain.ErrInvalidEmailOrPassword, http.StatusUnauthorized)
			return
		}

		response.Json(w, r, 200, map[string]interface{}{
			"status" : "success",
			"data" : user,
		})
	}
}

// signup
func (a *AuthHandler) signup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var signupRequest SignupRequest
		if err := request.DecodeJson(r, &signupRequest); err != nil {
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		if err := validation.Validator(&signupRequest); err != nil {
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		ctx  := r.Context()
		user := domain.User{
			Status    : domain.USER_STATUS_ACTIVE,
			Email     : signupRequest.Email,
			FirstName : signupRequest.FirstName,
			LastName  : signupRequest.LastName,
			Password  : signupRequest.Password,
		}

		if err := a.UserUsecase.Store(ctx, &user); err != nil {
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		response.Json(w, r, 200, map[string]interface{}{
			"status" : "success",
			"data" : signupRequest,
		})
	}
}

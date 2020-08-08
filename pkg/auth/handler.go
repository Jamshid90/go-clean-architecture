package auth

import (
	"github.com/Jamshid90/go-clean-architecture/pkg/config"
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
	"github.com/Jamshid90/go-clean-architecture/pkg/errors"
	"github.com/Jamshid90/go-clean-architecture/pkg/hash"
	"github.com/Jamshid90/go-clean-architecture/pkg/http/rest/request"
	"github.com/Jamshid90/go-clean-architecture/pkg/http/rest/response"
	"github.com/Jamshid90/go-clean-architecture/pkg/http/rest/middleware"
	"github.com/Jamshid90/go-clean-architecture/pkg/token"
	"github.com/Jamshid90/go-clean-architecture/pkg/validation"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AuthHandler struct {
	logger *zap.Logger
	config *config.Config
	userUsecase domain.UserUsecase
	refreshTokenUsecase domain.RefreshTokenUsecase
}

// New user handler
func NewAuthHandler(r chi.Router, userUsecase domain.UserUsecase, refreshTokenUsecase domain.RefreshTokenUsecase, config *config.Config, logger *zap.Logger)  {
	handler := AuthHandler{
		logger: logger,
		config: config,
		userUsecase: userUsecase,
		refreshTokenUsecase: refreshTokenUsecase,
	}

	r.Post("/auth/login", handler.login())
	r.Post("/auth/signup", handler.signup())
	r.Post("/auth/refresh-token", handler.refreshToken())

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(config.Jwt.Secret))
		r.Get("/auth/logout", handler.logout())
	})
}

// login
func (a *AuthHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var loginRequest LoginRequest
		if err := request.DecodeJson(r, &loginRequest); err != nil {
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		// validation login request
		if err := validation.Validator(&loginRequest); err != nil {
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		ctx  := r.Context()
		// find user by email
		user, err := a.userUsecase.FindByEmail(ctx, loginRequest.Email)
		if err != nil {
			a.logger.Error("auth login find by email", zap.Error(err))
			response.Error(w, r, errors.ErrInvalidEmailOrPassword, http.StatusUnauthorized)
			return
		}

		// check password
		if hash.CheckPasswordHash(loginRequest.Password, user.Password) == false {
			response.Error(w, r, errors.ErrInvalidEmailOrPassword, http.StatusUnauthorized)
			return
		}

		// generate token
		access_token, refresh_token, err := token.GenerateToken(a.config.Jwt.Secret, a.config.Jwt.AccessTTL, a.config.Jwt.RefreshTTL, user.ID)
		if err != nil {
			a.logger.Error("auth login generate token", zap.Error(err))
			response.Error(w, r, errors.ErrInternalServerError, http.StatusInternalServerError)
			return
		}

		// create refresh token
		if err = a.refreshTokenUsecase.Store(ctx, &domain.RefreshToken{
			UserID: user.ID,
			Token: refresh_token,
		}); err != nil {
			a.logger.Error("auth login refresh token store", zap.Error(err))
			response.Error(w, r, errors.ErrInternalServerError, response.GetStatusCodeErr(err))
			return
		}

		userInfo := User{
			ID: user.ID,
			Email: user.Email,
			Phone: user.Phone,
			Gender: user.Gender,
			FirstName: user.FirstName,
			LastName: user.LastName,
			BirthDate: user.BirthDate,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		response.Json(w, r, 200, map[string]interface{}{
			"status" : "success",
			"data" : userInfo,
			"token" : map[string]string{
				"type" : "Bearer",
				"access" : access_token,
				"refresh" : refresh_token,
			},
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

		birthDate, err := time.Parse("2006-01-02", signupRequest.BirthDate)
		if err != nil {
			a.logger.Error("auth signup parse birth date", zap.Error(err))
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		ctx  := r.Context()
		user := domain.User{
			Status    : domain.USER_STATUS_ACTIVE,
			Email     : signupRequest.Email,
			Phone     : signupRequest.Phone,
			Gender    : signupRequest.Gender,
			BirthDate : birthDate,
			FirstName : signupRequest.FirstName,
			LastName  : signupRequest.LastName,
			Password  : signupRequest.Password,
		}

		if err := a.userUsecase.Store(ctx, &user); err != nil {
			a.logger.Error("auth signup user store", zap.Error(err))
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		userInfo := User{
			ID: user.ID,
			Email: user.Email,
			Phone: user.Phone,
			Gender: user.Gender,
			BirthDate: user.BirthDate,
			FirstName: user.FirstName,
			LastName: user.LastName,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		response.Json(w, r, 200, map[string]interface{}{
			"status" : "success",
			"data" : userInfo,
		})
	}
}

// logout
func (a *AuthHandler) logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(*domain.User)
		if !ok {
			response.Error(w, r, errors.ErrInternalServerError, http.StatusInternalServerError)
			return
		}

		ctx  := r.Context()
		if err := a.refreshTokenUsecase.DeleteByUserId(ctx, user.ID); err != nil {
			a.logger.Error("auth logout delete refresh token", zap.Error(err))
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		response.Json(w, r, 200, map[string]interface{}{
			"status" : "success",
		})
	}
}

// refresh token
func (a *AuthHandler) refreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var refreshTokenRequest RefreshTokenRequest
		if err := request.DecodeJson(r, &refreshTokenRequest); err != nil {
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		// validation login request
		if err := validation.Validator(&refreshTokenRequest); err != nil {
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		ctx  := r.Context()
		refreshToken, err := a.refreshTokenUsecase.Find(ctx, refreshTokenRequest.Token)
		if err != nil {
			a.logger.Error("auth refresh token find", zap.Error(err))
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		if _, err := token.ParseJwtToken(refreshToken.Token, a.config.Jwt.Secret); err != nil {
			if err := a.refreshTokenUsecase.Delete(ctx, refreshToken.Token); err != nil {
				a.logger.Error("auth refresh token delete", zap.Error(err))
				response.Error(w, r, err, response.GetStatusCodeErr(err))
				return
			}
			response.Error(w, r, err, 400)
			return
		}

		// generate token
		access_token, refresh_token, err := token.GenerateToken(a.config.Jwt.Secret, a.config.Jwt.AccessTTL, a.config.Jwt.RefreshTTL, refreshToken.UserID)
		if err != nil {
			a.logger.Error("auth refresh token generate", zap.Error(err))
			response.Error(w, r, errors.ErrInternalServerError, http.StatusInternalServerError)
			return
		}

		// create refresh token
		if err = a.refreshTokenUsecase.Store(ctx, &domain.RefreshToken{
			UserID: refreshToken.UserID,
			Token: refresh_token,
		}); err != nil {
			a.logger.Error("auth refresh token store", zap.Error(err))
			response.Error(w, r, errors.ErrInternalServerError, response.GetStatusCodeErr(err))
			return
		}

		if err := a.refreshTokenUsecase.Delete(ctx, refreshToken.Token); err != nil {
			a.logger.Error("auth refresh token delete old token", zap.Error(err))
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		response.Json(w, r, 200, map[string]interface{}{
			"status" : "success",
			"token" : map[string]string{
				"type" : "Bearer",
				"access" : access_token,
				"refresh" : refresh_token,
			},
		})
	}
}
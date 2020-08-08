package user

import (
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
	"github.com/Jamshid90/go-clean-architecture/pkg/http/rest/request"
	"github.com/Jamshid90/go-clean-architecture/pkg/http/rest/response"
	"github.com/Jamshid90/go-clean-architecture/pkg/validation"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type UserHandler struct {
	logger *zap.Logger
	userUsecase domain.UserUsecase
}

// New user handler
func NewUserHandler(r chi.Router, userUsecase domain.UserUsecase, logger *zap.Logger)  {
	handler := UserHandler{
		userUsecase: userUsecase,
		logger: logger,
	}
	r.Get("/user", handler.findAll())
	r.Get("/user/{id}", handler.find())
	r.Post("/user", handler.store())
	r.Put("/user", handler.update())
	r.Delete("/user/{id}", handler.delete())
}

// convert domain user to user model
func (uh *UserHandler) convert(user *domain.User) *User  {
	return &User{
		ID        : user.ID,
		Status    : user.Status,
		Email     : user.Email,
		FirstName : user.FirstName,
		LastName  : user.LastName,
		Password  : user.Password,
		CreatedAt : user.CreatedAt,
		UpdatedAt : user.UpdatedAt,
	}
}

// convert items
func (uh *UserHandler) convertItems(items []*domain.User) []*User  {
	var users []*User
	for _, item := range items  {
		users = append(users, uh.convert(item))
	}
	return users
}

// store
func (uh *UserHandler) store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userRequest CreateUserRequest
		if err := request.DecodeJson(r, &userRequest); err != nil {
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		if err := validation.Validator(&userRequest); err != nil {
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		birthDate, err := time.Parse("2006-01-02", userRequest.BirthDate)
		if err != nil {
			uh.logger.Error("user store parse birth date", zap.Error(err))
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		ctx  := r.Context()
		user := domain.User{
			Status    : userRequest.Status,
			Email     : userRequest.Email,
			Phone     : userRequest.Phone,
			Gender    : userRequest.Gender,
			FirstName : userRequest.FirstName,
			LastName  : userRequest.LastName,
			Password  : userRequest.Password,
			BirthDate : birthDate,
		}
		if err := uh.userUsecase.Store(ctx, &user); err != nil {
			uh.logger.Error("user store", zap.Error(err))
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		response.Json(w, r, 200, map[string]interface{}{
			"status" : "success",
			"data" : uh.convert(&user).Sanitize(),
		})
	}
}

// update
func (uh *UserHandler) update() http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {

		var userRequest UpdateUserRequest
		if err := request.DecodeJson(r, &userRequest); err != nil {
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		if err := validation.Validator(&userRequest); err != nil {
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		birthDate, err := time.Parse("2006-01-02", userRequest.BirthDate)
		if err != nil {
			uh.logger.Error("user update parse birth date", zap.Error(err))
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		ctx  := r.Context()
		user := domain.User{
			ID        : userRequest.ID,
			Status    : userRequest.Status,
			Email     : userRequest.Email,
			Phone     : userRequest.Phone,
			Gender    : userRequest.Gender,
			FirstName : userRequest.FirstName,
			LastName  : userRequest.LastName,
			BirthDate : birthDate,
		}

		if err := uh.userUsecase.Update(ctx, &user); err != nil {
			uh.logger.Error("user update", zap.Error(err))
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		response.Json(w, r, 200, map[string]interface{}{
			"status" : "success",
			"user" : uh.convert(&user),
		})
	}
}

// delete
func (uh *UserHandler) delete() http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx  := r.Context()
		if err := uh.userUsecase.Delete(ctx, chi.URLParam(r, "id")); err != nil {
			uh.logger.Error("user delete", zap.Error(err))
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		response.Json(w, r, 200, map[string]interface{}{
			"status" : "success",
		})
	}
}

// find
func (uh *UserHandler) find() http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx  := r.Context()
		user, err := uh.userUsecase.Find(ctx, chi.URLParam(r, "id"))
		if  err != nil {
			uh.logger.Error("user find", zap.Error(err))
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		response.Json(w, r, 200, map[string]interface{}{
			"status" : "success",
			"data" : uh.convert(user).Sanitize(),
		})
	}
}

// find all
func (uh *UserHandler) findAll() http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			limit  = 10
			offset = 0
			params = make(map[string]interface{})
		)

		if _limit, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
			uh.logger.Error("user find all limit", zap.Error(err))
			limit = _limit
		}

		if _offset, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil {
			uh.logger.Error("user find all offset", zap.Error(err))
			offset = _offset
		}

		for i, v := range r.URL.Query() {
			params[i] = v
		}

		ctx  := r.Context()
		items, err := uh.userUsecase.FindAll(ctx, limit, offset, params)
		if  err != nil {
			uh.logger.Error("user find all", zap.Error(err))
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		response.Json(w, r, 200, map[string]interface{}{
			"status" : "success",
			"items" : uh.convertItems(items),
		})
	}
}
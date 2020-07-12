package user

import (
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
	"github.com/Jamshid90/go-clean-architecture/pkg/request"
	"github.com/Jamshid90/go-clean-architecture/pkg/response"
	"github.com/Jamshid90/go-clean-architecture/pkg/validation"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
}

// New user handler
func NewUserHandler(r chi.Router, userUsecase domain.UserUsecase)  {
	handler := UserHandler{ UserUsecase:userUsecase }
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

		user := domain.User{
			Status    : userRequest.Status,
			Email     : userRequest.Email,
			FirstName : userRequest.FirstName,
			LastName  : userRequest.LastName,
			Password  : userRequest.Password,
		}
		if err := uh.UserUsecase.Store(&user); err != nil {
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

		user := domain.User{
			ID        : userRequest.ID,
			Status    : userRequest.Status,
			Email     : userRequest.Email,
			FirstName : userRequest.FirstName,
			LastName  : userRequest.LastName,
		}

		if err := uh.UserUsecase.Update(&user); err != nil {
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
		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
		if err != nil {
			err := domain.NewErrNotFound("user")
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		if err := uh.UserUsecase.Delete(id); err != nil {
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
		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
		if err != nil {
			err := domain.NewErrNotFound("user")
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		user, err := uh.UserUsecase.Find(id)
		if  err != nil {
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
			limit = _limit
		}

		if _offset, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil {
			offset = _offset
		}

		for i, v := range r.URL.Query() {
			params[i] = v
		}

		items, err := uh.UserUsecase.FindAll(limit, offset, params)
		if  err != nil {
			response.Error(w, r, err, response.GetStatusCodeErr(err))
			return
		}

		response.Json(w, r, 200, map[string]interface{}{
			"status" : "success",
			"items" : uh.convertItems(items),
		})
	}
}
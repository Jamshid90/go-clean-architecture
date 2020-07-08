package http

import (
	"encoding/json"
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

type ResourceHandler struct {
	RUsecase domain.ResourceUsecase
}

func NewresourceHandler(r chi.Router, resUsecase domain.ResourceUsecase)  {

	handler := ResourceHandler{ RUsecase:resUsecase }

	r.Post("/resource", handler.Store())
	r.Get("/resource", handler.GetList())
	r.Get("/resource/{id}", handler.GetList())
}

func (rh *ResourceHandler) Store() http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {

		var resource domain.Resource

		if err := json.NewDecoder(r.Body).Decode(&resource); err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"error":err.Error(),
			})
		}

		if err := rh.RUsecase.Store(&resource); err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"error":err.Error(),
			})
		}

		json.NewEncoder(w).Encode(map[string]bool{
			"success" : true,
		})
	}
}

func (rh *ResourceHandler) GetList() http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {

		limit  := 10
		offset := 10

		list, err := rh.RUsecase.FindAll(limit, offset, nil)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"error":err.Error(),
			})
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"list" : list,
		})
	}
}

func (rh *ResourceHandler) GetById() http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"error":err.Error(),
			})
		}

		resource, err := rh.RUsecase.Find(id)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"error":err.Error(),
			})
		}

		json.NewEncoder(w).Encode(resource)
	}
}
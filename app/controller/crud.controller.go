package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	mw "de.server/app/middleware"
	m "de.server/app/model"
	s "de.server/app/store"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Controller[T m.Model] struct {
	model  string
	router *chi.Mux
	Store  *s.CrudStore[T]
}

func NewController[T m.Model](database *gorm.DB, router *chi.Mux, model string) *Controller[T] {
	return &Controller[T]{
		model:  model,
		router: router,
		Store:  &s.CrudStore[T]{DB: database, TableName: model},
	}
}

// CRUD Create godoc
//
//	@Summary		CRUD Create
//	@Accept			json
//	@Produce		json
//	@Param			modelName	path	string	true	"Model Name"
//	@Success		200	{object}	interface{}
//	@Failure		400	{object}	interface{}
//	@Failure		404	{object}	interface{}
//	@Failure		500	{object}	interface{}
//	@Router			/{modelName} [put]
func (controller *Controller[T]) Create(path string, handlers ...mw.Middleware) {
	fmt.Println("PUT", path)
	controller.router.
		With(mw.GetMiddlewares(handlers)...).
		Put(path, func(w http.ResponseWriter, r *http.Request) {
			response, err := m.CreateFromBody[T](r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			response, err = controller.Store.AddOne(response)
			HandleResponse(w, response, err)
		})
}

// CRUD FindOne godoc
//
//	@Summary		CRUD FindOne
//	@Accept			json
//	@Produce		json
//	@Param			modelName	path	string	true	"Model Name"
//	@Success		200	{object}	interface{}
//	@Failure		400	{object}	interface{}
//	@Failure		404	{object}	interface{}
//	@Failure		500	{object}	interface{}
//	@Router			/{modelName} [get]
func (controller *Controller[T]) FindOne(path string, preload *[]string, handlers ...mw.Middleware) {
	fmt.Println("GET", path)
	controller.router.
		With(mw.GetMiddlewares(handlers)...).
		Get(path, func(w http.ResponseWriter, r *http.Request) {
			data, err := m.CreateFromBody[T](r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response, err := controller.Store.FindOne(data, preload)
			HandleResponse(w, response, err)
		})
}

// CRUD FindMany godoc
//
//	@Summary		CRUD FindMany
//	@Accept			json
//	@Produce		json
//	@Param			modelName	path	string	true	"Model Name"
//	@Success		200	{object}	interface{}
//	@Failure		400	{object}	interface{}
//	@Failure		404	{object}	interface{}
//	@Failure		500	{object}	interface{}
//	@Router			/{modelName}/list [get]
func (controller *Controller[T]) FindMany(path string, preload *[]string, filters *[]string, handlers ...mw.Middleware) {
	fmt.Println("GET", path)
	controller.router.
		With(mw.GetMiddlewares(handlers)...).
		Get(path, func(w http.ResponseWriter, r *http.Request) {
			data, err := m.CreateFromBody[T](r)
			var response *[]T
			// no body in request
			if data == nil && err == nil {
				response, err = controller.Store.Filter(r, preload, filters)
			} else {
				response, err = controller.Store.FindMany(r, data, preload)
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			HandleResponse(w, response, err)
		})
}

// CRUD Update godoc
//
//	@Summary		CRUD Update
//	@Accept			json
//	@Produce		json
//	@Param			modelName	path	string	true	"Model Name"
//	@Success		200	{object}	interface{}
//	@Failure		400	{object}	interface{}
//	@Failure		404	{object}	interface{}
//	@Failure		500	{object}	interface{}
//	@Router			/{modelName} [patch]
func (controller *Controller[T]) Update(path string, handlers ...mw.Middleware) {
	fmt.Println("PATCH", path)
	controller.router.
		With(mw.GetMiddlewares(handlers)...).
		Patch(path, func(w http.ResponseWriter, r *http.Request) {
			data, err := m.CreateFromBody[T](r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response, err := controller.Store.Update(data)
			HandleResponse(w, response, err)
		})
}

// CRUD Delete godoc
//
//	@Summary		CRUD Delete
//	@Accept			json
//	@Produce		json
//	@Param			modelName	path	string	true	"Model Name"
//	@Success		200	{object}	interface{}
//	@Failure		400	{object}	interface{}
//	@Failure		404	{object}	interface{}
//	@Failure		500	{object}	interface{}
//	@Router			/{modelName} [delete]
func (controller *Controller[T]) Delete(path string, preload *[]string, handlers ...mw.Middleware) {
	fmt.Println("DELETE", path)
	controller.router.
		With(mw.GetMiddlewares(handlers)...).
		Delete(path, func(w http.ResponseWriter, r *http.Request) {
			data, err := m.CreateFromBody[T](r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response, err := controller.Store.FindOne(data, preload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response, err = controller.Store.Delete(response)
			HandleResponse(w, response, err)
		})
}

func HandleResponse(w http.ResponseWriter, response interface{}, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponse(w, &response)
}

func WriteResponse(w http.ResponseWriter, response *interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

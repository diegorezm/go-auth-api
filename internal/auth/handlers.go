package auth

import (
	"log"
	"net/http"

	"github.com/diegorezm/ticketing/internal/json"
	"github.com/diegorezm/ticketing/internal/responses"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var params loginParams

	if err := json.Read(r, &params); err != nil {
		log.Println(err)
		responses.Fail(w, http.StatusBadRequest, err.Error())
		return
	}

	responses.Ok(w, http.StatusOK, params)
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var params registerParams

	if err := json.Read(r, &params); err != nil {
		log.Println(err)
		responses.Fail(w, http.StatusBadRequest, err.Error())
		return
	}
	responses.Ok(w, http.StatusCreated, params)
}

func (h *handler) Mount(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/register", h.Register)
	})
}

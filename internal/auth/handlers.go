package auth

import (
	"log"
	"net/http"

	"github.com/diegorezm/ticketing/internal/json"
	"github.com/diegorezm/ticketing/internal/jwt"
	"github.com/diegorezm/ticketing/internal/middlewares"
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

	result, err := h.service.Login(r.Context(), params)

	if err != nil {
		handleAuthError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    result.Token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // true in prod
		MaxAge:   60 * 60 * 24,
	})

	responses.Ok(w, http.StatusOK, result)
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var params registerParams

	if err := json.Read(r, &params); err != nil {
		log.Println(err)
		responses.Fail(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Register(r.Context(), params); err != nil {
		handleAuthError(w, err)
		return
	}

	responses.OkMessage(w, http.StatusCreated, "User was created successfully!")
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // true in prod
		MaxAge:   -1,
	})

	responses.OkMessage(w, http.StatusOK, "logged out")
}

func (h *handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(jwt.UserIDKey).(string)
	if !ok || userID == "" {
		responses.Fail(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.service.Me(r.Context(), userID)
	if err != nil {
		handleAuthError(w, err)
		return
	}

	responses.Ok(w, http.StatusOK, user)
}

func (h *handler) Mount(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/register", h.Register)
	})

	r.Group(func(r chi.Router) {
		r.Use(middlewares.JWT)
		r.Post("/logout", h.Logout)
		r.Get("/me", h.Me)
	})
}

func handleAuthError(w http.ResponseWriter, err error) {
	switch err {
	case ErrInvalidCredentials:
		responses.Fail(w, http.StatusUnauthorized, err.Error())

	case ErrEmailAlreadyExists:
		responses.Fail(w, http.StatusConflict, err.Error())

	default:
		responses.Fail(w, http.StatusInternalServerError, "internal server error")
	}
}

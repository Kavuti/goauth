package users

import (
	"encoding/json"
	"net/http"

	"github.com/Kavuti/goauth/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type UserHandler interface {
	Routes() chi.Router
	Registration(w http.ResponseWriter, r *http.Request)
	Verify(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	chi.Router
	service UserService
}

func (h *userHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/registration", h.Registration)
	r.Route("/{name}", func(r chi.Router) {
		r.Post("/verify", h.Verify)
	})

	return r
}

func NewUserHandler(r chi.Router, db *sqlx.DB) UserHandler {
	handler := &userHandler{
		Router:  r,
		service: NewUserService(db),
	}

	return handler
}

func (h *userHandler) Registration(w http.ResponseWriter, r *http.Request) {
	defer utils.RecoverIfError(w, r)
	request := RegistrationRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	utils.CheckError(err)

	err = utils.ValidateStruct(request)
	utils.CheckError(err)

	err = h.service.Registration(request.FirstName, request.LastName, request.Email, request.Password)
	utils.CheckError(err)
}

func (h *userHandler) Verify(w http.ResponseWriter, r *http.Request) {
	defer utils.RecoverIfError(w, r)
	email := chi.URLParam(r, "email")
	err := h.service.Verify(email)
	utils.CheckError(err)
}

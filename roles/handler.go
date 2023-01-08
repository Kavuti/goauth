package roles

import (
	"encoding/json"
	"net/http"

	"github.com/Kavuti/goauth/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

type RoleHandler interface {
	Routes() chi.Router

	SearchByVisibleName(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type rolesHandler struct {
	chi.Router
	service RoleService
}

func (h *rolesHandler) SearchByVisibleName(w http.ResponseWriter, r *http.Request) {
	defer utils.RecoverIfError(w, r)
	visibleName := r.URL.Query().Get("visibleName")
	roles, err := h.service.SearchByVisibleName(visibleName)
	utils.CheckError(err)

	render.Render(w, r, &MultipleRoleResponse{Roles: roles})
}

func (h *rolesHandler) Get(w http.ResponseWriter, r *http.Request) {
	defer utils.RecoverIfError(w, r)
	name := chi.URLParam(r, "name")
	role, err := h.service.Get(name)
	utils.CheckError(err)

	render.Render(w, r, &SingleRoleResponse{Role: *role})
}

func (h *rolesHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer utils.RecoverIfError(w, r)
	request := RoleCreationRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	utils.CheckError(err)

	err = utils.ValidateStruct(request)
	utils.CheckError(err)

	err = h.service.Create(&request)
	utils.CheckError(err)
}

func (h *rolesHandler) Update(w http.ResponseWriter, r *http.Request) {
	defer utils.RecoverIfError(w, r)
	name := chi.URLParam(r, "name")
	request := RoleUpdateRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	utils.CheckError(err)

	err = utils.ValidateStruct(request)
	utils.CheckError(err)

	err = h.service.Update(name, &request)
	utils.CheckError(err)

}

func (h *rolesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	defer utils.RecoverIfError(w, r)
	name := chi.URLParam(r, "name")
	err := h.service.Delete(name)
	utils.CheckError(err)
}

func (h *rolesHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.SearchByVisibleName)
	r.Post("/", h.Create)

	r.Route("/{name}", func(r chi.Router) {
		r.Get("/", h.Get)
		r.Put("/", h.Update)
		r.Delete("/", h.Delete)
	})

	return r
}

func NewRoleHandler(r chi.Router, db *sqlx.DB) RoleHandler {
	handler := &rolesHandler{
		Router:  r,
		service: NewRoleService(db),
	}

	return handler
}

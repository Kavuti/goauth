package roles

import "net/http"

type Role struct {
	Name        string
	VisibleName string `db:"visible_name"`
}

type MultipleRoleResponse struct {
	Roles []Role `json:"roles"`
}

type SingleRoleResponse struct {
	Role Role `json:"role"`
}

type RoleCreationRequest struct {
	Name        string `json:"name" validate:"required,uppercase,max=255"`
	VisibleName string `json:"visibleName" db:"visible_name" validate:"required,max=255"`
}

type RoleUpdateRequest struct {
	VisibleName string `json:"visibleName" db:"visible_name" validate:"required,max=255"`
}

func (resp *MultipleRoleResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (resp *SingleRoleResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

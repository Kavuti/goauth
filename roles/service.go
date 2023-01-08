package roles

import (
	"net/http"

	"github.com/Kavuti/goauth/utils"
	"github.com/jmoiron/sqlx"
)

type RoleService interface {
	SearchByVisibleName(name string) ([]Role, error)
	Get(name string) (*Role, error)
	Create(req *RoleCreationRequest) error
	Update(name string, req *RoleUpdateRequest) error
	Delete(name string) error
}

type roleService struct {
	db *sqlx.DB
}

func (s *roleService) SearchByVisibleName(name string) ([]Role, error) {
	tx := s.db.MustBegin()
	defer tx.Rollback()
	var roles []Role
	query := "SELECT * FROM roles"
	var err error
	if name != "" {
		query = query + " WHERE visible_name like CONCAT('%', $1, '%')"
		err = tx.Select(&roles, query, name)
	} else {
		err = tx.Select(&roles, query)
	}
	if err != nil {
		return nil, utils.ServiceError(err.Error(), http.StatusInternalServerError)
	}
	return roles, nil
}

func (s *roleService) Get(name string) (*Role, error) {
	if name == "" {
		return nil, utils.ServiceError("Name parameter is mandatory", http.StatusBadRequest)
	}

	tx := s.db.MustBegin()
	defer tx.Rollback()

	var role Role
	err := tx.Get(&role, "SELECT * FROM roles WHERE name=$1", name)
	if err != nil {
		return nil, utils.ServiceError(err.Error(), http.StatusBadRequest)
	}
	err = tx.Commit()
	if err != nil {
		return nil, utils.ServiceError(err.Error(), http.StatusInternalServerError)
	}
	return &role, nil
}

func (s *roleService) Create(req *RoleCreationRequest) error {
	err := utils.ValidateStruct(req)
	if err != nil {
		return utils.ServiceError(err.Error(), http.StatusBadRequest)
	}

	tx := s.db.MustBegin()
	defer tx.Rollback()

	var roles []Role
	err = tx.Select(&roles, "SELECT * FROM roles WHERE name=$1", req.Name)
	if err != nil {
		return utils.ServiceError(err.Error(), http.StatusInternalServerError)
	}
	if len(roles) > 0 {
		return utils.ServiceError("Role already existing", http.StatusConflict)
	}

	_, err = tx.NamedExec("INSERT INTO roles (name, visible_name) VALUES (:name, :visible_name)", req)
	if err != nil {
		return utils.ServiceError(err.Error(), http.StatusInternalServerError)
	}
	err = tx.Commit()
	if err != nil {
		return utils.ServiceError(err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (s *roleService) Update(name string, req *RoleUpdateRequest) error {
	if name == "" {
		return utils.ServiceError("Name parameter is mandatory", http.StatusBadRequest)
	}

	err := utils.ValidateStruct(req)
	if err != nil {
		return utils.ServiceError(err.Error(), http.StatusBadRequest)
	}

	tx := s.db.MustBegin()
	defer tx.Rollback()

	rows, err := tx.MustExec("UPDATE roles SET visible_name = $1 WHERE name = $2", req.VisibleName, name).RowsAffected()
	if err != nil {
		return utils.ServiceError(err.Error(), http.StatusInternalServerError)
	}
	if rows == 0 {
		return utils.ServiceError("No role found with the given name", http.StatusNotFound)
	}
	err = tx.Commit()
	if err != nil {
		return utils.ServiceError(err.Error(), http.StatusBadRequest)
	}
	return nil
}

func (s *roleService) Delete(name string) error {
	if name == "" {
		return utils.ServiceError("Name parameter is mandatory", http.StatusBadRequest)
	}

	tx := s.db.MustBegin()
	defer tx.Rollback()

	rows, err := tx.MustExec("DELETE FROM roles WHERE name=$1", name).RowsAffected()
	if err != nil {
		return utils.ServiceError(err.Error(), http.StatusInternalServerError)
	}
	if rows == 0 {
		return utils.ServiceError("No role found with the given name", http.StatusNotFound)
	}
	err = tx.Commit()
	if err != nil {
		return utils.ServiceError(err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func NewRoleService(db *sqlx.DB) RoleService {
	return &roleService{db: db}
}

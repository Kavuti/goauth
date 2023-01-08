package users

import (
	"net/http"

	"github.com/Kavuti/goauth/utils"
	"github.com/jmoiron/sqlx"
)

type UserService interface {
	Registration(firstName string, lastName string, email string, password string) error
	Verify(email string) error
}

type userService struct {
	db *sqlx.DB
}

func NewUserService(db *sqlx.DB) UserService {
	return &userService{db: db}
}

func (s *userService) Registration(firstName string, lastName string, email string, password string) error {
	tx := s.db.MustBegin()
	defer tx.Rollback()

	var existingUsers []User
	err := tx.Select(&existingUsers, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		return utils.ServiceError(err.Error(), http.StatusInternalServerError)
	}
	if len(existingUsers) > 0 {
		return utils.ServiceError("User already exists", http.StatusBadRequest)
	}
	user := NewUserForRegistration(firstName, lastName, email, password)

	_, err = tx.NamedExec(`INSERT INTO users (first_name, last_name, email, password, verified) 
		VALUES (:first_name, :last_name, :email, :password, :verified)`, user)
	if err != nil {
		return utils.ServiceError(err.Error(), http.StatusInternalServerError)
	}
	err = tx.Commit()
	if err != nil {
		return utils.ServiceError(err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (s *userService) Verify(email string) error {
	tx := s.db.MustBegin()
	defer tx.Rollback()

	var user User
	err := tx.Get(&user, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		return utils.ServiceError(err.Error(), http.StatusBadRequest)
	}

	tx.MustExec("UPDATE users SET verified = true WHERE email=$1", email)
	err = tx.Commit()
	if err != nil {
		return utils.ServiceError(err.Error(), http.StatusInternalServerError)
	}
	return nil
}

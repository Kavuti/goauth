package users

import (
	"net/http"
	"os"

	"github.com/Kavuti/goauth/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Verified  bool   `json:"verified"`
}

type RegistrationRequest struct {
	FirstName string `json:"first_name" validate:"required,len=50"`
	LastName  string `json:"last_name" validate:"required,len=50"`
	Email     string `json:"email" validate:"required,emaillen=255"`
	Password  string `json:"password" validate:"required"`
}

func NewUserForRegistration(firstName string, lastName string, email string, password string) *User {
	bytes, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("SECRET_KEY")), bcrypt.DefaultCost)
	if err != nil {
		utils.ServiceError("Error generating a secure password", http.StatusInternalServerError)
	}
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(bytes),
		Verified:  false,
	}
}

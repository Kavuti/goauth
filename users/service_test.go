package users

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func Test_Registration(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").WillReturnRows(sqlmock.NewRows([]string{"first_name", "last_name", "email", "password", "verified"}))
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	service := &userService{db: sqlx.NewDb(db, "sqlmock")}
	err = service.Registration("Test", "Test", "Test", "Test")
	if err != nil {
		t.Fatalf("Error executing Registration test: %s\n", err.Error())
	}
}

func Test_Registration_UserAlreadyExisting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM users WHERE .+").WillReturnRows(sqlmock.NewRows([]string{"first_name", "last_name", "email", "password", "verified"}).AddRow("test", "test", "test", "test", false))
	mock.ExpectRollback()

	service := &userService{db: sqlx.NewDb(db, "sqlmock")}
	err = service.Registration("Test", "Test", "Test", "Test")
	if err == nil {
		t.Fatal("Error executing Registration_UserAlreadyExisting test: no error returned")
	}
}

func Test_Registration_ErrorSelecting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM users WHERE .+").WillReturnError(errors.New("Random error"))
	mock.ExpectRollback()

	service := &userService{db: sqlx.NewDb(db, "sqlmock")}
	err = service.Registration("Test", "Test", "Test", "Test")
	if err == nil {
		t.Fatal("Error executing Registration_ErrorSelecting test: no error returned")
	}
}

func Test_Verify(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM users WHERE .+").WillReturnRows(sqlmock.NewRows([]string{"first_name", "last_name", "email", "password", "verified"}).AddRow("test", "test", "test", "test", false))
	mock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	service := &userService{db: sqlx.NewDb(db, "sqlmock")}
	err = service.Verify("test")
	if err != nil {
		t.Fatalf("Error executing Verify test: %s\n", err.Error())
	}
}

func Test_Verify_MissingUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM users WHERE .+").WillReturnError(errors.New("Not found"))
	mock.ExpectRollback()

	service := &userService{db: sqlx.NewDb(db, "sqlmock")}
	err = service.Verify("test")
	if err == nil {
		t.Fatal("Error executing Verify_MissingUser test: no error returned")
	}

}

package roles

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func Test_SearchByVisibleName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM roles").WillReturnRows(sqlmock.NewRows([]string{"name", "visible_name"}))
	mock.ExpectCommit()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	_, err = service.SearchByVisibleName("")
	if err != nil {
		t.Fatalf("Error executing SearchByVisibleName test: %s\n", err.Error())
	}
}

func Test_SearchByVisibleName_ValuePopulated(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM roles WHERE (.+)").WillReturnRows(sqlmock.NewRows([]string{"name", "visible_name"}).AddRow("ttestt", "ttestt"))
	mock.ExpectCommit()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	_, err = service.SearchByVisibleName("test")
	if err != nil {
		t.Fatalf("Error executing SearchByVisibleName_ValuePopulated test: %s\n", err.Error())
	}
}

func Test_SearchByVisibleName_ErrorSelecting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM roles WHERE (.+)").WillReturnError(errors.New("Random error"))
	mock.ExpectRollback()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	_, err = service.SearchByVisibleName("test")
	if err == nil {
		t.Fatalf("Error executing SearchByVisibleName_ErrorSelecting test: no error returned")
	}
}

func Test_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM roles WHERE (.+)").WillReturnRows(sqlmock.NewRows([]string{"name", "visible_name"}).AddRow("test", "test"))
	mock.ExpectCommit()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	_, err = service.Get("test")
	if err != nil {
		t.Fatalf("Error executing Get test: %s\n", err.Error())
	}
}

func Test_Get_EmptyValue(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	_, err = service.Get("")
	if err == nil {
		t.Fatal("Error executing Get_EmptyValue test: no error returned")
	}
}

func Test_Get_ErrorSelecting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM roles WHERE (.+)").WillReturnError(errors.New("Random error"))
	mock.ExpectRollback()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	_, err = service.Get("test")
	if err == nil {
		t.Fatal("Error executing Get_EmptyValue test: no error returned")
	}
}

func Test_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM roles WHERE (.+)").WillReturnRows(sqlmock.NewRows([]string{"name", "visible_name"}))
	mock.ExpectExec("INSERT INTO roles").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	err = service.Create(&RoleCreationRequest{
		Name:        "TEST",
		VisibleName: "Test",
	})
	if err != nil {
		t.Fatalf("Error executing Create test: %s\n", err.Error())
	}
}

func Test_Create_InvalidPayload(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	err = service.Create(&RoleCreationRequest{
		Name:        "test", // This sould be uppercase
		VisibleName: "Test",
	})

	if err == nil {
		t.Fatal("Error executing Create_InvalidPayload test: no error returned")
	}
}

func Test_Create_AlreadyExisting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM roles WHERE (.+)").WillReturnRows(sqlmock.NewRows([]string{"name", "visible_name"}).AddRow("test", "test"))
	mock.ExpectRollback()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	err = service.Create(&RoleCreationRequest{
		Name:        "TEST",
		VisibleName: "Test",
	})
	if err == nil {
		t.Fatal("Error executing Create_AlreadyExisting test: no error returned")
	}
}

func Test_Create_ErrorSelecting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM roles WHERE (.+)").WillReturnError(errors.New("Random error"))
	mock.ExpectRollback()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	err = service.Create(&RoleCreationRequest{
		Name:        "TEST",
		VisibleName: "Test",
	})
	if err == nil {
		t.Fatal("Error executing Create_ErrorSelecting test: no error returned")
	}
}

func Test_Create_ErrorInserting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM roles WHERE (.+)").WillReturnRows(sqlmock.NewRows([]string{"name", "visible_name"}))
	mock.ExpectExec("INSERT INTO roles").WillReturnError(errors.New("Random error"))
	mock.ExpectRollback()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	err = service.Create(&RoleCreationRequest{
		Name:        "TEST",
		VisibleName: "Test",
	})
	if err == nil {
		t.Fatal("Error executing Create_ErrorInserting test: no error returned")
	}
}

func Test_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE roles").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	err = service.Update("test", &RoleUpdateRequest{
		VisibleName: "Testt",
	})
	if err != nil {
		t.Fatalf("Error executing Update test: %s\n", err.Error())
	}
}

func Test_Update_InvalidPayload(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	err = service.Update("test", &RoleUpdateRequest{
		VisibleName: "",
	})
	if err == nil {
		t.Fatal("Error executing Update_InvalidPayload test: no error returned")
	}
}

func Test_Update_NotExisting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE roles").WillReturnResult(sqlmock.NewResult(1, 0))
	mock.ExpectRollback()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	err = service.Update("test", &RoleUpdateRequest{
		VisibleName: "test",
	})
	if err == nil {
		t.Fatal("Error executing Update_NotExisting test: no error returned")
	}
}

func Test_Update_ErrorUpdating(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE roles").WillReturnResult(sqlmock.NewErrorResult(errors.New("Random error")))
	mock.ExpectRollback()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	err = service.Update("test", &RoleUpdateRequest{
		VisibleName: "test",
	})
	if err == nil {
		t.Fatal("Error executing Update_ErrorUpdating test: no error returned")
	}
}

func Test_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM roles WHERE (.+)").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	err = service.Delete("test")
	if err != nil {
		t.Fatalf("Error executing Delete test: %s\n", err.Error())
	}
}

func Test_Delete_ErrorDeleting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM roles WHERE (.+)").WillReturnResult(sqlmock.NewErrorResult(errors.New("Random error")))
	mock.ExpectRollback()

	service := &roleService{db: sqlx.NewDb(db, "sqlmock")}
	err = service.Delete("test")
	if err == nil {
		t.Fatal("Error executing Delete_ErrorDeleting test: no error returned")
	}
}

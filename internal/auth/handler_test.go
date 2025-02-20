package auth_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"short-link/configs"
	"short-link/internal/auth"
	"short-link/internal/user"
	"short-link/pkg/db"
	"testing"
)

func boorstrap() (*auth.AuthHandler, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, fmt.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbGorm, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	userRepo := user.NewUserRepository(&db.Db{
		DB: dbGorm,
	})
	handler := auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: auth.NewAuthService(userRepo),
	}
	return &handler, mock, nil
}

func TestLoginHandlerSuccess(t *testing.T) {
	handler, mock, err := boorstrap()
	rows := sqlmock.NewRows([]string{"email", "password"}).AddRow("Nels40@hotmail.com", "$2a$10$q6rbe11z4EwfVVl8AWatRu8px0zD.oHYMOdzwM73umdXckvBxT3jm")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "Nels40@hotmail.com",
		Password: "qwerty",
	})
	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)
	handler.Login()(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Result().StatusCode, http.StatusOK)
	}

}

func TestRegisterHandlerSuccess(t *testing.T) {
	handler, mock, err := boorstrap()
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	rows := sqlmock.NewRows([]string{"email", "password", "name"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	mock.ExpectCommit()

	data, _ := json.Marshal(&auth.RegisterRequest{
		Email:    "Nels40@hotmail.com",
		Password: "qwerty",
		Name:     "Nels",
	})
	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", reader)
	handler.Register()(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("got %d, want %d", w.Result().StatusCode, http.StatusCreated)
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"short-link/internal/auth"
	"short-link/internal/user"
	"testing"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env.dev")
	if err != nil {
		log.Fatalf("Error loading .env file. ERROR: %s", err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "test@test.com",
		Password: "$2a$10$5GPnccfDxlJh6s1Y1KGLAu9FbRciani7Oqy5njJmberuBAE47Uvcm",
		Name:     "test",
	})
}

func deleteData(db *gorm.DB) {
	db.Unscoped().Where("email = ?", "test@test.com").Delete(&user.User{})
}

func TestLoginSuccess(t *testing.T) {
	// Prepare
	db := initDb()
	initData(db)
	defer deleteData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@test.com",
		Password: "qwerty",
	})

	post, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	if post.StatusCode != http.StatusOK {
		t.Fatalf("expected status OK, got %d", post.StatusCode)
	}

	loginResponse := &auth.LoginResponse{}

	err = json.NewDecoder(post.Body).Decode(loginResponse)
	if err != nil {
		t.Fatal(err)
	}
	if loginResponse.Token == "" {
		t.Fatal("expected token, got empty string")
	}

}

func TestLoginError(t *testing.T) {
	initDb()
	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "a@a.a",
		Password: "",
	})

	post, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	if post.StatusCode != http.StatusUnauthorized && post.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 or 401, got %d", post.StatusCode)
	}

	//loginResponse := &auth.LoginResponse{}
	//
	//err = json.NewDecoder(post.Body).Decode(loginResponse)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if loginResponse.Token == "" {
	//	t.Fatal("expected token, got empty string")
	//}
}

package auth

import (
	"fmt"
	"log"
	"net/http"
	"short-link/configs"
	"short-link/pkg/req"
	"short-link/pkg/res"
)

type AuthHandler struct {
	*configs.Config
}

type AuthHandlerDeps struct {
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{Config: deps.Config}
	router.HandleFunc("POST /auth/login", handler.login())
	router.HandleFunc("POST /auth/register", handler.register())
}

func (handler *AuthHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			log.Println("Error handling request:", err)
			return
		}
		res.Json(w, body, http.StatusOK)
		fmt.Println("Login")
	}
}

func (handler *AuthHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			log.Println("Error handling request:", err)
			return
		}
		res.Json(w, body, http.StatusOK)
		fmt.Println("Register")
	}
}

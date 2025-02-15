package auth

import (
	"log"
	"net/http"
	"short-link/configs"
	"short-link/pkg/jwt"
	"short-link/pkg/req"
	"short-link/pkg/res"
)

type AuthHandler struct {
	*configs.Config
	*AuthService
}

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
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
		_, err = handler.AuthService.Login(body.Email, body.Password)
		if err != nil && err.Error() == ErrWrongCredentials {
			res.Json(w, err.Error(), http.StatusUnauthorized)
			return
		} else if err != nil {
			res.Json(w, ErrWrongCredentials, http.StatusUnauthorized)
			return
		}
		token, err := jwt.New(handler.Config.Auth.Secret).CreateToken(body.Email)
		if err != nil {
			res.Json(w, ErrWrongCredentials, http.StatusInternalServerError)
			return
		}
		data := LoginResponse{
			Token: token,
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *AuthHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			log.Println("Error handling request:", err)
			return
		}
		_, err = handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil && err.Error() == ErrUserExists {
			res.Json(w, err.Error(), http.StatusConflict)
			return
		} else if err != nil {
			log.Println("Error registering user:", err)
			res.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token, err := jwt.New(handler.Config.Auth.Secret).CreateToken(body.Email)
		if err != nil {
			res.Json(w, ErrWrongCredentials, http.StatusInternalServerError)
			return
		}
		data := RegisterResponse{
			Token: token,
			Email: body.Email,
		}
		res.Json(w, data, http.StatusOK)
	}
}

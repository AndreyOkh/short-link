package auth

import (
	"context"
	"io"
	"log"
	"net/http"
	"short-link/configs"
	"short-link/pkg/jwt"
	"short-link/pkg/req"
	"short-link/pkg/res"
	page "short-link/views/pages"
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
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
	router.HandleFunc("GET /auth", handler.AuthPage())
}

func (handler *AuthHandler) AuthPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		component := page.Auth()
		var render io.Writer = w
		err := component.Render(context.Background(), render)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (handler *AuthHandler) Login() http.HandlerFunc {
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
		token, err := jwt.New(handler.Config.Auth.Secret).CreateToken(jwt.JWTData{Email: body.Email})
		if err != nil {
			res.Json(w, ErrWrongCredentials, http.StatusInternalServerError)
			return
		}
		data := LoginResponse{
			Token: token,
		}
		//w.Header().Set("Authorization", "Bearer "+token)
		http.SetCookie(w, &http.Cookie{
			Name:  "Authorization",
			Value: "Bearer " + token,
			Path:  "/",
		})
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
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
		token, err := jwt.New(handler.Config.Auth.Secret).CreateToken(jwt.JWTData{Email: body.Email})
		if err != nil {
			res.Json(w, ErrWrongCredentials, http.StatusInternalServerError)
			return
		}
		data := RegisterResponse{
			Token: token,
			Email: body.Email,
		}
		res.Json(w, data, http.StatusCreated)
	}
}

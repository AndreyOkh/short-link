package middleware

import (
	"context"
	"log"
	"net/http"
	"short-link/configs"
	"short-link/pkg/jwt"
	"strings"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
)

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" || !strings.HasPrefix(tokenHeader, "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte("Unauthorized")); err != nil {
				log.Println("Error writing response:", err)
			}
			return
		}
		token := strings.TrimPrefix(tokenHeader, "Bearer ")
		isValid, data := jwt.New(config.Auth.Secret).ParseToken(token)
		if !isValid {
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte("Unauthorized")); err != nil {
				log.Println("Error writing response:", err)
			}
			return
		}
		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

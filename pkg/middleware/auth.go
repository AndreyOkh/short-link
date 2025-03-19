package middleware

import (
	"context"
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
		var tokenHeader string
		tokenCookie, err := r.Cookie("Authorization")
		if err != nil {
			tokenHeader = r.Header.Get("Authorization")
		} else {
			tokenHeader = tokenCookie.Value
		}

		if tokenHeader == "" || !strings.HasPrefix(tokenHeader, "Bearer") {
			//w.WriteHeader(http.StatusUnauthorized)
			//if _, err := w.Write([]byte("Unauthorized")); err != nil {
			//	log.Println("Error writing response:", err)
			//}
			http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
			return
		}
		token := strings.TrimPrefix(tokenHeader, "Bearer ")
		isValid, data := jwt.New(config.Auth.Secret).ParseToken(token)
		if !isValid {
			//w.WriteHeader(http.StatusUnauthorized)
			//if _, err := w.Write([]byte("Unauthorized")); err != nil {
			//	log.Println("Error writing response:", err)
			//}
			http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
			return
		}
		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

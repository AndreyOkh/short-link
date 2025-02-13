package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Log struct {
	RemoteAddr string        `json:"remote_addr"`
	Method     string        `json:"method"`
	Path       string        `json:"path"`
	Time       time.Duration `json:"time"`
	StatusCode int           `json:"status_code"`
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tNow := time.Now()
		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapper, r)

		// Обработка на обратном пути
		msg, _ := json.Marshal(&Log{
			RemoteAddr: r.RemoteAddr,
			Method:     r.Method,
			Path:       r.URL.Path,
			Time:       time.Since(tNow),
			StatusCode: wrapper.StatusCode,
		})
		log.Println(string(msg))
	})
}

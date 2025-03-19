package pages

import (
	"context"
	"io"
	"net/http"
	"short-link/views"
)

type Handler struct{}

func NewHandler(router *http.ServeMux) {
	handler := &Handler{}

	router.HandleFunc("GET /", handler.mainPage())
}

func (handler *Handler) mainPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		component := views.Main()
		var render io.Writer = w
		err := component.Render(context.Background(), render)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

package pages

import (
	"context"
	"io"
	"net/http"
	page "short-link/views/pages"
)

type Handler struct{}

func NewHandler(router *http.ServeMux) {
	handler := &Handler{}

	router.HandleFunc("GET /", handler.mainPage())
}

func (handler *Handler) mainPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		component := page.CreateLink()
		var render io.Writer = w
		err := component.Render(context.Background(), render)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

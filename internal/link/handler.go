package link

import (
	"log"
	"net/http"
	"short-link/pkg/req"
	"short-link/pkg/res"
)

type LinkHandler struct {
	LinkRepository *LinkRepository
}

type LinkHandlerDeps struct {
	LinkRepository *LinkRepository
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{LinkRepository: deps.LinkRepository}

	router.HandleFunc("POST /link", handler.create())
	router.HandleFunc("PATCH /link/{id}", handler.update())
	router.HandleFunc("DELETE /link/{id}", handler.delete())
	router.HandleFunc("GET /{hash}", handler.get())
}

func (handler *LinkHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[CreateRequest](&w, r)
		if err != nil {
			return
		}
		link := NewLink(body.URL)
		createdLink, err := handler.LinkRepository.Create(link)
		if err != nil {
			res.Json(w, "error creating link: "+err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, createdLink, http.StatusCreated)
	}
}

func (handler *LinkHandler) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (handler *LinkHandler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		log.Println("Delete request received", id)
	}
}

func (handler *LinkHandler) get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

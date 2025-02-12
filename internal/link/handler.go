package link

import (
	"log"
	"net/http"
	"short-link/pkg/req"
	"short-link/pkg/res"
	"strconv"
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
		link.generateHash(handler)
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
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			res.Json(w, "invalid id: "+err.Error(), http.StatusBadRequest)
			return
		}
		if err := handler.LinkRepository.DeleteByID(id); err != nil {
			res.Json(w, "error deleting link: "+err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, "link deleted", http.StatusOK)
	}
}

func (handler *LinkHandler) get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		log.Println("Get request received", hash)
		link, err := handler.LinkRepository.FindByHash(hash)
		if err != nil {
			res.Json(w, "error getting link: "+err.Error(), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, link.URL, http.StatusTemporaryRedirect)
	}
}

func (link *Link) generateHash(handler *LinkHandler) {
	for ok := false; !ok; {
		_, err := handler.LinkRepository.FindByHash(link.Hash)
		if err != nil && err.Error() == "record not found" {
			ok = true
			continue
		} else {
			link.Hash = RandStringRunes(6)
		}
	}
}

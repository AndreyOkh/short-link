package link

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"io"
	"net/http"
	"short-link/configs"
	"short-link/pkg/event"
	"short-link/pkg/middleware"
	"short-link/pkg/req"
	"short-link/pkg/res"
	"short-link/views"
	"strconv"
)

type LinkHandler struct {
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
}

type LinkHandlerDeps struct {
	LinkRepository *LinkRepository
	Config         *configs.Config
	EventBus       *event.EventBus
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
		EventBus:       deps.EventBus,
	}

	//router.Handle("POST /link", middleware.IsAuthed(handler.create(), deps.Config))
	router.HandleFunc("POST /link", handler.create())
	router.Handle("PATCH /link/{id}", middleware.IsAuthed(handler.update(), deps.Config))
	router.Handle("DELETE /link/{id}", middleware.IsAuthed(handler.delete(), deps.Config))
	router.HandleFunc("GET /{hash}", handler.get())
	router.Handle("GET /link", middleware.IsAuthed(handler.getAll(), deps.Config))
}

func (handler *LinkHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[CreateRequest](&w, r)
		if err != nil {
			return
		}
		link := NewLink(body.URL)
		for {
			existedLink, _ := handler.LinkRepository.FindByHash(link.Hash)
			if existedLink == nil {
				break
			}
			link.generateHash()
		}
		createdLink, err := handler.LinkRepository.Create(link)
		if err != nil {
			res.Json(w, "error creating link: "+err.Error(), http.StatusInternalServerError)
			return
		}

		shortUrl := r.Host + "/" + createdLink.Hash
		component := views.ShortUrl(createdLink.Hash, shortUrl)
		var render io.Writer = w
		err = component.Render(context.Background(), render)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		//res.Json(w, createdLink, http.StatusCreated)
	}
}

func (handler *LinkHandler) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			res.Json(w, "invalid id: "+err.Error(), http.StatusBadRequest)
			return
		}
		body, err := req.HandleBody[UpdateRequest](&w, r)
		if err != nil {
			res.Json(w, "error read request: "+err.Error(), http.StatusBadRequest)
			return
		}
		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{ID: uint(id)},
			URL:   body.URL,
			Hash:  body.Hash,
		})
		if err != nil {
			res.Json(w, "error updating link: "+err.Error(), http.StatusInternalServerError)
			return
		}
		emailUser, ok := r.Context().Value(middleware.ContextEmailKey).(string)
		if ok {
			fmt.Println("USER EMAIL: ", emailUser)
		}

		res.Json(w, link, http.StatusCreated)
	}
}

func (handler *LinkHandler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			res.Json(w, "invalid id: "+err.Error(), http.StatusBadRequest)
			return
		}
		if err := handler.LinkRepository.DeleteByID(id); err != nil {
			res.Json(w, "error deleting link: "+err.Error(), http.StatusNotFound)
			return
		}
		res.Json(w, "link deleted", http.StatusOK)
	}
}

func (handler *LinkHandler) get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		link, err := handler.LinkRepository.FindByHash(hash)
		if err != nil {
			res.Json(w, "error getting link: "+err.Error(), http.StatusNotFound)
			return
		}
		//handler.StatRepository.AddClick(link.ID)
		go handler.EventBus.Publish(event.Event{
			Type: event.LinkVisitedEvent,
			Data: link.ID,
		})
		http.Redirect(w, r, link.URL, http.StatusTemporaryRedirect)
	}
}

func (handler *LinkHandler) getAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			res.Json(w, "invalid limit: "+err.Error(), http.StatusBadRequest)
			return
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			res.Json(w, "invalid offset: "+err.Error(), http.StatusBadRequest)
			return
		}
		links, err := handler.LinkRepository.GetLinks(offset, limit)
		if err != nil {
			res.Json(w, "error getting links: "+err.Error(), http.StatusInternalServerError)
			return
		}
		countLinks, err := handler.LinkRepository.Count()
		if err != nil {
			res.Json(w, "error getting links: "+err.Error(), http.StatusInternalServerError)
			return
		}
		response := &GetAllLinkResponse{
			Count: countLinks,
			Links: links,
		}
		res.Json(w, response, http.StatusOK)
	}
}

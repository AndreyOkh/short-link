package stat

import (
	"net/http"
	"short-link/configs"
	"short-link/pkg/middleware"
	"short-link/pkg/res"
	"time"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

type StatHandler struct {
	StatRepository *StatRepository
}

type StatHandlerDeps struct {
	StatRepository *StatRepository
	Config         *configs.Config
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepository: deps.StatRepository,
	}

	router.Handle("GET /stat", middleware.IsAuthed(handler.getStat(), deps.Config))
}

func (handler *StatHandler) getStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fromParam, err := time.Parse("2006-01-02", r.URL.Query().Get("from"))
		if err != nil {
			res.Json(w, "invalid 'from' param: "+err.Error(), http.StatusBadRequest)
			return
		}
		toParam, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
		if err != nil {
			res.Json(w, "invalid 'to' param: "+err.Error(), http.StatusBadRequest)
			return
		}
		byParam := r.URL.Query().Get("by")
		if !(byParam == GroupByMonth) && !(byParam == GroupByDay) {
			res.Json(w, "invalid 'by' param, use only 'day' or 'month': "+err.Error(), http.StatusBadRequest)
			return
		}
		clicks := handler.StatRepository.GetStats(byParam, fromParam, toParam)
		res.Json(w, clicks, http.StatusOK)
	}
}

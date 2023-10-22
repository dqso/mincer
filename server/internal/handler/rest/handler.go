package rest

import (
	"context"
	"github.com/dqso/mincer/server/internal/log"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	logger  log.Logger
	usecase usecase
	router  *mux.Router
}

type usecase interface {
	AcquireToken(ctx context.Context) (uint64, []byte, error)
}

func NewHandler(logger log.Logger, usecase usecase) *Handler {
	h := &Handler{
		logger:  logger.With(log.Module("rest_handler")),
		usecase: usecase,
	}
	h.router = mux.NewRouter()

	h.router.HandleFunc("/token", h.AcquireToken).Methods(http.MethodPost)

	return h
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

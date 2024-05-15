package httpHandlers

import (
	"github.com/Sanchir01/Grasp/internal/config"
	"github.com/Sanchir01/Grasp/internal/db/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"log/slog"
	"net/http"
)

type Router struct {
	chiRouter *chi.Mux
	config    *config.Config
	storage   *storage.Storage
	logger    *slog.Logger
}

func NewChiRouter(chi *chi.Mux, cfg *config.Config, storage *storage.Storage, logger *slog.Logger) *Router {
	return &Router{chiRouter: chi, config: cfg, storage: storage, logger: logger}
}

func (r *Router) StartHttpHandlers() http.Handler {
	r.routerCors()

	r.ProductsHandlers()

	return r.chiRouter
}

func (r *Router) routerCors() {
	r.chiRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge:           300,
	}))
}

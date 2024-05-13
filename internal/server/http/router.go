package httpHandlers

import (
	"github.com/Sanchir01/Grasp/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

type Router struct {
	chiRouter *chi.Mux
	config    *config.Config
}

func NewChiRouter(chi *chi.Mux) *Router {
	return &Router{chiRouter: chi}
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

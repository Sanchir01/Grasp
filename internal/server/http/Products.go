package httpHandlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (r *Router) ProductsHandlers() {
	r.chiRouter.Route("/products", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("AllProducts"))
		})

		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			w.Write([]byte(id))
		})
	})
}

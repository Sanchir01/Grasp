package httpHandlers

import (
	"context"
	"github.com/Sanchir01/Grasp/internal/db/model"
	resp "github.com/Sanchir01/Grasp/pkg/lib/api/response"
	"github.com/Sanchir01/Grasp/pkg/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	Name  string `json:"name" validate:"required"`
	Price int    `json:"price" validate:"required"`
}

type Response struct {
	resp.Response
	Product model.Products
}

func (rout *Router) ProductsHandlers() {
	rout.chiRouter.Route("/products", func(r chi.Router) {
		r.Get("/", rout.getAllProduct)
		r.Post("/", rout.createProduct)
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			w.Write([]byte(id))
		})
	})
}

func (rout *Router) getAllProduct(w http.ResponseWriter, r *http.Request) {
	products, err := rout.storage.GetAllProducts(context.Background())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	render.JSON(w, r, products)

}

func (rout *Router) createProduct(writer http.ResponseWriter, request *http.Request) {
	var req Request
	err := render.DecodeJSON(request.Body, &req)

	if err != nil {
		rout.logger.Error("failed to decode json body", err.Error())
		render.JSON(writer, request, resp.Error("failed to decode request body"))
		return
	}
	rout.logger.Info("request decoded", slog.Any("request", req))

	if err := validator.New().Struct(req); err != nil {
		validatorErr := err.(validator.ValidationErrors)

		rout.logger.Error("validation error", sl.Err(err))

		render.JSON(writer, request, resp.Error("invalid request"))

		render.JSON(writer, request, resp.ValidationError(validatorErr))
		return
	}
	id, err := rout.storage.CreateProduct(context.Background(), model.Products{
		Name:  req.Name,
		Price: req.Price,
	})

	if err != nil {
		rout.logger.Warn("Product already exists", slog.String("id", err.Error()))

		render.JSON(writer, request, resp.Error("Product already exists"))

		return
	}
	rout.logger.Info("product created", slog.Int64("id", id))
	render.JSON(writer, request, Response{
		Response: resp.OK(),
		Product: model.Products{
			Price: req.Price,
			Name:  req.Name,
		},
	})
}

package httpHandlers

import (
	"context"
	"github.com/Sanchir01/Grasp/internal/db/model"
	resp "github.com/Sanchir01/Grasp/pkg/lib/api/response"
	"github.com/Sanchir01/Grasp/pkg/lib/logger/sl"
	"github.com/Sanchir01/Grasp/pkg/lib/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	Name         string `json:"name" validate:"required"`
	Price        int32  `json:"price" validate:"required"`
	CategoryName string `json:"category_name" validate:"required"`
}

type ResponseProduct struct {
	resp.Response
	Product model.Products
}

func (rout *Router) ProductsHandlers() {
	rout.chiRouter.Route("/products", func(r chi.Router) {
		r.Get("/", rout.getAllProduct)
		r.Post("/", rout.createProduct)
		r.Get("/op/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			w.Write([]byte(id))
		})
	})
}

func (rout *Router) getAllProduct(w http.ResponseWriter, r *http.Request) {
	products, err := rout.productsStr.GetAllProducts(context.Background())
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
	categorySlug, err := utils.Slugify(req.CategoryName)
	if err != nil {
		render.JSON(writer, request, resp.Error("Вы не ввели категорию товара"))
	}
	category, _ := rout.categoriesStr.GetCategoryBySlug(context.Background(), categorySlug)
	if category == nil {
		render.JSON(writer, request, resp.Error("category not found"))
		return
	}
	id, err := rout.productsStr.CreateProduct(context.Background(), model.Products{
		Name:       req.Name,
		Price:      req.Price,
		CategoryId: category.Id,
	})

	if err != nil {
		rout.logger.Warn("Product already exists", slog.String("id", err.Error()))

		render.JSON(writer, request, resp.Error("Product already exists"))

		return
	}
	rout.logger.Info("product created", slog.Int64("id", id))
	render.JSON(writer, request, ResponseProduct{
		Response: resp.OK(),
		Product: model.Products{
			Price:      req.Price,
			Name:       req.Name,
			CategoryId: category.Id,
		},
	})
}

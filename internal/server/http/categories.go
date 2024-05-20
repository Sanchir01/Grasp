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

type CrateCategoryRequest struct {
	CategoryName string `json:"name"`
	Description  string `json:"description"`
}
type ResponseCategory struct {
	resp.Response
	Category model.Categories
}

func (rout *Router) CategoriesHandlers() {
	rout.chiRouter.Route("/categories", func(r chi.Router) {
		r.Get("/", rout.getAllCategories)
		r.Post("/", rout.createCategories)
		r.Get("/{slug}", func(w http.ResponseWriter, r *http.Request) {
			slug := chi.URLParam(r, "slug")
			w.Write([]byte(slug))
		})
	})
}

func (rout *Router) getAllCategories(writer http.ResponseWriter, request *http.Request) {
	categrories, err := rout.categoriesStr.GetAllCategories(context.Background())
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	render.JSON(writer, request, categrories)
}

func (rout *Router) createCategories(writer http.ResponseWriter, request *http.Request) {
	var req CrateCategoryRequest
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
	slugCategory, err := utils.Slugify(req.CategoryName)
	if err != nil {
		render.JSON(writer, request, resp.Error("Вы не ввели категорию товара"))
		return
	}
	id, err := rout.categoriesStr.CreateCategories(context.Background(), model.Categories{
		Name:        req.CategoryName,
		Slug:        slugCategory,
		Description: req.Description,
	})
	if err != nil {
		return
	}
	rout.logger.Info("product created", slog.Int64("id", id))

	render.JSON(writer, request, ResponseCategory{
		Response: resp.Response{Status: "ok"},
		Category: model.Categories{
			Name:        req.CategoryName,
			Slug:        slugCategory,
			Description: req.Description,
		},
	})
}

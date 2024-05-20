package storage

import (
	"context"
	"github.com/Sanchir01/Grasp/internal/db/model"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"time"
)

type CategoriesPostgresStorage struct {
	db *sqlx.DB
}

func NewCategoriesStorage(db *sqlx.DB) *CategoriesPostgresStorage {
	return &CategoriesPostgresStorage{db: db}
}

func (s *CategoriesPostgresStorage) GetAllCategories(ctx context.Context) ([]model.Categories, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var categories []dbCategories
	if err = conn.SelectContext(ctx, &categories, "SELECT * FROM categories"); err != nil {
		return nil, err
	}

	return lo.Map(categories, func(categorie dbCategories, _ int) model.Categories {
		return model.Categories(categorie)
	}), nil
}

func (s *CategoriesPostgresStorage) CreateCategories(ctx context.Context, categories model.Categories) (int64, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	var id int64

	row := conn.QueryRowContext(
		ctx,
		`INSERT INTO categories(name, slug, description) VALUES($1, $2, $3) RETURNING id`,
		categories.Name, categories.Slug, categories.Description,
	)
	if err := row.Err(); err != nil {
		return 0, err
	}
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *CategoriesPostgresStorage) GetCategoryBySlug(ctx context.Context, slug string) (*model.Categories, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var category dbCategories
	if err := conn.GetContext(ctx, &category, "SELECT * FROM categories WHERE slug = $1", slug); err != nil {
		return nil, err
	}
	return (*model.Categories)(&category), nil
}

type dbCategories struct {
	Id          int32     `db:"id"`
	Name        string    `db:"name"`
	Slug        string    `db:"slug"`
	Description string    `db:"description"`
	UpdatedAt   time.Time `db:"updated_at"`
	CreatedAt   time.Time `db:"created_at"`
}

package storage

import (
	"context"
	"github.com/Sanchir01/Grasp/internal/db/model"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"time"
)

type ProductPostgresStorage struct {
	db *sqlx.DB
}

func NewProductStorage(db *sqlx.DB) *ProductPostgresStorage {
	return &ProductPostgresStorage{db: db}
}
func (s *ProductPostgresStorage) GetAllProducts(ctx context.Context) ([]model.Products, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var products []dbProducts

	if err := conn.SelectContext(ctx, &products, "SELECT * FROM products"); err != nil {
		return nil, err
	}
	return lo.Map(products, func(products dbProducts, _ int) model.Products {
		return model.Products(products)
	}), nil
}

func (s *ProductPostgresStorage) CreateProduct(ctx context.Context, product model.Products) (int64, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	var id int64

	row := conn.QueryRowContext(ctx,
		"INSERT INTO products(name, price,category_id) VALUES($1, $2, $3) RETURNING id",
		product.Name,
		product.Price,
		product.CategoryId,
	)

	if err := row.Err(); err != nil {
		return 0, err
	}
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

type dbProducts struct {
	Id         int32     `db:"id"`
	Name       string    `db:"name"`
	Price      int32     `db:"price"`
	CategoryId int32     `db:"category_id"`
	UpdatedAt  time.Time `db:"updated_at"`
	CreatedAt  time.Time `db:"created_at"`
}

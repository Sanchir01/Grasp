package storage

import (
	"context"
	"github.com/Sanchir01/Grasp/internal/db/model"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"time"
)

type Storage struct {
	db *sqlx.DB
}

func NewProductStorage(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}
func (s *Storage) GetAllProducts(ctx context.Context) ([]model.Products, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var products []DBProducts

	if err := conn.SelectContext(ctx, &products, "SELECT * FROM products"); err != nil {
		return nil, err
	}
	return lo.Map(products, func(products DBProducts, _ int) model.Products {
		return model.Products(products)
	}), nil
}

func (s *Storage) CreateProduct(ctx context.Context, product model.Products) (int64, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	var id int64

	row := conn.QueryRowContext(ctx,
		"INSERT INTO products(name, price) VALUES($1, $2) RETURNING id",
		product.Name,
		product.Price,
	)

	if err := row.Err(); err != nil {
		return 0, err
	}
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

type DBProducts struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	Price     int       `db:"price"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}

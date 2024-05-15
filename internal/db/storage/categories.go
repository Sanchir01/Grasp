package storage

import "github.com/jmoiron/sqlx"

func NewCategoriesStorage(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

package model

import "time"

type Products struct {
	Id        int
	Name      string
	Price     int
	UpdatedAt time.Time
	CreatedAt time.Time
}

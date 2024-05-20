package model

import "time"

type Products struct {
	Id         int32
	Name       string
	Price      int32
	CategoryId int32
	UpdatedAt  time.Time
	CreatedAt  time.Time
}

type Categories struct {
	Id          int32
	Name        string
	Slug        string
	Description string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

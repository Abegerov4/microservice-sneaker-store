package model

import "time"

type Product struct {
	ID          string
	Name        string
	Brand       string
	Description string
	Price       float64
	Sizes       []string
	Stock       int
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductStats struct {
	TotalProducts int
	TotalBrands   int
	TotalStock    int
	AveragePrice  float64
}

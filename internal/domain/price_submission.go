package domain

import "time"

type PriceSubmission struct {
	ID          int64     `json:"id"`
	ProductName string    `json:"product_name"`
	StoreName   string    `json:"store_name"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	PhotoURL    string    `json:"photo_url,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	ObservedAt  time.Time `json:"observed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

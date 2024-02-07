package models

import "time"

type Order struct {
	ID            int64     `json:"id"`
	Price         float64   `json:"price"`
	DeliveryPrice float64   `json:"delivery_price"`
	Lng           float64   `json:"lng"`
	Lat           float64   `json:"lat"`
	IsDelivered   bool      `json:"is_delivere"`
	CreatedAt     time.Time `json:"created_at"`
}

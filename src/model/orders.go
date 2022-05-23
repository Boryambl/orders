package model

import "time"

type Order struct {
	UUID      string           `json:"uuid"`
	Number    int              `json:"number"`
	UserID    string           `json:"user_id"`
	Products  []ProductInOrder `json:"products"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt *time.Time       `json:"updated_at,omitempty"`
}

type ProductToShow struct {
	Product string `json:"product"`
	Count   int    `json:"count"`
}

type ProductInOrder struct {
	ProductID string `json:"product_id"`
	Count     int    `json:"count"`
}

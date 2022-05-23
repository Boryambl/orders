package model

import "time"

type Product struct {
	UUID        string     `json:"uuid"`
	Description string     `json:"description,omitempty"`
	Price       []Price    `json:"price"`
	LeftInStock int        `json:"left_in_stock"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type Price struct {
	Сurrency string `json:"currency"`
	Value    int    `json:"value"`
}

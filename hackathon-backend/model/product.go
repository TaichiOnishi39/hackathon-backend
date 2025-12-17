package model

import "time"

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Price       int       `json:"price"`
	Description string    `json:"description"`
	UserID      string    `json:"user_id"` // ここは User の ID を入れる
	UserName    string    `json:"user_name"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductReq struct {
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}

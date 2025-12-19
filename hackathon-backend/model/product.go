package model

import "time"

type Product struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Price         int       `json:"price"`
	Description   string    `json:"description"`
	UserID        string    `json:"user_id"` // ここは User の ID を入れる
	UserName      string    `json:"user_name"`
	ImageURL      string    `json:"image_url"`
	BuyerID       string    `json:"buyer_id"`
	CreatedAt     time.Time `json:"created_at"`
	LikeCount     int       `json:"like_count"`
	IsLiked       bool      `json:"is_liked"`
	BuyerName     string    `json:"buyer_name"`
	UserImageURL  string    `json:"user_image_url"`
	BuyerImageURL string    `json:"buyer_image_url"`
}

type ProductPage struct {
	Products []*Product `json:"products"`
	Total    int        `json:"total"`
}

type ProductReq struct {
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}

// AI商品説明生成のリクエスト
type GenerateReq struct {
	Name     string `json:"name"`
	Keywords string `json:"keywords"`
}

// AI商品説明生成のレスポンス
type GenerateRes struct {
	Description string `json:"description"`
}

type GenerateImageRes struct {
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Keywords    string `json:"keywords"`
	Description string `json:"description"`
}

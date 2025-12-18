package model

import "time"

type Message struct {
	ID          string    `json:"id"`
	SenderID    string    `json:"sender_id"`
	ReceiverID  string    `json:"receiver_id"`
	Content     string    `json:"content"`
	ProductID   string    `json:"product_id,omitempty"`
	ProductName string    `json:"product_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// 送信するときのリクエスト用
type SendMessageReq struct {
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
	ProductID  string `json:"product_id"`
}

// ChatListRes: チャット一覧画面用のレスポンス
type ChatListRes struct {
	PartnerID   string    `json:"partner_id"`
	PartnerName string    `json:"partner_name"`
	LastMessage string    `json:"last_message"`
	LastTime    time.Time `json:"last_time"`
}

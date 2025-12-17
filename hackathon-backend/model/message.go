package model

import "time"

type Message struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

// 送信するときのリクエスト用
type SendMessageReq struct {
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
}

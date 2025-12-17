package dao

import (
	"database/sql"
	"hackathon-backend/model"
)

type MessageDao struct {
	db *sql.DB
}

func NewMessageDao(db *sql.DB) *MessageDao {
	return &MessageDao{db: db}
}

// Create: メッセージを保存
func (d *MessageDao) Create(msg *model.Message) error {
	query := `
		INSERT INTO messages (id, sender_id, receiver_id, content, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := d.db.Exec(query, msg.ID, msg.SenderID, msg.ReceiverID, msg.Content, msg.CreatedAt)
	return err
}

// GetMessagesBetween: 2人の間のメッセージを時系列順に取得
func (d *MessageDao) GetMessagesBetween(userA, userB string) ([]*model.Message, error) {
	// Aが送ってBが受け取った or Bが送ってAが受け取った メッセージを取得
	query := `
		SELECT id, sender_id, receiver_id, content, created_at
		FROM messages
		WHERE (sender_id = ? AND receiver_id = ?) 
		   OR (sender_id = ? AND receiver_id = ?)
		ORDER BY created_at ASC
	`
	rows, err := d.db.Query(query, userA, userB, userB, userA)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		m := &model.Message{}
		if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

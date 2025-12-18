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
		INSERT INTO messages (id, sender_id, receiver_id, content, created_at, product_id, is_read)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := d.db.Exec(query, msg.ID, msg.SenderID, msg.ReceiverID, msg.Content, msg.CreatedAt, msg.ProductID, msg.IsRead)
	return err
}

// GetMessagesBetween: 2人の間のメッセージを時系列順に取得
func (d *MessageDao) GetMessagesBetween(userA, userB string) ([]*model.Message, error) {
	// Aが送ってBが受け取った or Bが送ってAが受け取った メッセージを取得
	query := `
	   SELECT 
           m.id, m.sender_id, m.receiver_id, m.content, m.created_at, 
           m.product_id, p.name, m.is_read
       FROM messages m
       LEFT JOIN products p ON m.product_id = p.id
       WHERE (m.sender_id = ? AND m.receiver_id = ?) 
          OR (m.sender_id = ? AND m.receiver_id = ?)
       ORDER BY m.created_at ASC
	`
	rows, err := d.db.Query(query, userA, userB, userB, userA)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		m := &model.Message{}
		var productID sql.NullString
		var productName sql.NullString
		if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt, &productID, &productName, &m.IsRead); err != nil {
			return nil, err
		}

		if productID.Valid {
			m.ProductID = productID.String
		}
		if productName.Valid {
			m.ProductName = productName.String
		}
		messages = append(messages, m)
	}
	return messages, nil
}

// FindAllByUserID: 自分に関係する全てのメッセージを新しい順に取得
func (d *MessageDao) FindAllByUserID(userID string) ([]*model.Message, error) {
	query := `
	   SELECT 
           m.id, m.sender_id, m.receiver_id, m.content, m.created_at,
           m.product_id, p.name, m.is_read
       FROM messages m
       LEFT JOIN products p ON m.product_id = p.id
       WHERE m.sender_id = ? OR m.receiver_id = ? 
       ORDER BY m.created_at DESC
	`
	rows, err := d.db.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		m := &model.Message{}
		var productID sql.NullString
		var productName sql.NullString
		if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt, &productID, &productName, &m.IsRead); err != nil {
			return nil, err
		}

		if productID.Valid {
			m.ProductID = productID.String
		}
		if productName.Valid {
			m.ProductName = productName.String
		}
		messages = append(messages, m)
	}
	return messages, nil
}

// 特定の相手からのメッセージを全て既読にする
func (d *MessageDao) MarkAsRead(myUserID, partnerID string) error {
	// 自分が受信者(Receiver)で、相手が送信者(Sender)のメッセージを既読(TRUE)にする
	query := `UPDATE messages SET is_read = TRUE WHERE receiver_id = ? AND sender_id = ? AND is_read = FALSE`
	_, err := d.db.Exec(query, myUserID, partnerID)
	return err
}

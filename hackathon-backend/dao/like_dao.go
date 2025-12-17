package dao

import (
	"database/sql"
)

type LikeDao struct {
	db *sql.DB
}

func NewLikeDao(db *sql.DB) *LikeDao {
	return &LikeDao{db: db}
}

// AddLike: いいねを追加
func (d *LikeDao) AddLike(userID, productID string) error {
	// IGNORE: 既にいいね済みならエラーにせず無視
	query := "INSERT IGNORE INTO likes (user_id, product_id) VALUES (?, ?)"
	_, err := d.db.Exec(query, userID, productID)
	return err
}

// RemoveLike: いいねを解除
func (d *LikeDao) RemoveLike(userID, productID string) error {
	query := "DELETE FROM likes WHERE user_id = ? AND product_id = ?"
	_, err := d.db.Exec(query, userID, productID)
	return err
}

// HasLiked: 自分がいいねしているか確認
func (d *LikeDao) HasLiked(userID, productID string) (bool, error) {
	query := "SELECT COUNT(*) FROM likes WHERE user_id = ? AND product_id = ?"
	var count int
	err := d.db.QueryRow(query, userID, productID).Scan(&count)
	if err != nil {
		return false, err
	}
	// 1件以上あれば true
	return count > 0, nil
}

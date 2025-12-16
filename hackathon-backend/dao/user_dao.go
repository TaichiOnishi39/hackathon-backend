package dao

import (
	"database/sql"
	"fmt"

	"hackathon-backend/model"
)

type UserDao struct {
	DB *sql.DB
}

func NewUserDao(db *sql.DB) *UserDao {
	return &UserDao{DB: db}
}

func (dao *UserDao) FindByName(name string) ([]model.UserResForHTTPGet, error) {
	rows, err := dao.DB.Query("SELECT id, name, firebase_uid FROM users WHERE name = ?", name)
	if err != nil {
		return nil, fmt.Errorf("fail: db.Query, %v", err)
	}
	defer rows.Close()

	users := make([]model.UserResForHTTPGet, 0)
	for rows.Next() {
		var u model.UserResForHTTPGet
		if err := rows.Scan(&u.Id, &u.Name, &u.FirebaseUID); err != nil {
			return nil, fmt.Errorf("fail: rows.Scan, %v", err)
		}
		users = append(users, u)
	}
	return users, nil
}

func (dao *UserDao) RegisterUser(ulid string, firebaseUID string, name string) (*model.User, error) {
	tx, err := dao.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("fail: db.Begin, %v", err)
	}
	defer tx.Rollback()

	// 1. 新規ならULIDを使ってINSERT、既存なら名前だけUPDATE
	// (既存の場合、ここで渡したulidは無視されます)
	_, err = tx.Exec(`
		INSERT INTO users (id, firebase_uid, name) VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE name = ?
	`, ulid, firebaseUID, name, name)

	if err != nil {
		return nil, fmt.Errorf("fail: tx.Exec, %v", err)
	}

	// 2. 実際にDBに入っているIDを取得する
	// (新規作成ならさっきのULID、既存なら昔作られたULIDが返ってくる)
	var user model.User
	err = tx.QueryRow("SELECT id, name, firebase_uid FROM users WHERE firebase_uid = ?", firebaseUID).
		Scan(&user.Id, &user.Name, &user.FirebaseUID)
	if err != nil {
		return nil, fmt.Errorf("fail: tx.QueryRow, %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("fail: tx.Commit, %v", err)
	}

	return &user, nil
}

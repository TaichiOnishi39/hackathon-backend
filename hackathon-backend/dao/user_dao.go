package dao

import (
	"database/sql"
	"fmt"

	"hackathon-backend/model"
)

type UserDao struct {
	db *sql.DB
}

func NewUserDao(db *sql.DB) *UserDao {
	return &UserDao{db: db}
}

func (dao *UserDao) FindByName(name string) ([]model.UserRes, error) {
	rows, err := dao.db.Query("SELECT id, name, firebase_uid FROM users WHERE name = ?", name)
	if err != nil {
		return nil, fmt.Errorf("fail: db.Query, %v", err)
	}
	defer rows.Close()

	users := make([]model.UserRes, 0)
	for rows.Next() {
		var u model.UserRes
		if err := rows.Scan(&u.ID, &u.Name, &u.FirebaseUID); err != nil {
			return nil, fmt.Errorf("fail: rows.Scan, %v", err)
		}
		users = append(users, u)
	}
	return users, nil
}

func (dao *UserDao) FindByFirebaseUID(firebaseUID string) (*model.User, error) {
	var user model.User
	// 1件だけ取得するので QueryRow を使います
	row := dao.db.QueryRow("SELECT id, name, firebase_uid FROM users WHERE firebase_uid = ?", firebaseUID)

	if err := row.Scan(&user.ID, &user.Name, &user.FirebaseUID); err != nil {
		if err == sql.ErrNoRows {
			// ユーザーが見つからない場合は nil, nil を返す設計にします
			// (呼び出し元の Usecase や Controller で 404 エラーにするため)
			return nil, nil
		}
		// その他のDBエラー
		return nil, fmt.Errorf("fail: row.Scan, %v", err)
	}

	return &user, nil
}

func (dao *UserDao) FindByID(id string) (*model.User, error) {
	var user model.User
	row := dao.db.QueryRow("SELECT id, name, firebase_uid FROM users WHERE id = ?", id)
	if err := row.Scan(&user.ID, &user.Name, &user.FirebaseUID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (dao *UserDao) CreateOrUpdate(ulid string, firebaseUID string, name string) (*model.User, error) {
	tx, err := dao.db.Begin()
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
	// 確定したユーザー情報を取得して返す
	var user model.User
	err = tx.QueryRow("SELECT id, name, firebase_uid FROM users WHERE firebase_uid = ?", firebaseUID).
		Scan(&user.ID, &user.Name, &user.FirebaseUID)
	if err != nil {
		return nil, fmt.Errorf("fail: tx.QueryRow, %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("fail: tx.Commit, %v", err)
	}

	return &user, nil
}

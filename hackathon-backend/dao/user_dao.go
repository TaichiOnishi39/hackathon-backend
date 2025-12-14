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
	rows, err := dao.DB.Query("SELECT id, name, age FROM user WHERE name = ?", name)
	if err != nil {
		return nil, fmt.Errorf("fail: db.Query, %v", err)
	}
	defer rows.Close()

	users := make([]model.UserResForHTTPGet, 0)
	for rows.Next() {
		var u model.UserResForHTTPGet
		if err := rows.Scan(&u.Id, &u.Name, &u.Age); err != nil {
			return nil, fmt.Errorf("fail: rows.Scan, %v", err)
		}
		users = append(users, u)
	}
	return users, nil
}

func (dao *UserDao) Insert(id string, name string, age int) error {
	tx, err := dao.DB.Begin()
	if err != nil {
		return fmt.Errorf("fail: db.Begin, %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO user (id, name, age) VALUES (?, ?, ?)", id, name, age)
	if err != nil {
		return fmt.Errorf("fail: tx.Exec, %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("fail: tx.Commit, %v", err)
	}
	return nil
}

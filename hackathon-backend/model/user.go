package model

import (
	"errors"
	"fmt"
)

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	FirebaseUID string `json:"firebase_uid"`
}
type CreateUserReq struct {
	Name string `json:"name"`
}

func (u *CreateUserReq) Validate() error {
	if u.Name == "" {
		return errors.New("name is empty")
	}
	if len(u.Name) > 50 {
		return fmt.Errorf("name is too long: max 50 chars, but got %d", len(u.Name))
	}
	return nil
}

type UserRes struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	FirebaseUID string `json:"firebase_uid"`
}

package model

import (
	"errors"
	"fmt"
)

type UserResForHTTPGet struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	FirebaseUID string `json:"firebase_uid"`
}
type UserReqForHTTPPost struct {
	Name string `json:"name"`
}

func (u *UserReqForHTTPPost) Validate() error {
	if u.Name == "" {
		return errors.New("name is empty")
	}
	if len(u.Name) > 50 {
		return fmt.Errorf("name is too long: max 50 chars, but got %d", len(u.Name))
	}
	return nil
}

type UserResForHTTPPost struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	FirebaseUID string `json:"firebase_uid"`
}
type User struct {
	Id          string
	Name        string
	FirebaseUID string
}

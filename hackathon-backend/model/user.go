package model

import (
	"errors"
	"fmt"
)

type UserResForHTTPGet struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type UserReqForHTTPPost struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (u *UserReqForHTTPPost) Validate() error {
	if u.Name == "" {
		return errors.New("name is empty")
	}
	if len(u.Name) > 50 {
		return fmt.Errorf("name is too long: max 50 chars, but got %d", len(u.Name))
	}

	if u.Age < 20 || u.Age > 80 {
		return fmt.Errorf("age must be between 20 and 80, but got %d", u.Age)
	}
	return nil
}

type UserResForHTTPPost struct {
	Id string `json:"id"`
}
type User struct {
	Id   string
	Name string
	Age  int
}

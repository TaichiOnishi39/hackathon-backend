package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"hackathon-backend/usecase"
)

type SearchUserController struct {
	Usecase *usecase.SearchUserUsecase
}

func NewSearchUserController(u *usecase.SearchUserUsecase) *SearchUserController {
	return &SearchUserController{Usecase: u}
}

func (c *SearchUserController) Handle(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		log.Println("fail: name is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := c.Usecase.SearchUser(name)
	if err != nil {
		log.Printf("fail: SearchUser, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(users)
	if err != nil {
		log.Printf("fail: json.Marshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

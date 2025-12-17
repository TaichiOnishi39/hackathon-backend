package controller

import (
	"encoding/json"
	"hackathon-backend/model"
	"hackathon-backend/usecase"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
)

type RegisterUserController struct {
	BaseController
	Usecase *usecase.RegisterUserUsecase
}

func NewRegisterUserController(u *usecase.RegisterUserUsecase, auth *auth.Client) *RegisterUserController {
	return &RegisterUserController{
		BaseController: BaseController{AuthClient: auth},
		Usecase:        u,
	}
}

func (c *RegisterUserController) Handle(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	var reqBody model.CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("fail: json.NewDecoder, %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Usecase実行
	resUser, err := c.Usecase.RegisterUser(reqBody, firebaseUID)
	if err != nil {
		// 特定のエラーハンドリング
		if strings.Contains(err.Error(), "name is empty") || strings.Contains(err.Error(), "too long") {
			c.respondError(w, http.StatusBadRequest, err)
			return
		}
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	// 共通関数でレスポンス
	c.respondJSON(w, http.StatusOK, resUser)
}

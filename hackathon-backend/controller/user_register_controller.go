package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"hackathon-backend/model"
	"hackathon-backend/usecase"
)

type RegisterUserController struct {
	Usecase *usecase.RegisterUserUsecase
}

func NewRegisterUserController(u *usecase.RegisterUserUsecase) *RegisterUserController {
	return &RegisterUserController{Usecase: u}
}

func (c *RegisterUserController) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("fail: io.ReadAll, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var reqBody model.UserReqForHTTPPost
	if err := json.Unmarshal(body, &reqBody); err != nil {
		log.Printf("fail: json.Unmarshal, %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Usecase実行
	newIDStr, err := c.Usecase.RegisterUser(reqBody)
	if err != nil {
		// エラー内容に応じてステータスコードを分ける簡易実装
		if err.Error() == "name is empty or too long" || err.Error() == "age must be between 20 and 80" {
			log.Printf("fail: validation, %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err.Error())
			return
		}
		log.Printf("fail: RegisterUser, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resBody := model.UserResForHTTPPost{Id: newIDStr}
	bytes, err := json.Marshal(resBody)
	if err != nil {
		log.Printf("fail: json.Marshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

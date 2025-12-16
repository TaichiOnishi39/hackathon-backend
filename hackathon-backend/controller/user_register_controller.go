package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"hackathon-backend/model"
	"hackathon-backend/usecase"

	"firebase.google.com/go/auth"
)

type RegisterUserController struct {
	Usecase    *usecase.RegisterUserUsecase
	AuthClient *auth.Client
}

func NewRegisterUserController(u *usecase.RegisterUserUsecase, auth *auth.Client) *RegisterUserController {
	return &RegisterUserController{Usecase: u, AuthClient: auth}
}

func (c *RegisterUserController) Handle(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	idToken := strings.Replace(authHeader, "Bearer ", "", 1)
	if idToken == "" {
		log.Println("fail: No token provided")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := c.AuthClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		log.Printf("fail: Invalid token, %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	firebaseUID := token.UID

	var reqBody model.UserReqForHTTPPost
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("fail: json.NewDecoder, %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Usecase実行
	resUser, err := c.Usecase.RegisterUser(reqBody, firebaseUID)
	if err != nil {
		if err.Error() == "name is empty" || strings.Contains(err.Error(), "too long") {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err.Error())
			return
		}
		log.Printf("fail: RegisterUser, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resUser)
}

package controller

import (
	"encoding/json"
	"hackathon-backend/model"
	"hackathon-backend/usecase"
	"net/http"

	"firebase.google.com/go/auth"
)

type UserUpdateController struct {
	BaseController
	Usecase *usecase.UserUpdateUsecase
}

func NewUserUpdateController(u *usecase.UserUpdateUsecase, auth *auth.Client) *UserUpdateController {
	return &UserUpdateController{
		BaseController: BaseController{AuthClient: auth},
		Usecase:        u,
	}
}

func (c *UserUpdateController) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	var req model.UpdateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondError(w, http.StatusBadRequest, err)
		return
	}

	if req.Name == "" {
		c.respondError(w, http.StatusBadRequest, nil) // 名前は必須
		return
	}

	user, err := c.Usecase.UpdateUser(firebaseUID, req.Name, req.Bio)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	c.respondJSON(w, http.StatusOK, user)
}

package controller

import (
	"encoding/json"
	"hackathon-backend/model"
	"hackathon-backend/usecase"
	"net/http"

	"firebase.google.com/go/auth"
)

type ProductRegisterController struct {
	BaseController
	Usecase *usecase.ProductRegisterUsecase
}

func NewProductRegisterController(u *usecase.ProductRegisterUsecase, auth *auth.Client) *ProductRegisterController {
	return &ProductRegisterController{
		BaseController: BaseController{AuthClient: auth},
		Usecase:        u,
	}
}

func (c *ProductRegisterController) Handler(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	var reqBody model.ProductReq
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		c.respondError(w, http.StatusBadRequest, err)
		return
	}

	//  Usecase の実行
	product, err := c.Usecase.RegisterProduct(
		firebaseUID,
		reqBody.Name,
		reqBody.Description,
		reqBody.Price,
	)

	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	c.respondJSON(w, http.StatusCreated, product)
}

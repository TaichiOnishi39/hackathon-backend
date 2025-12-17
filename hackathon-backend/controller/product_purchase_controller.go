package controller

import (
	"hackathon-backend/usecase"
	"net/http"

	"firebase.google.com/go/auth"
)

type ProductPurchaseController struct {
	BaseController
	Usecase *usecase.ProductPurchaseUsecase
}

func NewProductPurchaseController(u *usecase.ProductPurchaseUsecase, auth *auth.Client) *ProductPurchaseController {
	return &ProductPurchaseController{BaseController: BaseController{AuthClient: auth}, Usecase: u}
}

func (c *ProductPurchaseController) HandlePurchaseProduct(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}
	productID := r.PathValue("id")

	err = c.Usecase.PurchaseProduct(productID, firebaseUID)
	if err != nil {
		c.respondError(w, http.StatusBadRequest, err)
		return
	}
	c.respondJSON(w, http.StatusOK, map[string]string{"status": "purchased"})
}

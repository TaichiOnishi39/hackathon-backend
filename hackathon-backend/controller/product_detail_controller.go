package controller

import (
	"hackathon-backend/usecase"
	"net/http"

	"firebase.google.com/go/auth"
)

type ProductDetailController struct {
	BaseController
	Usecase *usecase.ProductDetailUsecase
}

func NewProductDetailController(u *usecase.ProductDetailUsecase, auth *auth.Client) *ProductDetailController {
	return &ProductDetailController{BaseController: BaseController{AuthClient: auth}, Usecase: u}
}

func (c *ProductDetailController) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	product, err := c.Usecase.GetProductByID(productID)
	if err != nil {
		c.respondError(w, http.StatusNotFound, err)
		return
	}
	c.respondJSON(w, http.StatusOK, product)
}

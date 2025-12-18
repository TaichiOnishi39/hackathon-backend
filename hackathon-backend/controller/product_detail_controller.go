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
	// ★追加: 閲覧者IDを取得（未ログインなら空文字）
	viewerID := ""
	if uid, err := c.verifyToken(r); err == nil {
		viewerID = uid
	}

	productID := r.PathValue("id")

	// ★引数に viewerID を追加
	product, err := c.Usecase.GetProductByID(productID, viewerID)
	if err != nil {
		c.respondError(w, http.StatusNotFound, err)
		return
	}
	c.respondJSON(w, http.StatusOK, product)
}

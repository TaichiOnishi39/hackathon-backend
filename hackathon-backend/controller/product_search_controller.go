package controller

import (
	"hackathon-backend/usecase"
	"net/http"

	"firebase.google.com/go/auth"
)

type ProductSearchController struct {
	BaseController
	Usecase *usecase.ProductSearchUsecase
}

func NewProductSearchController(u *usecase.ProductSearchUsecase, auth *auth.Client) *ProductSearchController {
	return &ProductSearchController{
		BaseController: BaseController{AuthClient: auth},
		Usecase:        u,
	}
}

// HandleListProducts が GET /products の処理です
func (c *ProductSearchController) HandleListProducts(w http.ResponseWriter, r *http.Request) {
	// Usecase から商品一覧を取得
	products, err := c.Usecase.GetAllProducts()
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	// 取得したリストをJSONで返す
	c.respondJSON(w, http.StatusOK, products)
}

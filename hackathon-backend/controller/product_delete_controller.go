package controller

import (
	"fmt"
	"net/http"

	"hackathon-backend/usecase"

	"firebase.google.com/go/auth"
)

type ProductDeleteController struct {
	BaseController
	Usecase *usecase.ProductDeleteUsecase
}

func NewProductDeleteController(u *usecase.ProductDeleteUsecase, auth *auth.Client) *ProductDeleteController {
	return &ProductDeleteController{
		BaseController: BaseController{AuthClient: auth},
		Usecase:        u,
	}
}

// HandleDeleteProduct は DELETE /products?id=xxx を処理します
func (c *ProductDeleteController) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	// 1. 認証チェック
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	// 2. 削除したい商品IDを取得 (?id=xxx)
	productID := r.URL.Query().Get("id")
	if productID == "" {
		c.respondError(w, http.StatusBadRequest, fmt.Errorf("product id is required"))
		return
	}

	// 3. 削除実行
	err = c.Usecase.DeleteProduct(productID, firebaseUID)
	if err != nil {
		// 権限がない、または商品がない場合
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	// 4. 成功レスポンス
	c.respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

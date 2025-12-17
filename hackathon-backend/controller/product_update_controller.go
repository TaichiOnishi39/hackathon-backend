package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"hackathon-backend/model"
	"hackathon-backend/usecase"

	"firebase.google.com/go/auth"
)

type ProductUpdateController struct {
	BaseController
	Usecase *usecase.ProductUpdateUsecase
}

func NewProductUpdateController(u *usecase.ProductUpdateUsecase, auth *auth.Client) *ProductUpdateController {
	return &ProductUpdateController{
		BaseController: BaseController{AuthClient: auth},
		Usecase:        u,
	}
}

// HandleUpdateProduct は PUT /products?id=xxx を処理します
func (c *ProductUpdateController) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	// 1. 認証チェック
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	// 2. 更新したい商品IDを取得
	productID := r.URL.Query().Get("id")
	if productID == "" {
		c.respondError(w, http.StatusBadRequest, fmt.Errorf("product id is required"))
		return
	}

	// 3. リクエストボディの解析
	var req model.ProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondError(w, http.StatusBadRequest, err)
		return
	}

	// 4. 更新実行
	updatedProduct, err := c.Usecase.UpdateProduct(
		productID,
		firebaseUID,
		req.Name,
		req.Description,
		req.Price,
	)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	// 5. 成功レスポンス
	c.respondJSON(w, http.StatusOK, updatedProduct)
}

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
	keyword := r.URL.Query().Get("q")
	// Usecase から商品一覧を取得
	products, err := c.Usecase.SearchProduct(keyword)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	// 取得したリストをJSONで返す
	c.respondJSON(w, http.StatusOK, products)
}

// GET /users/{id}/products (公開ユーザーページ用)
func (c *ProductSearchController) HandleGetByUserID(w http.ResponseWriter, r *http.Request) {
	// URLのパスパラメータからIDを取得
	userID := r.PathValue("id")

	products, err := c.Usecase.GetProductsByUserID(userID)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}
	c.respondJSON(w, http.StatusOK, products)
}

// HandleGetSelling: GET /users/me/products (出品履歴)
func (c *ProductSearchController) HandleGetSelling(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}
	products, err := c.Usecase.GetSellingProducts(firebaseUID)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}
	c.respondJSON(w, http.StatusOK, products)
}

// HandleGetPurchased: GET /users/me/purchases (購入履歴)
func (c *ProductSearchController) HandleGetPurchased(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}
	products, err := c.Usecase.GetPurchasedProducts(firebaseUID)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}
	c.respondJSON(w, http.StatusOK, products)
}

// HandleGetLiked: GET /users/me/likes (いいね一覧)
func (c *ProductSearchController) HandleGetLiked(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}
	products, err := c.Usecase.GetLikedProducts(firebaseUID)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}
	c.respondJSON(w, http.StatusOK, products)
}

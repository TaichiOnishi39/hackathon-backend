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
	// ★追加: ログインしていれば閲覧者IDを取得（未ログインなら空文字）
	viewerID := ""
	if uid, err := c.verifyToken(r); err == nil {
		viewerID = uid
	}

	keyword := r.URL.Query().Get("q")
	sortOrder := r.URL.Query().Get("sort")
	status := r.URL.Query().Get("status")

	// ★引数に viewerID を追加して呼び出し
	products, err := c.Usecase.SearchProduct(keyword, sortOrder, status, viewerID)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	c.respondJSON(w, http.StatusOK, products)
}

// GET /users/{id}/products (公開ユーザーページ用)
func (c *ProductSearchController) HandleGetByUserID(w http.ResponseWriter, r *http.Request) {
	// ★追加: 閲覧者IDを取得
	viewerID := ""
	if uid, err := c.verifyToken(r); err == nil {
		viewerID = uid
	}

	// URLのパスパラメータからIDを取得
	userID := r.PathValue("id")
	sortOrder := r.URL.Query().Get("sort")
	status := r.URL.Query().Get("status")

	// ★引数に viewerID を追加
	products, err := c.Usecase.GetProductsByUserID(userID, sortOrder, status, viewerID)
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

	// ★自分が出品者(target)であり、閲覧者(viewer)でもある
	products, err := c.Usecase.GetSellingProducts(firebaseUID, firebaseUID)
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

	// ★自分が購入者(target)であり、閲覧者(viewer)でもある
	products, err := c.Usecase.GetPurchasedProducts(firebaseUID, firebaseUID)
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

	// ★自分がいいねした人(target)であり、閲覧者(viewer)でもある
	products, err := c.Usecase.GetLikedProducts(firebaseUID, firebaseUID)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}
	c.respondJSON(w, http.StatusOK, products)
}

package controller

import (
	"hackathon-backend/usecase"
	"net/http"

	"firebase.google.com/go/auth"
)

type ProductLikeController struct {
	BaseController
	Usecase *usecase.ProductLikeUsecase
}

func NewProductLikeController(u *usecase.ProductLikeUsecase, auth *auth.Client) *ProductLikeController {
	return &ProductLikeController{
		BaseController: BaseController{AuthClient: auth},
		Usecase:        u,
	}
}

// HandleToggleLike: POST /products/{id}/like
func (c *ProductLikeController) HandleToggleLike(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	productID := r.PathValue("id")
	// 切り替え実行
	liked, err := c.Usecase.ToggleLike(productID, firebaseUID)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	// 結果(true/false)を返す
	c.respondJSON(w, http.StatusOK, map[string]bool{"liked": liked})
}

// HandleGetLikeStatus: GET /products/{id}/like
func (c *ProductLikeController) HandleGetLikeStatus(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		// ログインしていなければ "liked": false で返す手もあるが、今回は401
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	productID := r.PathValue("id")
	liked, err := c.Usecase.GetLikeStatus(productID, firebaseUID)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	c.respondJSON(w, http.StatusOK, map[string]bool{"liked": liked})
}

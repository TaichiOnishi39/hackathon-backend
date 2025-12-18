package controller

import (
	"encoding/json"
	"hackathon-backend/model"
	"hackathon-backend/usecase"
	"net/http"

	"firebase.google.com/go/auth"
)

type ProductDescriptionController struct {
	BaseController
	Usecase *usecase.ProductDescriptionUsecase
}

func NewProductDescriptionController(u *usecase.ProductDescriptionUsecase, auth *auth.Client) *ProductDescriptionController {
	return &ProductDescriptionController{
		BaseController: BaseController{AuthClient: auth},
		Usecase:        u,
	}
}

// HandleGenerate: POST /products/generate-description
func (c *ProductDescriptionController) HandleGenerate(w http.ResponseWriter, r *http.Request) {
	// ログインチェック
	_, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	var req model.GenerateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondError(w, http.StatusBadRequest, err)
		return
	}

	if req.Name == "" {
		c.respondError(w, http.StatusBadRequest, nil)
		return
	}

	// 生成実行
	desc, err := c.Usecase.Generate(r.Context(), req.Name, req.Keywords)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	c.respondJSON(w, http.StatusOK, model.GenerateRes{Description: desc})
}

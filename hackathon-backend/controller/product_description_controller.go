package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hackathon-backend/model"
	"hackathon-backend/usecase"
	"io"
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

// ★追加: 画像アップロード解析ハンドラ
func (c *ProductDescriptionController) HandleGenerateFromImage(w http.ResponseWriter, r *http.Request) {
	_, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	// 1. 画像ファイルの取得 (form key: "image")
	file, header, err := r.FormFile("image")
	if err != nil {
		c.respondError(w, http.StatusBadRequest, fmt.Errorf("failed to get image file: %w", err))
		return
	}
	defer file.Close()

	// 2. 画像データをバイト列に読み込む
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	// 3. MIMEタイプの取得
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg" // デフォルト
	}

	// 4. Usecase実行
	res, err := c.Usecase.GenerateInfoFromImage(r.Context(), buf.Bytes(), mimeType)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	// 5. 結果を返す
	c.respondJSON(w, http.StatusOK, res)
}

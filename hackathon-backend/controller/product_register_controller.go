package controller

import (
	"fmt"
	"hackathon-backend/usecase"
	"net/http"
	"strconv"

	"firebase.google.com/go/auth"
)

type ProductRegisterController struct {
	BaseController
	Usecase *usecase.ProductRegisterUsecase
}

func NewProductRegisterController(u *usecase.ProductRegisterUsecase, auth *auth.Client) *ProductRegisterController {
	return &ProductRegisterController{
		BaseController: BaseController{AuthClient: auth},
		Usecase:        u,
	}
}

func (c *ProductRegisterController) Handler(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	//  multipart/form-data の解析 (最大10MBまでメモリ展開)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		c.respondError(w, http.StatusBadRequest, err)
		return
	}

	//  フォームデータの取得
	name := r.FormValue("name")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	// 価格を数値に変換
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		c.respondError(w, http.StatusBadRequest, err)
		return
	}

	//  画像ファイルの取得
	file, header, err := r.FormFile("image") // フロント側で "image" というキーで送る

	// ★変更: 画像がない場合はエラーにする
	if err != nil {
		c.respondError(w, http.StatusBadRequest, fmt.Errorf("image is required"))
		return
	}
	defer file.Close()
	//  Usecase 実行
	product, err := c.Usecase.RegisterProduct(
		firebaseUID,
		name,
		description,
		price,
		file,
		header.Filename,
	)

	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	c.respondJSON(w, http.StatusCreated, product)
}

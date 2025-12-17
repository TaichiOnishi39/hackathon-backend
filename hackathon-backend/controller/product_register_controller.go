package controller

import (
	"hackathon-backend/usecase"
	"io"
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
	// 画像がない場合も許容するならエラーハンドリングを調整
	// 今回は「画像なしでもOK」とする場合:
	var fileReader io.Reader
	var fileName string
	if err == nil {
		defer file.Close()
		fileReader = file
		fileName = header.Filename
	}

	//  Usecase 実行
	product, err := c.Usecase.RegisterProduct(
		firebaseUID,
		name,
		description,
		price,
		fileReader, // 画像データ
		fileName,   // ファイル名
	)

	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	c.respondJSON(w, http.StatusCreated, product)
}

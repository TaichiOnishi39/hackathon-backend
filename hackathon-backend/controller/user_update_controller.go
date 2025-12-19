package controller

import (
	"hackathon-backend/usecase"
	"net/http"

	"firebase.google.com/go/auth"
)

type UserUpdateController struct {
	BaseController
	Usecase *usecase.UserUpdateUsecase
}

func NewUserUpdateController(u *usecase.UserUpdateUsecase, auth *auth.Client) *UserUpdateController {
	return &UserUpdateController{
		BaseController: BaseController{AuthClient: auth},
		Usecase:        u,
	}
}

func (c *UserUpdateController) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	// ★ JSONデコードではなく MultipartForm のパースに変更
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
		c.respondError(w, http.StatusBadRequest, err)
		return
	}

	// フォーム値の取得
	name := r.FormValue("name")
	bio := r.FormValue("bio")

	if name == "" {
		c.respondError(w, http.StatusBadRequest, nil) // 名前は必須
		return
	}

	// ★ 画像ファイルの取得
	file, header, err := r.FormFile("image")
	// ファイルがない場合は err が返るが、画像なし更新も許可したいのでチェック
	if err == nil {
		defer file.Close()
		// 画像がある場合は header を渡す
	} else if err == http.ErrMissingFile {
		// 画像なしの場合は nil を渡す
		header = nil
	} else {
		// その他のエラー
		c.respondError(w, http.StatusBadRequest, err)
		return
	}

	// UseCase 呼び出し
	user, err := c.Usecase.UpdateUser(r.Context(), firebaseUID, name, bio, header)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	c.respondJSON(w, http.StatusOK, user)
}

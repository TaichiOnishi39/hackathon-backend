package controller

import (
	"fmt"
	"hackathon-backend/usecase"
	"log"
	"net/http"

	"firebase.google.com/go/auth"
)

type SearchUserController struct {
	BaseController
	Usecase *usecase.SearchUserUsecase
}

func NewSearchUserController(u *usecase.SearchUserUsecase, auth *auth.Client) *SearchUserController {
	return &SearchUserController{
		BaseController: BaseController{AuthClient: auth},
		Usecase:        u,
	}
}

func (c *SearchUserController) HandleSearch(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		log.Println("fail: name is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := c.Usecase.SearchUser(name)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	c.respondJSON(w, http.StatusOK, users)
}

func (c *SearchUserController) HandleGetMe(w http.ResponseWriter, r *http.Request) {
	firebaseUID, err := c.verifyToken(r)
	if err != nil {
		c.respondError(w, http.StatusUnauthorized, err)
		return
	}

	user, err := c.Usecase.GetUserByFirebaseUID(firebaseUID)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		// 404 Not Found を返す (これでフロントエンドが「未登録だ」と気づける)
		c.respondError(w, http.StatusNotFound, fmt.Errorf("user not found"))
		return
	}

	c.respondJSON(w, http.StatusOK, user)
}

// ★追加: GET /users/{id}
func (c *SearchUserController) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	// URLパラメータからIDを取得 (例: /users/01HXYZ...)
	userID := r.PathValue("id")

	user, err := c.Usecase.GetUserByID(userID)
	if err != nil {
		c.respondError(w, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		c.respondError(w, http.StatusNotFound, nil) // ユーザーがいない場合
		return
	}

	// 公開して良い情報だけ返すのが理想ですが、現状のUserモデル(名前, Bio, ID)ならそのまま返してOK
	c.respondJSON(w, http.StatusOK, user)
}

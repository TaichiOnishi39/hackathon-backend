package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"hackathon-backend/usecase"

	"firebase.google.com/go/auth"
)

type SearchUserController struct {
	Usecase    *usecase.SearchUserUsecase
	AuthClient *auth.Client
}

func NewSearchUserController(u *usecase.SearchUserUsecase, auth *auth.Client) *SearchUserController {
	return &SearchUserController{Usecase: u, AuthClient: auth}
}

func (c *SearchUserController) Handle(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		log.Println("fail: name is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := c.Usecase.SearchUser(name)
	if err != nil {
		log.Printf("fail: SearchUser, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(users)
	if err != nil {
		log.Printf("fail: json.Marshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func (c *SearchUserController) GetMe(w http.ResponseWriter, r *http.Request) {
	// 1. Authorizationヘッダー確認
	authHeader := r.Header.Get("Authorization")
	idToken := strings.Replace(authHeader, "Bearer ", "", 1)
	if idToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 2. トークン検証
	token, err := c.AuthClient.VerifyIDToken(r.Context(), idToken)
	if err != nil {
		log.Printf("fail: VerifyIDToken, %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// ★★★ デバッグログその1: 検証後のUID確認 ★★★
	firebaseUID := token.UID
	log.Printf("DEBUG: Verified UID: %s", firebaseUID)

	// 3. Usecase経由でDB検索 (token.UIDを使う)
	user, err := c.Usecase.GetUserByFirebaseUID(token.UID)
	if err != nil {
		// ★★★ デバッグログその2: DBエラー確認 ★★★
		log.Printf("ERROR: GetUserByFirebaseUID failed for UID %s: %v\n", token.UID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 4. ユーザーがいなかった場合 (404)
	if user == nil {
		// ★★★ デバッグログその3: ユーザー未登録として扱われている ★★★
		log.Printf("INFO: User not found in DB for UID: %s (Returning 404)", token.UID)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// 5. JSONを返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("fail: json.Encode, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

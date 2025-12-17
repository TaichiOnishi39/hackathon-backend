package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
)

// BaseController は共通機能を提供します
// 他のコントローラーはこの構造体を埋め込むか、この関数を呼び出して使います
type BaseController struct {
	AuthClient *auth.Client
}

// verifyToken: AuthorizationヘッダーからUIDを取得する共通関数
func (b *BaseController) verifyToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	idToken := strings.Replace(authHeader, "Bearer ", "", 1)
	if idToken == "" {
		return "", fmt.Errorf("no token provided")
	}

	token, err := b.AuthClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return "", err
	}
	return token.UID, nil
}

// respondJSON: JSONレスポンスを返す共通関数
func (b *BaseController) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			fmt.Printf("fail: json encode response, %v\n", err)
		}
	}
}

// respondError: エラーレスポンスを返す共通関数
func (b *BaseController) respondError(w http.ResponseWriter, status int, err error) {
	// ログ出しなどをここに入れると便利
	fmt.Printf("Error: %v\n", err)
	http.Error(w, err.Error(), status)
}

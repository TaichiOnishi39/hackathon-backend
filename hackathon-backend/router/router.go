package router

import (
	"hackathon-backend/controller"
	"net/http"
)

// NewRouter は全コントローラーを受け取り、ルーティングを設定して返します
func NewRouter(
	registerUserCtrl *controller.RegisterUserController,
	searchUserCtrl *controller.SearchUserController,
	productRegisterCtrl *controller.ProductRegisterController,
	productSearchCtrl *controller.ProductSearchController,
) http.Handler {
	mux := http.NewServeMux()

	// --- ルーティング定義 ---

	// /users (GET: Search, POST: Register)
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		// CORS設定 (共通化も可能ですが、まずはここに)
		if !enableCORS(w, r) {
			return
		}

		switch r.Method {
		case http.MethodGet:
			searchUserCtrl.HandleSearch(w, r)
		case http.MethodPost:
			registerUserCtrl.Handle(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// /users/me
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		if r.Method == http.MethodGet {
			searchUserCtrl.HandleGetMe(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// /products
	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		switch r.Method {
		case http.MethodPost:
			productRegisterCtrl.Handler(w, r)
		case http.MethodGet:
			productSearchCtrl.HandleListProducts(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return mux
}

// enableCORS: CORSヘッダーをセットし、OPTIONSリクエストならfalseを返して終了させる
func enableCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// OPTIONSリクエスト（プリフライト）の場合はここで処理を終える
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return false
	}
	return true
}

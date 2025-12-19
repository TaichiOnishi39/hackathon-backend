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
	productDeleteCtrl *controller.ProductDeleteController,
	productUpdateCtrl *controller.ProductUpdateController,
	productDetailCtrl *controller.ProductDetailController,
	productPurchaseCtrl *controller.ProductPurchaseController,
	messageCtrl *controller.MessageController,
	productLikeCtrl *controller.ProductLikeController,
	userUpdateCtrl *controller.UserUpdateController,
	productDescCtrl *controller.ProductDescriptionController,
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
		switch r.Method {
		case http.MethodGet:
			searchUserCtrl.HandleGetMe(w, r)
		case http.MethodPut: // ★追加: 更新はPUT
			userUpdateCtrl.HandleUpdate(w, r)
		default:
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
		case http.MethodDelete:
			productDeleteCtrl.HandleDeleteProduct(w, r)
		case http.MethodPut:
			productUpdateCtrl.HandleUpdateProduct(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// /products/{id}
	mux.HandleFunc("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		if r.Method == http.MethodGet {
			productDetailCtrl.HandleGetProduct(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/{id}/products", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		if r.Method == http.MethodGet {
			productSearchCtrl.HandleGetByUserID(w, r)
		}
	})

	mux.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		if r.Method == http.MethodGet {
			searchUserCtrl.HandleGetUser(w, r)
		}
	})

	// /products/{id}/purchase
	mux.HandleFunc("/products/{id}/purchase", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		if r.Method == http.MethodPost {
			productPurchaseCtrl.HandlePurchaseProduct(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// /messages
	mux.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		switch r.Method {
		case http.MethodPost:
			messageCtrl.HandleSendMessage(w, r)
		case http.MethodGet:
			messageCtrl.HandleGetChat(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// /messages/list
	mux.HandleFunc("/messages/list", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		if r.Method == http.MethodGet {
			messageCtrl.HandleGetChatList(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// /messages/read
	mux.HandleFunc("/messages/read", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		if r.Method == http.MethodPost {
			messageCtrl.HandleMarkAsRead(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// いいね機能
	// GET: 状態確認, POST: 切り替え
	mux.HandleFunc("/products/{id}/like", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		switch r.Method {
		case http.MethodGet:
			productLikeCtrl.HandleGetLikeStatus(w, r)
		case http.MethodPost:
			productLikeCtrl.HandleToggleLike(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// 出品した商品
	mux.HandleFunc(" /users/me/products", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		if r.Method == http.MethodGet {
			productSearchCtrl.HandleGetSelling(w, r)
		}
	})
	// 購入した商品
	mux.HandleFunc("/users/me/purchases", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		if r.Method == http.MethodGet {
			productSearchCtrl.HandleGetPurchased(w, r)
		}
	})
	// いいねした商品
	mux.HandleFunc("/users/me/likes", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		if r.Method == http.MethodGet {
			productSearchCtrl.HandleGetLiked(w, r)
		}
	})

	// ★AI生成エンドポイント
	mux.HandleFunc("/products/generate-description", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		if r.Method == http.MethodPost {
			productDescCtrl.HandleGenerate(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// AI生成 (画像から)
	mux.HandleFunc("/products/generate-from-image", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		if r.Method == http.MethodPost {
			productDescCtrl.HandleGenerateFromImage(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/messages/{id}/unsend", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		if r.Method == http.MethodPut {
			messageCtrl.HandleUnsendMessage(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/messages/{id}", func(w http.ResponseWriter, r *http.Request) {
		if !enableCORS(w, r) {
			return
		}
		if r.Method == http.MethodDelete {
			messageCtrl.HandleDeleteMessage(w, r)
		} else {
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

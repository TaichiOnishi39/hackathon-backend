package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	firebase "firebase.google.com/go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	// 各パッケージをインポート (module名は適宜書き換えてください)
	"hackathon-backend/controller"
	"hackathon-backend/dao"
	"hackathon-backend/usecase"
)

var db *sql.DB

func main() {
	// --- 1. DB接続設定 ---
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlUserPwd := os.Getenv("MYSQL_PASSWORD")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	mysqlHost := os.Getenv("MYSQL_HOST")
	dsn := fmt.Sprintf("%s:%s@%s/%s?parseTime=true", mysqlUser, mysqlUserPwd, mysqlHost, mysqlDatabase)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("fail: sql.Open, %v\n", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("fail: db.Ping, %v\n", err)
	}
	log.Println("Successfully connected to the database!")

	// --- Firebase初期化 (追加) ---
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "term8-taichi-onishi"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error initializing auth client: %v\n", err)
	}

	userDAO := dao.NewUserDao(db)
	registerUsecase := usecase.NewRegisterUserUsecase(userDAO)
	searchUsecase := usecase.NewSearchUserUsecase(userDAO)
	registerController := controller.NewRegisterUserController(registerUsecase, authClient)
	searchController := controller.NewSearchUserController(searchUsecase)
	// --- 3. ルーティング設定 ---
	// 単一のエンドポイントでメソッドによって振り分ける場合
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		switch r.Method {
		case http.MethodGet:
			searchController.Handle(w, r)
		case http.MethodPost:
			registerController.Handle(w, r)
		default:
			log.Printf("fail: HTTP Method is %s\n", r.Method)
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	// シャットダウン処理のセットアップ
	closeDBWithSysCall()

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // ローカルで動かすとき用
	}

	log.Printf("Listening on :%s...", port)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func closeDBWithSysCall() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		log.Printf("received syscall, %v", s)

		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
		log.Printf("success: db.Close()")
		os.Exit(0)
	}()
}

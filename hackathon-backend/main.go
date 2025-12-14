package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	// --- 2. 依存関係の注入 (DI) ---
	// DAOの初期化
	userDAO := dao.NewUserDao(db)

	// Usecaseの初期化
	searchUsecase := usecase.NewSearchUserUsecase(userDAO)
	registerUsecase := usecase.NewRegisterUserUsecase(userDAO)

	// Controllerの初期化
	searchController := controller.NewSearchUserController(searchUsecase)
	registerController := controller.NewRegisterUserController(registerUsecase)

	// --- 3. ルーティング設定 ---
	// 単一のエンドポイントでメソッドによって振り分ける場合
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
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

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

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"hackathon-backend/controller"
	"hackathon-backend/dao"
	"hackathon-backend/router"
	"hackathon-backend/service"
	"hackathon-backend/usecase"
)

var db *sql.DB

func main() {
	// --- DB接続設定 ---
	_ = godotenv.Load()

	db := initDB()
	defer db.Close()
	// --- Firebase初期化 ---
	authClient := initFirebase()

	// --- GCS初期化  ---
	ctx := context.Background()
	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("fail: storage.NewClient, %v", err)
	}
	defer gcsClient.Close()

	bucketName := os.Getenv("GCS_BUCKET_NAME")
	storageService := service.NewStorageService(gcsClient, bucketName)
	
	projectID := "term8-taichi-onishi"
	location := "asia-northeast1"
	modelName := "gemini-2.5-flash"
	geminiService, err := service.NewGeminiService(ctx, projectID, location, modelName)
	if err != nil {
		log.Fatalf("failed to init gemini: %v", err)
	}
	defer geminiService.Close()

	//DAO
	userDAO := dao.NewUserDao(db)
	productDAO := dao.NewProductDAO(db)
	messageDAO := dao.NewMessageDao(db)
	likeDAO := dao.NewLikeDao(db)

	//Usecase
	registerUsecase := usecase.NewRegisterUserUsecase(userDAO)
	searchUsecase := usecase.NewSearchUserUsecase(userDAO)
	productRegisterUsecase := usecase.NewProductRegisterUsecase(productDAO, userDAO, storageService)
	productSearchUsecase := usecase.NewProductSearchUsecase(productDAO, userDAO, storageService)
	productDeleteUsecase := usecase.NewProductDeleteUsecase(productDAO, userDAO)
	productUpdateUsecase := usecase.NewProductUpdateUsecase(productDAO, userDAO)
	productDetailUsecase := usecase.NewProductDetailUsecase(productDAO, storageService)
	productPurchaseUsecase := usecase.NewProductPurchaseUsecase(productDAO, userDAO)
	messageUsecase := usecase.NewMessageUsecase(messageDAO, userDAO)
	productLikeUsecase := usecase.NewProductLikeUsecase(likeDAO, userDAO)
	userUpdateUsecase := usecase.NewUserUpdateUsecase(userDAO)
	productDescUsecase := usecase.NewProductDescriptionUsecase(geminiService)

	//Controller
	registerUserCtrl := controller.NewRegisterUserController(registerUsecase, authClient)
	searchUserCtrl := controller.NewSearchUserController(searchUsecase, authClient)
	productRegisterCtrl := controller.NewProductRegisterController(productRegisterUsecase, authClient)
	productSearchCtrl := controller.NewProductSearchController(productSearchUsecase, authClient)
	productDeleteCtrl := controller.NewProductDeleteController(productDeleteUsecase, authClient)
	productUpdateCtrl := controller.NewProductUpdateController(productUpdateUsecase, authClient)
	productDetailCtrl := controller.NewProductDetailController(productDetailUsecase, authClient)
	productPurchaseCtrl := controller.NewProductPurchaseController(productPurchaseUsecase, authClient)
	messageCtrl := controller.NewMessageController(messageUsecase, authClient)
	productLikeCtrl := controller.NewProductLikeController(productLikeUsecase, authClient)
	userUpdateCtrl := controller.NewUserUpdateController(userUpdateUsecase, authClient)
	productDescCtrl := controller.NewProductDescriptionController(productDescUsecase, authClient)

	// --- 3. ルーティング設定 ---
	mux := router.NewRouter(
		registerUserCtrl,
		searchUserCtrl,
		productRegisterCtrl,
		productSearchCtrl,
		productDeleteCtrl,
		productUpdateCtrl,
		productDetailCtrl,
		productPurchaseCtrl,
		messageCtrl,
		productLikeCtrl,
		userUpdateCtrl,
		productDescCtrl,
	)

	// シャットダウン処理のセットアップ
	closeDBWithSysCall()

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // ローカルで動かすとき用
	}

	log.Printf("Listening on :%s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func initDB() *sql.DB {
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlUserPwd := os.Getenv("MYSQL_PASSWORD")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	mysqlHost := os.Getenv("MYSQL_HOST")
	dsn := fmt.Sprintf("%s:%s@%s/%s?parseTime=true", mysqlUser, mysqlUserPwd, mysqlHost, mysqlDatabase)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("fail: sql.Open, %v\n", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("fail: db.Ping, %v\n", err)
	}
	log.Println("Successfully connected to the database!")
	return db
}

func initFirebase() *auth.Client {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "term8-taichi-onishi"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error initializing auth client: %v", err)
	}
	return client
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

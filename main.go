package main

import (
	"fmt"
	"log"
	"os"

	"voice-link/domain/model"
	"voice-link/infrastructure/persistence"
	"voice-link/interface/handler/auth"
	"voice-link/interface/handler/user"
	"voice-link/interface/router"
	"voice-link/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// JWT_SECRETの設定
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "your-secret-key-change-in-production")
	}

	// データベース接続
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// マイグレーション
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 依存関係の注入
	userRepo := persistence.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	authHandler := auth.NewAuthHandler(userUseCase)
	userHandler := user.NewUserHandler(userUseCase)

	// Echoのインスタンスを作成
	e := echo.New()

	// ルーティングの設定
	r := router.NewRouter(e, authHandler, userHandler)
	r.Setup()

	// サーバーの起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is starting on port %s", port)
	if err := e.Start(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

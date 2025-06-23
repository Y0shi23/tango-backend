package main

import (
	"log"

	"backend/auth"
	"backend/config"
	"backend/database"
	"backend/routes"
)

func main() {
	// 設定を読み込み
	cfg := config.Load()

	// JWT秘密鍵を初期化
	auth.InitJWT(cfg.JWT.Secret)

	// データベース接続
	if err := database.Connect(cfg.Database); err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer database.Close()

	log.Println("Database tables initialized")

	// ルーターを設定
	r := routes.SetupRouter(cfg)

	// サーバー起動
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

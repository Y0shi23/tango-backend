package main

import (
	"log"
	"net/http"

	"backend/auth"
	"backend/config"
	"backend/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	// Ginルーターを設定
	r := gin.Default()

	// CORS設定
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.Server.AllowOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// ルートを定義
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Backend API is running!",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		// データベース接続状態もチェック
		dbStatus := "connected"
		if err := database.ValidateConnection(); err != nil {
			dbStatus = "disconnected"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"database": dbStatus,
		})
	})

	// 認証関連のルート
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", auth.RegisterHandler)
		authGroup.POST("/login", auth.LoginHandler)
	}

	// API v1 ルートグループ
	v1 := r.Group("/api/v1")
	{
		v1.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "API v1 test endpoint",
			})
		})

		// 認証が必要なルート
		protected := v1.Group("/")
		protected.Use(auth.Middleware())
		{
			protected.GET("/profile", auth.ProfileHandler)
		}
	}

	// サーバー起動
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

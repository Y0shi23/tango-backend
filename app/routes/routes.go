package routes

import (
	"backend/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter はGinルーターを設定し、全てのルートを登録します
func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// CORS設定
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.Server.AllowOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// 基本ルートを設定
	setupBaseRoutes(r)

	// 認証ルートを設定
	setupAuthRoutes(r)

	// API v1 ルートを設定
	setupAPIRoutes(r)

	return r
}

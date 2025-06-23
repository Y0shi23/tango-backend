package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

// setupBaseRoutes は基本的なルート（ヘルスチェック等）を設定します
func setupBaseRoutes(r *gin.Engine) {
	// ルートエンドポイント
	r.GET("/", handlers.RootHandler)

	// ヘルスチェックエンドポイント
	r.GET("/health", handlers.HealthHandler)
}

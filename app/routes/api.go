package routes

import (
	"backend/auth"
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

// setupAPIRoutes はAPI v1関連のルートを設定します
func setupAPIRoutes(r *gin.Engine) {
	// API v1 ルートグループ
	v1 := r.Group("/api/v1")
	{
		v1.GET("/test", handlers.TestHandler)

		// 認証が必要なルート
		protected := v1.Group("/")
		protected.Use(auth.Middleware())
		{
			protected.GET("/profile", auth.ProfileHandler)
			protected.PUT("/profile", auth.UpdateProfileHandler)
		}
	}
}

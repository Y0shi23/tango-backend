package routes

import (
	"backend/auth"

	"github.com/gin-gonic/gin"
)

// setupAuthRoutes は認証関連のルートを設定します
func setupAuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", auth.RegisterHandler)
		authGroup.POST("/login", auth.LoginHandler)
	}
}

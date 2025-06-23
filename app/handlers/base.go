package handlers

import (
	"net/http"

	"backend/database"

	"github.com/gin-gonic/gin"
)

// RootHandler はルートエンドポイントのハンドラーです
func RootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Backend API is running!",
	})
}

// HealthHandler はヘルスチェックエンドポイントのハンドラーです
func HealthHandler(c *gin.Context) {
	// データベース接続状態もチェック
	dbStatus := "connected"
	if err := database.ValidateConnection(); err != nil {
		dbStatus = "disconnected"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "healthy",
		"database": dbStatus,
	})
}

// TestHandler はAPI v1テストエンドポイントのハンドラーです
func TestHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "API v1 test endpoint",
	})
}

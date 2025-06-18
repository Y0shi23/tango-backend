package auth

import (
	"net/http"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
)

// RegisterHandler はユーザー登録ハンドラーです
func RegisterHandler(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザー名の重複チェック
	var existingUser models.User
	err := database.GetDB().QueryRow("SELECT id FROM users WHERE username = $1 OR email = $2", req.Username, req.Email).Scan(&existingUser.ID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	// パスワードをハッシュ化
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// ユーザーをデータベースに挿入
	var userID int
	err = database.GetDB().QueryRow(
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
		req.Username, req.Email, hashedPassword,
	).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// JWTトークンを生成
	token, err := GenerateJWT(userID, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// レスポンスを返す
	user := models.User{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
	}

	response := models.AuthResponse{
		Message: "User created successfully",
		Token:   token,
		User:    user.ToResponse(),
	}

	c.JSON(http.StatusCreated, response)
}

// LoginHandler はログインハンドラーです
func LoginHandler(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザーをデータベースから取得
	var user models.User
	err := database.GetDB().QueryRow(
		"SELECT id, username, email, password FROM users WHERE username = $1",
		req.Username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// パスワードを検証
	if !CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// JWTトークンを生成
	token, err := GenerateJWT(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// レスポンスを返す
	response := models.AuthResponse{
		Message: "Login successful",
		Token:   token,
		User:    user.ToResponse(),
	}

	c.JSON(http.StatusOK, response)
}

// ProfileHandler はプロフィール取得ハンドラーです（認証が必要）
func ProfileHandler(c *gin.Context) {
	userID := c.GetInt("user_id")

	var user models.User
	err := database.GetDB().QueryRow(
		"SELECT id, username, email FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user.ToResponse(),
	})
}

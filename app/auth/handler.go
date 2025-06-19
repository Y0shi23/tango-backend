package auth

import (
	"net/http"
	"time"

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

	// ユーザー名またはメールの重複チェック
	var existingUser models.User
	err := database.GetDB().Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error
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

	// デフォルト値を設定
	preferredAccent := req.PreferredAccent
	if preferredAccent == "" {
		preferredAccent = "US"
	}

	studyLevel := req.StudyLevel
	if studyLevel == "" {
		studyLevel = "BEGINNER"
	}

	// ユーザーをデータベースに挿入
	user := models.User{
		Username:        req.Username,
		Email:           req.Email,
		PasswordHash:    hashedPassword,
		PreferredAccent: preferredAccent,
		StudyLevel:      studyLevel,
		CreatedAt:       time.Now(),
	}

	if err := database.GetDB().Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
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
	err := database.GetDB().Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// パスワードを検証
	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// last_loginを更新
	now := time.Now()
	user.LastLogin = &now
	database.GetDB().Model(&user).Update("last_login", now)

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
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	var user models.User
	err := database.GetDB().Where("user_id = ?", userIDUint).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user.ToResponse(),
	})
}

// UpdateProfileRequest はプロフィール更新リクエスト構造体です
type UpdateProfileRequest struct {
	PreferredAccent string `json:"preferred_accent,omitempty"`
	StudyLevel      string `json:"study_level,omitempty"`
}

// UpdateProfileHandler はプロフィール更新ハンドラーです（認証が必要）
func UpdateProfileHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザーを取得
	var user models.User
	err := database.GetDB().Where("user_id = ?", userIDUint).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// フィールドを更新
	updateData := make(map[string]interface{})
	if req.PreferredAccent != "" {
		if req.PreferredAccent != "US" && req.PreferredAccent != "UK" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid preferred_accent. Must be US or UK"})
			return
		}
		updateData["preferred_accent"] = req.PreferredAccent
		user.PreferredAccent = req.PreferredAccent
	}
	if req.StudyLevel != "" {
		if req.StudyLevel != "BEGINNER" && req.StudyLevel != "INTERMEDIATE" && req.StudyLevel != "ADVANCED" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid study_level. Must be BEGINNER, INTERMEDIATE, or ADVANCED"})
			return
		}
		updateData["study_level"] = req.StudyLevel
		user.StudyLevel = req.StudyLevel
	}

	if len(updateData) > 0 {
		err = database.GetDB().Model(&user).Updates(updateData).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    user.ToResponse(),
	})
}

package models

import "time"

// User構造体
type User struct {
	ID              uint       `json:"id" gorm:"primary_key;column:user_id"`
	Username        string     `json:"username" gorm:"not null;unique;size:50"`
	Email           string     `json:"email" gorm:"not null;unique;size:100"`
	PasswordHash    string     `json:"-" gorm:"not null;size:255;column:password_hash"` // JSONには含めない
	PreferredAccent string     `json:"preferred_accent" gorm:"default:'US';size:10;check:preferred_accent in ('US', 'UK')"`
	StudyLevel      string     `json:"study_level" gorm:"default:'BEGINNER';size:20;check:study_level in ('BEGINNER', 'INTERMEDIATE', 'ADVANCED')"`
	CreatedAt       time.Time  `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	LastLogin       *time.Time `json:"last_login" gorm:"column:last_login"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// ログインリクエスト構造体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ユーザー登録リクエスト構造体
type RegisterRequest struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	PreferredAccent string `json:"preferred_accent,omitempty"`
	StudyLevel      string `json:"study_level,omitempty"`
}

// ユーザーレスポンス構造体（パスワードを除外）
type UserResponse struct {
	ID              uint       `json:"id"`
	Username        string     `json:"username"`
	Email           string     `json:"email"`
	PreferredAccent string     `json:"preferred_accent"`
	StudyLevel      string     `json:"study_level"`
	CreatedAt       time.Time  `json:"created_at"`
	LastLogin       *time.Time `json:"last_login,omitempty"`
}

// ログイン/登録レスポンス構造体
type AuthResponse struct {
	Message string       `json:"message"`
	Token   string       `json:"token"`
	User    UserResponse `json:"user"`
}

// ToResponse はUserをUserResponseに変換します
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:              u.ID,
		Username:        u.Username,
		Email:           u.Email,
		PreferredAccent: u.PreferredAccent,
		StudyLevel:      u.StudyLevel,
		CreatedAt:       u.CreatedAt,
		LastLogin:       u.LastLogin,
	}
}

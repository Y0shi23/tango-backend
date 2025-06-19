package models

// User構造体
type User struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	Username  string `json:"username" gorm:"not null;unique"`
	Email     string `json:"email" gorm:"not null;unique"`
	Password  string `json:"-" gorm:"not null"` // JSONには含めない
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

// ログインリクエスト構造体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ユーザー登録リクエスト構造体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// ユーザーレスポンス構造体（パスワードを除外）
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
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
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}

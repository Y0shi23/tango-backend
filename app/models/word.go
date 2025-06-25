package models

import "time"

// Word構造体 - 冗長なフィールドを削除
type Word struct {
	ID             uint      `json:"id" gorm:"primary_key;column:id"`
	Word           string    `json:"word" gorm:"column:word;size:100;not null"`
	IsSystem       bool      `json:"is_system" gorm:"column:is_system;default:true;not null"`
	Level          *int      `json:"level" gorm:"column:level"`
	MainCategoryID *int      `json:"main_category_id" gorm:"column:main_category_id"`
	SubCategoryID  *int      `json:"sub_category_id" gorm:"column:sub_category_id"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}

// TableName specifies the table name for the Word model
func (Word) TableName() string {
	return "words"
}

// WordRequest 単語作成/更新リクエスト構造体
type WordRequest struct {
	Word           string `json:"word" binding:"required"`
	IsSystem       *bool  `json:"is_system,omitempty"`
	Level          *int   `json:"level,omitempty"`
	MainCategoryID *int   `json:"main_category_id,omitempty"`
	SubCategoryID  *int   `json:"sub_category_id,omitempty"`
}

// WordResponse 単語レスポンス構造体
type WordResponse struct {
	ID             uint      `json:"id"`
	Word           string    `json:"word"`
	IsSystem       bool      `json:"is_system"`
	Level          *int      `json:"level"`
	MainCategoryID *int      `json:"main_category_id"`
	SubCategoryID  *int      `json:"sub_category_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ToResponse はWordをWordResponseに変換します
func (w *Word) ToResponse() WordResponse {
	return WordResponse{
		ID:             w.ID,
		Word:           w.Word,
		IsSystem:       w.IsSystem,
		Level:          w.Level,
		MainCategoryID: w.MainCategoryID,
		SubCategoryID:  w.SubCategoryID,
		CreatedAt:      w.CreatedAt,
		UpdatedAt:      w.UpdatedAt,
	}
}

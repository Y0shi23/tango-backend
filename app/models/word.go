package models

import "time"

// Word構造体 - 既存のテーブル構造に合わせて修正
type Word struct {
	ID             uint      `json:"id" gorm:"primary_key;column:id"`
	WordID         uint      `json:"word_id" gorm:"column:word_id;autoIncrement"`
	Word           string    `json:"word" gorm:"column:word;size:100;not null"`
	IsSystem       bool      `json:"is_system" gorm:"column:is_system;default:true;not null"`
	Level          *int      `json:"level" gorm:"column:level"`
	MainCategoryID *int      `json:"main_category_id" gorm:"column:main_category_id"`
	SubCategoryID  *int      `json:"sub_category_id" gorm:"column:sub_category_id"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
	// 新しいフィールド（NULL許可で追加）
	JapaneseMeaning *string `json:"japanese_meaning,omitempty" gorm:"column:japanese_meaning;type:text"`
	PartOfSpeech    *string `json:"part_of_speech,omitempty" gorm:"column:part_of_speech;size:20"`
	DifficultyLevel *int    `json:"difficulty_level,omitempty" gorm:"column:difficulty_level;check:difficulty_level >= 1 AND difficulty_level <= 3"`
}

// TableName specifies the table name for the Word model
func (Word) TableName() string {
	return "words"
}

// WordRequest 単語作成/更新リクエスト構造体
type WordRequest struct {
	Word            string  `json:"word" binding:"required"`
	IsSystem        *bool   `json:"is_system,omitempty"`
	Level           *int    `json:"level,omitempty"`
	MainCategoryID  *int    `json:"main_category_id,omitempty"`
	SubCategoryID   *int    `json:"sub_category_id,omitempty"`
	JapaneseMeaning *string `json:"japanese_meaning,omitempty"`
	PartOfSpeech    *string `json:"part_of_speech,omitempty"`
	DifficultyLevel *int    `json:"difficulty_level,omitempty" binding:"omitempty,min=1,max=3"`
}

// WordResponse 単語レスポンス構造体
type WordResponse struct {
	ID              uint      `json:"id"`
	WordID          uint      `json:"word_id"`
	Word            string    `json:"word"`
	IsSystem        bool      `json:"is_system"`
	Level           *int      `json:"level"`
	MainCategoryID  *int      `json:"main_category_id"`
	SubCategoryID   *int      `json:"sub_category_id"`
	JapaneseMeaning *string   `json:"japanese_meaning,omitempty"`
	PartOfSpeech    *string   `json:"part_of_speech,omitempty"`
	DifficultyLevel *int      `json:"difficulty_level,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ToResponse はWordをWordResponseに変換します
func (w *Word) ToResponse() WordResponse {
	return WordResponse{
		ID:              w.ID,
		WordID:          w.WordID,
		Word:            w.Word,
		IsSystem:        w.IsSystem,
		Level:           w.Level,
		MainCategoryID:  w.MainCategoryID,
		SubCategoryID:   w.SubCategoryID,
		JapaneseMeaning: w.JapaneseMeaning,
		PartOfSpeech:    w.PartOfSpeech,
		DifficultyLevel: w.DifficultyLevel,
		CreatedAt:       w.CreatedAt,
		UpdatedAt:       w.UpdatedAt,
	}
}

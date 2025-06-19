package models

import "time"

// Word構造体
type Word struct {
	ID              uint      `json:"id" gorm:"primary_key;column:word_id"`
	EnglishWord     string    `json:"english_word" gorm:"not null;size:100"`
	JapaneseMeaning string    `json:"japanese_meaning" gorm:"not null;type:text"`
	PartOfSpeech    string    `json:"part_of_speech" gorm:"not null;size:20;check:part_of_speech in ('NOUN', 'VERB', 'ADJECTIVE', 'ADVERB', 'PREPOSITION', 'CONJUNCTION', 'INTERJECTION', 'PRONOUN')"`
	DifficultyLevel int       `json:"difficulty_level" gorm:"default:2;check:difficulty_level >= 1 AND difficulty_level <= 3"`
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for the Word model
func (Word) TableName() string {
	return "words"
}

// WordRequest 単語作成/更新リクエスト構造体
type WordRequest struct {
	EnglishWord     string `json:"english_word" binding:"required"`
	JapaneseMeaning string `json:"japanese_meaning" binding:"required"`
	PartOfSpeech    string `json:"part_of_speech" binding:"required"`
	DifficultyLevel int    `json:"difficulty_level" binding:"min=1,max=3"`
}

// WordResponse 単語レスポンス構造体
type WordResponse struct {
	ID              uint      `json:"id"`
	EnglishWord     string    `json:"english_word"`
	JapaneseMeaning string    `json:"japanese_meaning"`
	PartOfSpeech    string    `json:"part_of_speech"`
	DifficultyLevel int       `json:"difficulty_level"`
	CreatedAt       time.Time `json:"created_at"`
}

// ToResponse はWordをWordResponseに変換します
func (w *Word) ToResponse() WordResponse {
	return WordResponse{
		ID:              w.ID,
		EnglishWord:     w.EnglishWord,
		JapaneseMeaning: w.JapaneseMeaning,
		PartOfSpeech:    w.PartOfSpeech,
		DifficultyLevel: w.DifficultyLevel,
		CreatedAt:       w.CreatedAt,
	}
}

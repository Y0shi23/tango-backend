package database

import (
	"database/sql"
	"fmt"
	"log"

	"backend/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect は設定を使用してデータベースに接続します
func Connect(dbConfig config.DatabaseConfig) error {
	// 接続文字列を作成
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// データベース接続確認
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return nil
}

// Close はデータベース接続を閉じます
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB はデータベースインスタンスを返します
func GetDB() *sql.DB {
	return DB
}

// ValidateConnection はデータベース接続を検証します
func ValidateConnection() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	return DB.Ping()
}

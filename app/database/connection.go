package database

import (
	"fmt"
	"log"

	"backend/config"
	"backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect は設定を使用してデータベースに接続します
func Connect(dbConfig config.DatabaseConfig) error {
	// 接続文字列を作成
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.Port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// マイグレーション実行
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Successfully connected to database and migrated tables")
	return nil
}

// Close はデータベース接続を閉じます
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB はデータベースインスタンスを返します
func GetDB() *gorm.DB {
	return DB
}

// ValidateConnection はデータベース接続を検証します
func ValidateConnection() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

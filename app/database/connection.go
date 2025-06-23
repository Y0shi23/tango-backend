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

	// カスタムマイグレーション処理
	if err := customMigration(); err != nil {
		return fmt.Errorf("failed to run custom migration: %w", err)
	}

	// マイグレーション実行
	if err := DB.AutoMigrate(&models.User{}, &models.Word{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Successfully connected to database and migrated tables")
	return nil
}

// customMigration はカスタムマイグレーション処理を実行します
func customMigration() error {
	// wordsテーブルが存在するかチェック
	if DB.Migrator().HasTable(&models.Word{}) {
		// 新しいカラムを安全に追加
		columnsToAdd := []struct {
			name    string
			sqlType string
		}{
			{"japanese_meaning", "TEXT"},
			{"part_of_speech", "VARCHAR(20)"},
			{"difficulty_level", "INTEGER"},
		}

		for _, col := range columnsToAdd {
			if !DB.Migrator().HasColumn(&models.Word{}, col.name) {
				sql := fmt.Sprintf("ALTER TABLE words ADD COLUMN %s %s", col.name, col.sqlType)
				if err := DB.Exec(sql).Error; err != nil {
					log.Printf("Warning: Failed to add column %s: %v", col.name, err)
					// エラーがあっても続行（カラムが既に存在する可能性）
				} else {
					log.Printf("Added %s column to words table", col.name)
				}
			}
		}

		// difficulty_levelにチェック制約を追加（存在しない場合のみ）
		checkConstraintSQL := `
			DO $$ 
			BEGIN
				IF NOT EXISTS (
					SELECT 1 FROM information_schema.table_constraints 
					WHERE constraint_name = 'words_difficulty_level_check' 
					AND table_name = 'words'
				) THEN
					ALTER TABLE words ADD CONSTRAINT words_difficulty_level_check 
					CHECK (difficulty_level >= 1 AND difficulty_level <= 3);
				END IF;
			END $$;
		`
		if err := DB.Exec(checkConstraintSQL).Error; err != nil {
			log.Printf("Warning: Failed to add difficulty_level check constraint: %v", err)
		} else {
			log.Println("Added difficulty_level check constraint")
		}
	}
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

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config はアプリケーション設定構造体
type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
}

// DatabaseConfig はデータベース設定
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// JWTConfig はJWT設定
type JWTConfig struct {
	Secret string
}

// ServerConfig はサーバー設定
type ServerConfig struct {
	Port         string
	AllowOrigins []string
}

// Load は設定を読み込みます
func Load() *Config {
	// 環境変数を読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", ""),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			AllowOrigins: []string{
				"http://localhost:3000",
				"http://127.0.0.1:3000",
				"http://10.0.2.2:8080",  // Android エミュレータからホストへ
				"http://localhost:8080", // ローカルホスト
				"http://192.168.1.0/24", // ローカルネットワーク（実機用）
				"*",                     // 開発時のみ - 本番環境では削除してください
			},
		},
	}

	// JWT秘密鍵の警告
	if config.JWT.Secret == "your-secret-key-change-this-in-production" {
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}

	return config
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返します
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

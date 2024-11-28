package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DB           *gorm.DB
	JWTSecret    string
	CookieDomain string
}

func LoadConfig() (*Config, error) {
    err := godotenv.Load()
    if err != nil {
        return nil, fmt.Errorf("error loading .env file: %w", err)
    }

    // Gunakan DATABASE_URL dari .env file
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return nil, fmt.Errorf("DATABASE_URL not set in .env")
    }

    db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
        // Tambahkan konfigurasi tambahan jika diperlukan
        // Misalnya: ConnectTimeout, IdleConnTimeout
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Enable UUID extension
    if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
        return nil, fmt.Errorf("failed to create UUID extension: %w", err)
    }

    return &Config{
        DB:           db,
        JWTSecret:    os.Getenv("JWT_SECRET"),
    }, nil
}
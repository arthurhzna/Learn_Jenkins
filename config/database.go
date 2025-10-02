package config

import (
	"fmt"
	"net/url"
	"time"

	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}
	name := os.Getenv("DB_NAME")

	encodedPassword := url.QueryEscape(password)
	uri := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		username,
		encodedPassword,
		host,
		port,
		name,
	)

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Duration(300) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(60) * time.Second)
	return db, nil
}

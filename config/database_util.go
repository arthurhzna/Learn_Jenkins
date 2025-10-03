package config

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"errors"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateTestDatabase() error {
	err := godotenv.Load("../.env")
	if err != nil {
		return err
	}

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME_TESTING")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid DB_PORT: %w", err)
	}

	encodedPassword := url.QueryEscape(password)
	uri := fmt.Sprintf("postgresql://%s:%s@%s:%d/postgres?sslmode=disable",
		username, encodedPassword, host, port)

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %w", err)
	}
	defer sqlDB.Close()

	var tmp int
	err = sqlDB.QueryRow("SELECT 1 FROM pg_database WHERE datname = $1", name).Scan(&tmp)
	if err == nil {
		fmt.Printf("database %s already exists\n", name)
		return nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	safeName := strings.ReplaceAll(name, `"`, `""`)

	_, err = sqlDB.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, safeName))
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	fmt.Println("Database Test created successfully")
	return nil
}

func DropTestDatabase() error {
	err := godotenv.Load("../.env")
	if err != nil {
		return fmt.Errorf("failed to load .env: %w", err)
	}

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME_TESTING")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid DB_PORT: %w", err)
	}

	encodedPassword := url.QueryEscape(password)
	uri := fmt.Sprintf("postgresql://%s:%s@%s:%d/postgres?sslmode=disable",
		username, encodedPassword, host, port)

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %w", err)
	}
	defer sqlDB.Close()

	safeName := strings.ReplaceAll(name, `"`, `""`)

	_, err = sqlDB.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, safeName))
	if err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	fmt.Println("Database Test dropped successfully")
	return nil
}

func InitTestDatabase() (*gorm.DB, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		return nil, fmt.Errorf("failed to load .env: %w", err)
	}

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME_TESTING")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	encodedPassword := url.QueryEscape(password)
	uri := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		username, encodedPassword, host, port, name)

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	return db, nil
}

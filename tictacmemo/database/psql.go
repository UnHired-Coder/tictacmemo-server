package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDatabase() *gorm.DB {

	dbHost := os.Getenv("DB_HOST")

	sslmode := "require"

	if os.Getenv("ENV") == "test" {
		dbHost = os.Getenv("TEST_DB_HOST")
		sslmode = "disable"
	}

	logMode := logger.Silent
	if os.Getenv("ENV") == "dev" {
		logMode = logger.Info
		sslmode = "disable"
	}

	var db *gorm.DB

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		dbHost,
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		sslmode,
	)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{}, &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})

	if err != nil {
		log.Fatal("Unable to connect to the database " + err.Error())
	}

	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetConnMaxLifetime(time.Hour)
		sqlDB.SetMaxOpenConns(100)
	} else {
		log.Fatal("Failed to set connection pool parameters")
	}

	return db
}

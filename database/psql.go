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

	if os.Getenv("ENV") == "test" {
		dbHost = os.Getenv("TEST_DB_HOST")
	}

	logMode := logger.Silent
	if os.Getenv("ENV") == "dev" {
		logMode = logger.Info
	}

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=prefer",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		dbHost,
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	log.Println(connectionString)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
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

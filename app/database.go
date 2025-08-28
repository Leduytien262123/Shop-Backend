package app

import (
	"backend/internal/model"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	// First, connect without database to create it if needed
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	// Configure GORM with less verbose logging
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // Chỉ hiển thị warning và error
	}

	// Connect without database first
	tempDB, err := gorm.Open(mysql.Open(dsnWithoutDB), config)
	if err != nil {
		log.Fatal("❌ Failed to connect to MySQL server:", err)
	}

	// Create database if not exists
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "backend"
	}
	
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)
	if err := tempDB.Exec(createDBSQL).Error; err != nil {
		log.Fatal("❌ Failed to create database:", err)
	}
	
	log.Printf("✅ Database '%s' ensured to exist", dbName)

	// Now connect to the specific database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		dbName,
	)

	// Connect to database
	database, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// Configure connection pool
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("❌ Failed to get database instance:", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = database

	// Run migrations
	if err := runMigrations(); err != nil {
		log.Fatal("❌ Failed to run migrations:", err)
	}

	log.Println("✅ Database connected and migrated successfully!")
}

func runMigrations() error {
	// Auto-migrate all models
	return DB.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Product{},
	)
}

func GetDB() *gorm.DB {
	return DB
}

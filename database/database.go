package database

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"your-app/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	instance *gorm.DB
	once     sync.Once
)

// GetDB returns the singleton database instance (SQLite)
func GetDB() *gorm.DB {
	once.Do(func() {
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "./data.db"
		}
		if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			log.Fatal("Failed to create database directory: ", err)
		}

		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatal("Failed to connect to database: ", err)
		}

		log.Println("connected to SQLite:", dbPath)
		db.Logger = logger.Default.LogMode(logger.Info)
		log.Println("running migrations")

		if err := db.AutoMigrate(&models.Lead{}); err != nil {
			log.Fatal("Failed to run migrations: ", err)
		}
		instance = db
	})
	return instance
}

// ConnectDb initializes the database connection
func ConnectDb() {
	GetDB()
}

package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"leads-import/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func getPostgresDSN() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := os.Getenv("DB_PASSWORD")
	dbname := getEnv("DB_NAME", "your_app")
	sslmode := getEnv("DB_SSLMODE", "disable")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// GetDB returns the singleton database instance (SQLite or PostgreSQL based on DB_DRIVER)
func GetDB() *gorm.DB {
	once.Do(func() {
		driver := getEnv("DB_DRIVER", "sqlite")
		var dialector gorm.Dialector

		switch driver {
		case "postgres":
			dialector = postgres.Open(getPostgresDSN())
		case "sqlite", "":
			dbPath := getEnv("DB_PATH", "./data.db")
			if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
				log.Fatal("Failed to create database directory: ", err)
			}
			dialector = sqlite.Open(dbPath)
		default:
			log.Fatalf("Unknown DB_DRIVER: %q (use sqlite or postgres)", driver)
		}

		db, err := gorm.Open(dialector, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatal("Failed to connect to database: ", err)
		}

		log.Printf("connected to %s", driver)
		db.Logger = logger.Default.LogMode(logger.Info)

		if driver == "postgres" {
			db.Exec("CREATE SCHEMA IF NOT EXISTS amigocare")
		}

		if err := db.AutoMigrate(
			&models.Lead{},
			&models.LeadImport{},
			&models.LeadSource{},
			&models.LeadChannel{},
			&models.Tag{},
			&models.ChatTag{},
			&models.Patient{},
			&models.MessagingAccount{},
		); err != nil {
			log.Fatal("Failed to auto-migrate: ", err)
		}
		log.Println("auto-migration completed")

		instance = db
	})
	return instance
}

// ConnectDb initializes the database connection
func ConnectDb() {
	GetDB()
}

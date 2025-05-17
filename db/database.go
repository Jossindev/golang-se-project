package db

import (
	"awesomeProject/utils"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Fallback for local development
		host := utils.GetEnv("DB_HOST", "localhost")
		port := utils.GetEnv("DB_PORT", "5432")
		user := utils.GetEnv("DB_USER", "postgres")
		password := utils.GetEnv("DB_PASSWORD", "postgres")
		dbname := utils.GetEnv("DB_NAME", "weatherapi")

		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user, password, host, port, dbname)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	sqlDB, err := DB.DB()
	if err != nil {
		panic("Failed to get database connection: " + err.Error())
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := RunMigrations(); err != nil {
		panic(fmt.Sprintf("Failed to run migrations: %v", err))
	}

}

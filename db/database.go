package db

import (
	"awesomeProject/models"
	"awesomeProject/utils"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {

	host := utils.GetEnv("DB_HOST", "localhost")
	port := utils.GetEnv("DB_PORT", "5432")
	user := utils.GetEnv("DB_USER", "postgres")
	password := utils.GetEnv("DB_PASSWORD", "postgres")
	dbname := utils.GetEnv("DB_NAME", "userapi")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		host, port, user, password, dbname)

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic("Failed to get database connection: " + err.Error())
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = DB.AutoMigrate(&models.Subscription{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	fmt.Println("Database connection established")
}

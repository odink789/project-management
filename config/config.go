package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	AppConfig *Config
)

type Config struct {
	AppPort          string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	JWTSecret        string
	JWTRefreshToken  string
	JWTExpire        string
}

//function file untuk load file .env

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found.")
	}
	AppConfig = &Config{
		AppPort:          getEnv("PORT", "3030"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBUser:           getEnv("DB_USER", "posgres"),
		DBPassword:       getEnv("DB_PASSWORD", "password"),
		JWTSecret:        getEnv("JWT_SECRET", "rahasia"),
		JWTExpire: getEnv("JWT_EXPIRY", "60"),
		JWTRefreshToken:  getEnv("REFRESH_TOKEN_EXPIRED", "24h"),
	}
}

func getEnv(key string, fallback string) string {
	value, exist := os.LookupEnv(key)
	if exist {
		return value
	} else {
		return fallback
	}
}

func ConnectDB() {
	cfg := AppConfig

	dsn := fmt.Sprintf("host=%s port=%s password=%s dbname=%s sslmode=disable",cfg.DBHost
	cfg.DBPort,cfg.DBUser,cfg.DBPassword,cfg.DBName)

	//open conecction ke db

	db,err := gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err != nil {
		log.Fatal("Failed to Connect to database",err)
	}

	sqlDB,err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance",err)

	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)


	DB = db
	


}

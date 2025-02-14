package Config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// DBConfig represents db configuration
type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"post"`
	User     string `json:"user"`
	DBName   string `json:"dbname"`
	Password string `json:"password"`
}

func BuildDBConfig() *DBConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbConfig := DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
	return &dbConfig
}

func DbURL(dbConfig *DBConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.Port,
	)
}

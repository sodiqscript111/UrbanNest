package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	KafkaBrokers  string
	RedisAddr     string
	RedisPassword string
	ResendAPIKey  string
	Port          string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		KafkaBrokers:  os.Getenv("KAFKA_BROKERS"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		ResendAPIKey:  os.Getenv("RESEND_API_KEY"),
		Port:          os.Getenv("PORT"),
	}
}

package config

import "os"

type Config struct {
	Port          string
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
	KafkaBrokers  string
	RedisAddr     string
	RedisPassword string
	ResendAPIKey  string
	JWTSecret     string
}

func LoadConfig() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBUser:        getEnv("DB_USER", "user"),
		DBPassword:    getEnv("DB_PASSWORD", "password"),
		DBName:        getEnv("DB_NAME", "airbnb"),
		DBPort:        getEnv("DB_PORT", "5432"),
		KafkaBrokers:  getEnv("KAFKA_BROKERS", "localhost:9092"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		ResendAPIKey:  getEnv("RESEND_API_KEY", ""),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

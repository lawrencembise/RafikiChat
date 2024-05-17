package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	SSLMode     string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file",err)
	}

	cfg := Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "rafikichat"),
		SSLMode:    getEnv("SSL_MODE", "disable"),
	}

	log.Printf("Config loaded: %+v", cfg)
	return cfg
}

// getEnv reads an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

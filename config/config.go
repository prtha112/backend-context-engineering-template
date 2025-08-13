package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		Name string
		Env  string
	}
	HTTP struct {
		Addr string
		Port string
	}
	DB struct {
		Driver   string
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	Log struct {
		Level string
	}
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{}

	config.App.Name = getEnv("APP_NAME", "product-service")
	config.App.Env = getEnv("APP_ENV", "development")

	config.HTTP.Addr = getEnv("HTTP_ADDR", "0.0.0.0")
	config.HTTP.Port = getEnv("HTTP_PORT", "8080")

	config.DB.Driver = getEnv("DB_DRIVER", "postgres")
	config.DB.Host = getEnv("DB_HOST", "localhost")
	config.DB.Port = getEnv("DB_PORT", "5432")
	config.DB.User = getEnv("DB_USER", "app_user")
	config.DB.Password = getEnv("DB_PASSWORD", "app_password")
	config.DB.Name = getEnv("DB_NAME", "product_db")
	config.DB.SSLMode = getEnv("DB_SSLMODE", "disable")

	config.Log.Level = getEnv("LOG_LEVEL", "info")

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

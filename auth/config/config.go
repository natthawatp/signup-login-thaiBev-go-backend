package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI  string
	MongoDB   string
	JWTSecret string
	Port      string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: .env file not found, using system env")
	}

	return &Config{
		MongoURI:  getEnv("MONGO_URI", ""),
		MongoDB:   getEnv("MONGO_DB", ""),
		JWTSecret: getEnv("JWT_SECRET", ""),
		Port:      getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

package bootstrap

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Env struct {
	ServerPort         string
	ContextTimeout     time.Duration
	DBUri              string
	DBName             string
	GeminiAPIKey       string
	AccessTokenSecret  string
	AccessTokenExpiry  time.Duration
	RefreshTokenSecret string
	RefreshTokenExpiry time.Duration
}

func NewEnv() *Env {
	env := Env{}
	env.loadEnv()
	env.ServerPort = os.Getenv("SERVER_PORT")
	env.ContextTimeout = env.getDuration("CONTEXT_TIMEOUT", 10*time.Second)
	env.DBUri = os.Getenv("MONGODB_URI")
	env.DBName = os.Getenv("DB_NAME")
	env.GeminiAPIKey = os.Getenv("GEMINI_API_KEY")
	env.AccessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
	env.AccessTokenExpiry = env.getDuration("ACCESS_TOKEN_EXPIRY", 15*time.Minute)
	env.RefreshTokenSecret = os.Getenv("REFRESH_TOKEN_SECRET")
	env.RefreshTokenExpiry = env.getDuration("REFRESH_TOKEN_EXPIRY", 24*time.Hour)

	// Debug logging
	log.Printf("Environment variables loaded:")
	log.Printf("SERVER_PORT: %s", env.ServerPort)
	log.Printf("MONGODB_URI: %s", env.DBUri)
	log.Printf("DB_NAME: %s", env.DBName)
	log.Printf("GEMINI_API_KEY: %s", env.GeminiAPIKey)
	log.Printf("ACCESS_TOKEN_SECRET: %s", env.AccessTokenSecret)
	log.Printf("REFRESH_TOKEN_SECRET: %s", env.RefreshTokenSecret)

	return &env
}

func (e *Env) loadEnv() {
	// Try multiple possible locations for the .env file
	envPaths := []string{
		".env",                         // Current directory
		"interview-coach-backend/.env", // Relative to current directory
		"../.env",                      // Parent directory
		"../../.env",                   // Two levels up
	}

	var loaded bool
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("Successfully loaded .env file from: %s", path)
			loaded = true
			break
		}
	}

	if !loaded {
		log.Println("Warning: No .env file found in any of the expected locations. Using default values.")
		log.Println("Expected locations:")
		for _, path := range envPaths {
			absPath, _ := filepath.Abs(path)
			log.Printf("- %s", absPath)
		}
	}
}

func (e *Env) getDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return time.Duration(value) * time.Second
}

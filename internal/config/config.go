package config

import (
    "log"
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type Config struct {
    Port            string
    DBHost          string
    DBUser          string
    DBPassword      string
    DBName          string
    DBPort          string
    DBSSLMode       string
    RedisHost       string
    RedisPort       string
    RedisPassword   string
    JWTSecret       string
    JWTExpiryHours  int
    BaseURL         string
}

func LoadConfig() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    jwtExpiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

    return &Config{
        Port:           getEnv("PORT", "8080"),
        DBHost:         getEnv("DB_HOST", "localhost"),
        DBUser:         getEnv("DB_USER", "postgres"),
        DBPassword:     getEnv("DB_PASSWORD", "postgres"),
        DBName:         getEnv("DB_NAME", "urlshortener"),
        DBPort:         getEnv("DB_PORT", "5432"),
        DBSSLMode:      getEnv("DB_SSL_MODE", "disable"),
        RedisHost:      getEnv("REDIS_HOST", "localhost"),
        RedisPort:      getEnv("REDIS_PORT", "6379"),
        RedisPassword:  getEnv("REDIS_PASSWORD", ""),
        JWTSecret:      getEnv("JWT_SECRET", "your-super-secret-key"),
        JWTExpiryHours: jwtExpiryHours,
        BaseURL:        getEnv("BASE_URL", "http://localhost:8080"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
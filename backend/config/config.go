package config

import (
	"os"
)

type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPass            string
	DBName            string
	JWTSecret         string
	ServerPort        string
	GDriveCredentials string
	GDriveFolderID    string
}

func Load() *Config {
	return &Config{
		DBHost:            getEnv("DB_HOST", "127.0.0.1"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPass:            getEnv("DB_PASS", "100%Bisa"),
		DBName:            getEnv("DB_NAME", "pkl_db"),
		JWTSecret:         getEnv("JWT_SECRET", "pkl-tracker-secret-key-2026"),
		ServerPort:        getEnv("SERVER_PORT", "8082"),
		GDriveCredentials: getEnv("GDRIVE_CREDENTIALS", ""),
		GDriveFolderID:    getEnv("GDRIVE_FOLDER_ID", ""),
	}
}

func (c *Config) DSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPass +
		" dbname=" + c.DBName +
		" sslmode=disable TimeZone=Asia/Jakarta"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

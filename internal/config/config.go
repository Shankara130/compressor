package config

import (
	"os"
	"runtime"
	"strconv"
)

type Config struct {
	ServerPort   string
	RedisAddr    string
	InputDir     string
	OutputDir    string
	MaxFileSize  int64
	WorkerCount  int
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int

	CleanupIntervalSec int
}

func Load() *Config {
	return &Config{
		ServerPort:         getEnv("SERVER_PORT", getEnv("PORT", "8001")),
		RedisAddr:          getEnv("REDIS_ADDR", "localhost:6379"),
		InputDir:           getEnv("INPUT_DIR", "tmp/input"),
		OutputDir:          getEnv("OUTPUT_DIR", "tmp/output"),
		MaxFileSize:        50 << 20,
		WorkerCount:        getEnvAsInt("WORKER_COUNT", runtime.NumCPU()),
		ReadTimeout:        getEnvAsInt("READ_TIMEOUT", 15),
		WriteTimeout:       getEnvAsInt("WRITE_TIMEOUT", 15),
		IdleTimeout:        getEnvAsInt("IDLE_TIMEOUT", 60),
		CleanupIntervalSec: getEnvAsInt("CLEANUP_INTERVAL_SEC", 3600),
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

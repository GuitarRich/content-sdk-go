package utils

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var envLoadOnce sync.Once

// loadDotEnvOnce loads variables from a local .env file if present.
// In deployed environments where .env is not present, this is a no-op.
func loadDotEnvOnce() {
	envLoadOnce.Do(func() {
		_ = godotenv.Load()
	})
}

// GetEnvVar returns the value of the environment variable named by key.
// It supports local development by loading a .env file if present.
func GetEnvVar(key string) string {
	loadDotEnvOnce()
	return os.Getenv(key)
}

// GetEnvVarOrDefault returns the env var value or the provided default if unset.
func GetEnvVarOrDefault(key string, defaultValue string) string {
	value := GetEnvVar(key)
	if value == "" {
		return defaultValue
	}
	return value
}

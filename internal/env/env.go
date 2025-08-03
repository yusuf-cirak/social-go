package env

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}

func GetInt(key string, fallback int) int {
	valueStr, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return fallback
	}
	return value
}

func GetBool(key string, fallback bool) bool {
	valueStr, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return fallback
	}
	return value
}

func GetFloat64(key string, fallback float64) float64 {
	valueStr, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return fallback
	}
	return value
}

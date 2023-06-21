package environment

import (
	"os"
	"strconv"
	"strings"
)

func GetString(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return strings.TrimSpace(value)
}

func GetInt64(key string, fallback int64) int64 {
	value := GetString(key, "")
	if len(value) == 0 {
		return fallback
	}
	valueInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}

	return valueInt
}

func GetBool(key string, fallback bool) bool {
	value := GetString(key, "")
	if len(value) == 0 {
		return fallback
	}
	valueInt, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return valueInt
}

package helpers

import (
	"os"
	"strconv"
)

// GetEnvironmentVariableAsString returns environment variable as a string
func GetEnvironmentVariableAsString(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// GetEnvironmentVariableAsBool returns environment variable as a boolean
func GetEnvironmentVariableAsBool(key string, fallback bool) bool {
	var newValue bool
	value, exists := os.LookupEnv(key)
	if !exists {
		newValue = fallback
	} else {
		newValue = value == "true"
	}
	return newValue
}

// GetEnvironmentVariableAsInteger returns environment variable as an integer
func GetEnvironmentVariableAsInteger(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsedValue
}

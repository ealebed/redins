package helpers

import (
	"os"
	"strconv"
)

// GetEnvironmentVariableAsString returns environment variable as a string
func GetEnvironmentVariableAsString(key string, fallback string) string {
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
	var newValue int
	value, exists := os.LookupEnv(key)
	if !exists {
		newValue = fallback
	} else {
		value, _ := strconv.Atoi(value)
		newValue = value
	}
	return newValue
}

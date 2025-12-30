package helpers

import (
	"os"
	"testing"
)

func TestGetEnvironmentVariableAsString(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		setValue string
		fallback string
		want     string
	}{
		{
			name:     "environment variable exists",
			key:      "TEST_STRING_VAR",
			setValue: "test-value",
			fallback: "default-value",
			want:     "test-value",
		},
		{
			name:     "environment variable does not exist",
			key:      "TEST_STRING_VAR_MISSING",
			setValue: "",
			fallback: "default-value",
			want:     "default-value",
		},
		{
			name:     "environment variable is empty string",
			key:      "TEST_STRING_VAR_EMPTY",
			setValue: "",
			fallback: "default-value",
			want:     "",
		},
		{
			name:     "fallback is empty string",
			key:      "TEST_STRING_VAR_FALLBACK_EMPTY",
			setValue: "",
			fallback: "",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up any existing environment variable
			os.Unsetenv(tt.key)
			defer os.Unsetenv(tt.key)

			// Set environment variable if test case requires it
			if tt.name != "environment variable does not exist" {
				os.Setenv(tt.key, tt.setValue)
			}

			got := GetEnvironmentVariableAsString(tt.key, tt.fallback)
			if got != tt.want {
				t.Errorf("GetEnvironmentVariableAsString(%q, %q) = %q, want %q", tt.key, tt.fallback, got, tt.want)
			}
		})
	}
}

func TestGetEnvironmentVariableAsBool(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		setValue string
		fallback bool
		want     bool
	}{
		{
			name:     "environment variable is 'true'",
			key:      "TEST_BOOL_VAR_TRUE",
			setValue: "true",
			fallback: false,
			want:     true,
		},
		{
			name:     "environment variable is 'false'",
			key:      "TEST_BOOL_VAR_FALSE",
			setValue: "false",
			fallback: true,
			want:     false,
		},
		{
			name:     "environment variable is not 'true'",
			key:      "TEST_BOOL_VAR_OTHER",
			setValue: "yes",
			fallback: true,
			want:     false,
		},
		{
			name:     "environment variable does not exist",
			key:      "TEST_BOOL_VAR_MISSING",
			setValue: "",
			fallback: true,
			want:     true,
		},
		{
			name:     "environment variable is empty string",
			key:      "TEST_BOOL_VAR_EMPTY",
			setValue: "",
			fallback: false,
			want:     false,
		},
		{
			name:     "environment variable is 'True' (case sensitive)",
			key:      "TEST_BOOL_VAR_CASE",
			setValue: "True",
			fallback: false,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up any existing environment variable
			os.Unsetenv(tt.key)
			defer os.Unsetenv(tt.key)

			// Set environment variable if test case requires it
			if tt.name != "environment variable does not exist" {
				os.Setenv(tt.key, tt.setValue)
			}

			got := GetEnvironmentVariableAsBool(tt.key, tt.fallback)
			if got != tt.want {
				t.Errorf("GetEnvironmentVariableAsBool(%q, %v) = %v, want %v", tt.key, tt.fallback, got, tt.want)
			}
		})
	}
}

func TestGetEnvironmentVariableAsInteger(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		setValue string
		fallback int
		want     int
	}{
		{
			name:     "environment variable is valid integer",
			key:      "TEST_INT_VAR_VALID",
			setValue: "42",
			fallback: 0,
			want:     42,
		},
		{
			name:     "environment variable is negative integer",
			key:      "TEST_INT_VAR_NEGATIVE",
			setValue: "-10",
			fallback: 0,
			want:     -10,
		},
		{
			name:     "environment variable is zero",
			key:      "TEST_INT_VAR_ZERO",
			setValue: "0",
			fallback: 100,
			want:     0,
		},
		{
			name:     "environment variable does not exist",
			key:      "TEST_INT_VAR_MISSING",
			setValue: "",
			fallback: 99,
			want:     99,
		},
		{
			name:     "environment variable is invalid (not a number)",
			key:      "TEST_INT_VAR_INVALID",
			setValue: "not-a-number",
			fallback: 50,
			want:     50, // Should return fallback on parse error
		},
		{
			name:     "environment variable is empty string",
			key:      "TEST_INT_VAR_EMPTY",
			setValue: "",
			fallback: 77,
			want:     77,
		},
		{
			name:     "environment variable is float string",
			key:      "TEST_INT_VAR_FLOAT",
			setValue: "3.14",
			fallback: 10,
			want:     10, // Should return fallback on parse error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up any existing environment variable
			os.Unsetenv(tt.key)
			defer os.Unsetenv(tt.key)

			// Set environment variable if test case requires it
			if tt.name != "environment variable does not exist" {
				os.Setenv(tt.key, tt.setValue)
			}

			got := GetEnvironmentVariableAsInteger(tt.key, tt.fallback)
			if got != tt.want {
				t.Errorf("GetEnvironmentVariableAsInteger(%q, %d) = %d, want %d", tt.key, tt.fallback, got, tt.want)
			}
		})
	}
}

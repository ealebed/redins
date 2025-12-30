package redis

import (
	"strings"
	"testing"

	"github.com/go-redis/redis"
)

func TestInitRedisClient(t *testing.T) {
	tests := []struct {
		name        string
		redisAddr   string
		redisPasswd string
		redisDb     int
		want        *redisClient
	}{
		{
			name:        "valid client with all parameters",
			redisAddr:   "localhost:6379",
			redisPasswd: "password123",
			redisDb:     0,
			want: &redisClient{
				Addr:       "localhost:6379",
				Password:   "password123",
				DB:         0,
				MaxRetries: 10,
			},
		},
		{
			name:        "valid client with empty password",
			redisAddr:   "127.0.0.1:6379",
			redisPasswd: "",
			redisDb:     1,
			want: &redisClient{
				Addr:       "127.0.0.1:6379",
				Password:   "",
				DB:         1,
				MaxRetries: 10,
			},
		},
		{
			name:        "valid client with custom database",
			redisAddr:   "redis.example.com:6379",
			redisPasswd: "secret",
			redisDb:     5,
			want: &redisClient{
				Addr:       "redis.example.com:6379",
				Password:   "secret",
				DB:         5,
				MaxRetries: 10,
			},
		},
		{
			name:        "valid client with negative database",
			redisAddr:   "localhost:6379",
			redisPasswd: "",
			redisDb:     -1,
			want: &redisClient{
				Addr:       "localhost:6379",
				Password:   "",
				DB:         -1,
				MaxRetries: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InitRedisClient(tt.redisAddr, tt.redisPasswd, tt.redisDb)

			if got == nil {
				t.Fatal("InitRedisClient returned nil")
			}

			if got.Addr != tt.want.Addr {
				t.Errorf("InitRedisClient().Addr = %q, want %q", got.Addr, tt.want.Addr)
			}

			if got.Password != tt.want.Password {
				t.Errorf("InitRedisClient().Password = %q, want %q", got.Password, tt.want.Password)
			}

			if got.DB != tt.want.DB {
				t.Errorf("InitRedisClient().DB = %d, want %d", got.DB, tt.want.DB)
			}

			if got.MaxRetries != tt.want.MaxRetries {
				t.Errorf("InitRedisClient().MaxRetries = %d, want %d", got.MaxRetries, tt.want.MaxRetries)
			}

			// Connection should be nil before Connect() is called
			if got.connection != nil {
				t.Error("InitRedisClient().connection should be nil before Connect() is called")
			}
		})
	}
}

func TestDisconnect(t *testing.T) {
	tests := []struct {
		name        string
		setupClient func() *redisClient
		wantErr     bool
	}{
		{
			name: "nil connection returns no error",
			setupClient: func() *redisClient {
				return &redisClient{
					Addr:       "localhost:6379",
					Password:   "",
					DB:         0,
					MaxRetries: 10,
					connection: nil,
				}
			},
			wantErr: false,
		},
		{
			name: "disconnect with mock connection",
			setupClient: func() *redisClient {
				// Create a mock client that implements UniversalClient
				// Using a real client but we'll test the disconnect path
				client := redis.NewClient(&redis.Options{
					Addr: "localhost:6379",
				})
				return &redisClient{
					Addr:       "localhost:6379",
					Password:   "",
					DB:         0,
					MaxRetries: 10,
					connection: client,
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := tt.setupClient()
			err := client.Disconnect()

			if (err != nil) != tt.wantErr {
				t.Errorf("Disconnect() error = %v, wantErr %v", err, tt.wantErr)
			}

			// For the mock connection test, verify behavior after first disconnect
			// Note: Redis client.Close() closes the connection and subsequent calls return an error
			if tt.name == "disconnect with mock connection" {
				// First disconnect should succeed
				if err != nil {
					t.Errorf("First Disconnect() should succeed, got error = %v", err)
				}

				// Second disconnect may return an error (client is closed), which is expected
				err2 := client.Disconnect()
				if err2 != nil {
					// Accept "client is closed" error as expected behavior
					if !strings.Contains(err2.Error(), "client is closed") {
						t.Errorf("Disconnect() on already disconnected client unexpected error = %v", err2)
					}
				}
			}
		})
	}
}

func TestSetValue_NoConnection(t *testing.T) {
	// Test that SetValue requires connection to be established first
	// Note: Current implementation will panic with nil connection
	client := InitRedisClient("localhost:6379", "", 0)

	// Verify connection is nil
	if client.connection != nil {
		t.Error("Expected connection to be nil before Connect()")
	}

	// SetValue will panic with nil connection, which documents that Connect() must be called first
	// This test verifies the expected behavior
	defer func() {
		if r := recover(); r == nil {
			t.Error("SetValue should panic with nil connection")
		}
	}()

	_ = client.SetValue("test-key", "test-value")
}

func TestQueryValue_NoConnection(t *testing.T) {
	// Test that QueryValue requires connection to be established first
	client := InitRedisClient("localhost:6379", "", 0)

	// Verify connection is nil
	if client.connection != nil {
		t.Error("Expected connection to be nil before Connect()")
	}

	// QueryValue will panic with nil connection, which documents that Connect() must be called first
	defer func() {
		if r := recover(); r == nil {
			t.Error("QueryValue should panic with nil connection")
		}
	}()

	_ = client.QueryValue("test-key")
}

func TestDeleteValue_NoConnection(t *testing.T) {
	// Test that DeleteValue requires connection to be established first
	client := InitRedisClient("localhost:6379", "", 0)

	// Verify connection is nil
	if client.connection != nil {
		t.Error("Expected connection to be nil before Connect()")
	}

	// DeleteValue will panic with nil connection, which documents that Connect() must be called first
	defer func() {
		if r := recover(); r == nil {
			t.Error("DeleteValue should panic with nil connection")
		}
	}()

	_ = client.DeleteValue("test-key")
}

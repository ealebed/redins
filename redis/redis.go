package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type redisClient struct {
	// host:port address.
	Addr string
	// Optional password. Must match the password specified in the
	// requirepass server configuration option.
	Password string
	// Database to be selected after connecting to the server.
	DB int
	// Maximum number of retries before giving up.
	// Default is to not retry failed commands.
	MaxRetries int

	connection redis.UniversalClient
}

// InitRedisClient initialize redisClient
func InitRedisClient(redisAddr, redisPasswd string, redisDb int) *redisClient {
	client := &redisClient{
		Addr:       redisAddr,
		Password:   redisPasswd,
		DB:         redisDb,
		MaxRetries: 10,
	}

	return client
}

// Connect to Redis
func (c *redisClient) Connect() (*redisClient, error) {
	var endpoints []string
	endpoints = append(endpoints, c.Addr)

	c.connection = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:      endpoints,
		Password:   c.Password,
		DB:         c.DB,
		MaxRetries: c.MaxRetries,
	})

	for {
		if pong, err := c.connection.Ping().Result(); err != nil || pong != "PONG" {
			fmt.Printf("Redis Connection Failed: %s\n", err)
			fmt.Println("Retry in 5 second..")
		} else {
			return c, nil
		}
		time.Sleep(5 * time.Second)
	}
}

// Disconnect from Redis
func (c *redisClient) Disconnect() error {
	if c.connection != nil {
		return c.connection.Close()
	}
	return nil
}

// SetValue writes value to Redis
func (c *redisClient) SetValue(key, value string) error {
	err := c.connection.Set(key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// QueryValue returns key value from Redis (as a plain string)
func (c *redisClient) QueryValue(key string) string {
	if val, err := c.connection.Get(key).Result(); err != nil {
		return fmt.Sprintf("Error happened with %v\n", err)
	} else {
		return fmt.Sprintf("Query key '%s', get return value '%s'\n", key, val)
	}
}

// DeleteValue remove key from redis
func (c *redisClient) DeleteValue(key string) error {
	_, err := c.connection.Del(key).Result()
	return err
}

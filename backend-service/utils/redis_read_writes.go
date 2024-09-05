package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis setter wrapper
func Set(c *redis.Client, key string, value interface{}, expiration time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(context.Background(), key, p, expiration).Err()
}

// Redis getter wrapper
func Get(c *redis.Client, key string, dest interface{}) error {
	p, err := c.Get(context.Background(), key).Result()
	if err != nil {
		fmt.Println("Error not unmarshaling")
		return err
	}
	return json.Unmarshal([]byte(p), &dest)
}

package store

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
)

// Client is a struct to store the redis client
type Client struct {
	*redis.Client
}

// NewClient is a function to create a new redis client
func NewClient() (*Client, error) {
	client := &Client{}

	client.Client = redis.NewClient(
		&redis.Options{
			Addr:     viper.GetString("redis.addr"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		},
	)

	// Check if the redis is connected
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return client, nil
}

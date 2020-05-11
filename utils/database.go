package utils

import (
	"os"

	"github.com/go-redis/redis/v7"
)

// RDB keeps track of the redis database client
var RDB *redis.Client

// InitRDB starts the connection
func InitRDB() (*redis.Client, error) {
	opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, err
	}

	opts.DB = 0 // default redis db
	client := redis.NewClient(opts)

	RDB = client
	return RDB, nil
}

// GetRDB returns the redis database client
func GetRDB() *redis.Client {
	return RDB
}

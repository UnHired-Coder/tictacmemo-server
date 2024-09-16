package database

import (
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetRedisDb() redis.UniversalClient {

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	if os.Getenv("ENV") == "test" {
		host = os.Getenv("TEST_REDIS_HOST")
		port = os.Getenv("TEST_REDIS_PORT")
	}

	addrs := []string{fmt.Sprintf("%s:%s", host, port)}

	if os.Getenv("ENV") == "prod" {
		// Having 2 or more elements in Addrs will initiate the client in cluster mode
		addrs = []string{fmt.Sprintf("%s:%s", host, port)}
	}

	Rdb := redis.NewUniversalClient(
		&redis.UniversalOptions{
			Addrs:        addrs,
			Password:     os.Getenv("REDIS_PASSWORD"),
			DialTimeout:  time.Second * 10,
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
			PoolSize:     10,
		},
	)
	return Rdb
}

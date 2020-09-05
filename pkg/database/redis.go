package database

import (
	"github.com/go-redis/redis"
	"log"
	"time"
)

type redisDB struct {
	client *redis.Client
}

func createRedisDatabase() (Database, error) {
	// create a new client with default options
	// Addr is defined by the docker-compose.yml file
	client := redis.NewClient(&redis.Options{
		Addr:     "service.redis:6379",
		Password: "",
		DB:       0,
	})
	var err error
	// Try to Ping the database, sometimes this doesn't work right away so we try it 10 times with a timeout of 1sec
	// between tries.
	for i := 0; i < 10; i++ {
		_, err = client.Ping().Result()
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return nil, &CreateDatabaseError{reason: err.Error()}
	}
	return &redisDB{client: client}, nil
}

// Set: Implementation of the database interface
func (r *redisDB) Set(key string, value interface{}) error {
	log.Println("(Redis DB): saving data:", value, " for key:", key)

	_, err := r.client.Set(key, value, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// Get: Implementation of the database interface
func (r *redisDB) Get(key string) (interface{}, error) {
	log.Println("(Redis DB): retrieving value from key:", key)

	value, err := r.client.Get(key).Result()
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Delete: Implementation of the database interface
func (r *redisDB) Delete(key string) error {
	log.Println("(Redis DB): deleting value from key:", key)

	_, err := r.client.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}

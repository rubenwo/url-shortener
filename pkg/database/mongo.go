package database

import (
	"github.com/go-redis/redis"
	"log"
	"time"
)

type mongoDb struct {
	client *redis.Client
}

func createMongoDatabase() (*mongoDb, error) {
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
	return &mongoDb{client: client}, nil
}

// Set takes a key (string) and value (interface{}) and stores it in the database. If an error occurred this will be returned.
// error is nil if all went well.
func (m *mongoDb) Set(key string, value interface{}) error {
	log.Println("(Redis DB): saving data:", value, " for key:", key)

	_, err := m.client.Set(key, value, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// Get takes a key (string) and returns a value (interface{}) and an error. Error is nil if everything went okay.
func (m *mongoDb) Get(key string) (interface{}, error) {
	log.Println("(Redis DB): retrieving value from key:", key)

	value, err := m.client.Get(key).Result()
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Delete takes a key (string) and deletes the value associated with that key. error is nil when all went well.
func (m *mongoDb) Delete(key string) error {
	log.Println("(Redis DB): deleting value from key:", key)

	_, err := m.client.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}

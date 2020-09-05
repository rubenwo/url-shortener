package main

import (
	"flag"
	"github.com/rubenwo/url-shortener/pkg/database"
	"github.com/rubenwo/url-shortener/pkg/server"
	"log"
)

func main() {
	dbStr := flag.String("database", "redis", "possible database type. Valid values are `jsonDB` and `redis`.")
	flag.Parse()

	db, err := database.Factory(*dbStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Run(db); err != nil {
		log.Fatal(err)
	}
}

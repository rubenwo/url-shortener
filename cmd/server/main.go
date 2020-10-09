package main

import (
	"flag"
	"log"

	"github.com/rubenwo/url-shortener/pkg/database"
	"github.com/rubenwo/url-shortener/pkg/malicious"
	"github.com/rubenwo/url-shortener/pkg/server"
)

func main() {
	dbStr := flag.String("database", "redis", "possible database type. Valid values are `jsonDB` and `redis`.")
	apiStr := flag.String("api", "sb", "possible api to request for malicious URLs. Valid values are `sb` and `dummy`.")
	flag.Parse()

	db, err := database.Factory(*dbStr)
	if err != nil {
		log.Fatal(err)
	}

	api, err := malicious.Factory(*apiStr)
	if err != nil {
		log.Fatal(err)
	}
	defer api.Close()

	if err := server.Run(db, api); err != nil {
		log.Fatal(err)
	}
}

package server

import (
	"github.com/rubenwo/url-shortener/pkg/database"
)

func Run(db database.Database) error {
	r := &api{db: db}
	return r.run()
}

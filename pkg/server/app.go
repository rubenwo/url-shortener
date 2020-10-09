package server

import (
	"github.com/rubenwo/url-shortener/pkg/database"
	"github.com/rubenwo/url-shortener/pkg/malicious"
)

func Run(db database.Database, malAPI malicious.API) error {
	r := &api{db: db, api: malAPI}
	return r.run()
}

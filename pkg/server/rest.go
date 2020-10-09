package server

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"
	"github.com/rs/cors"
	"github.com/rubenwo/url-shortener/pkg/database"
	"github.com/rubenwo/url-shortener/pkg/malicious"
)

type RedirectReq struct {
	Url string `json:"url"`
}

type RedirectResp struct {
	Slug string `json:"slug"`
	Url  string `json:"url"`
}

type api struct {
	db  database.Database
	api malicious.API
}

const slugLength = 5

var fs = http.FileServer(http.Dir("./public"))

func (a *api) run() error {
	rand.Seed(time.Now().UnixNano())
	router := chi.NewRouter()

	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(httprate.LimitByIP(1, 1*time.Second))

	router.Post("/shorten", a.Add)
	router.Handle("/{id:[A-Za-z0-9_!-]+}", http.HandlerFunc(a.redirect))
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		fs.ServeHTTP(w, r)
	})

	handler := cors.Default().Handler(router)

	tlsCfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:         ":6443",
		Handler:      handler,
		TLSConfig:    tlsCfg,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 20,
	}

	// Start the HTTP REST server.
	log.Println("SmartEnergyTable API is running on:", server.Addr)
	return server.ListenAndServeTLS("./certs/server.pem", "./certs/server.key")
}

func (a *api) redirect(w http.ResponseWriter, r *http.Request) {
	target := ""

	slug := chi.URLParam(r, "id")
	fmt.Printf("redirecting for slug: %s\n", slug)

	if slug != "" && len(slug) == slugLength {
		val, err := a.db.Get(slug)
		if err != nil {
			target = "https://" + r.Host
		} else {
			u := val.(string)
			target = u
		}
	} else {
		target = "https://" + r.Host + "/"
	}

	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}

	http.Redirect(w, r, target, http.StatusPermanentRedirect)
}

func (a *api) Add(w http.ResponseWriter, r *http.Request) {
	//start := time.Now()
	var msg RedirectReq

	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println(err)
		writeJsonError(w, err, http.StatusUnprocessableEntity)
		return
	}

	//fmt.Printf("reading request body & decoding json took %d microseconds\n", time.Since(start).Microseconds())
	//start = time.Now()
	if !strings.HasPrefix(msg.Url, "http") {
		msg.Url = "http://" + msg.Url
	}

	valid := isValidURL(msg.Url)
	if !valid {
		log.Printf("%s is not a valid url", msg.Url)
		writeJsonError(w, fmt.Errorf("'%s' is not a valid url", msg.Url), http.StatusUnprocessableEntity)
		return
	}

	done := make(chan malicious.APIResp)
	ctx := r.Context()
	// Call the api check
	go a.api.Check(ctx, done, msg.Url)

	select {
	case resp := <-done:
		if resp.Err != nil {
			log.Println(resp.Err)
			writeJsonError(w, resp.Err, http.StatusInternalServerError)
			return
		}

		if !resp.Valid {
			log.Printf("%s is a malicious url", msg.Url)
			writeJsonError(w, fmt.Errorf("'%s' is a malicious url", msg.Url), http.StatusForbidden)
			return
		}
	// In case of WriteTimeout reached
	case <-ctx.Done():
		err := ctx.Err()
		log.Println(err)
		writeJsonError(w, err, http.StatusInternalServerError)
		return
	}

	//fmt.Printf("validating url took %d microseconds\n", time.Since(start).Microseconds())
	//start = time.Now()

	slug := generateSlug(slugLength)

	//fmt.Printf("generating slug took %d microseconds\n", time.Since(start).Microseconds())
	//start = time.Now()

	if err := a.db.Set(slug, msg.Url); err != nil {
		log.Println(err)
		writeJsonError(w, err, http.StatusInternalServerError)
		return

	}
	//fmt.Printf("settings database values took %d microseconds\n", time.Since(start).Microseconds())
	//start = time.Now()

	res := RedirectResp{
		Slug: slug,
		Url:  msg.Url,
	}

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		log.Println(err)
		writeJsonError(w, err, http.StatusBadGateway)
		return
	}

	//fmt.Printf("sending result took %d microseconds\n", time.Since(start).Microseconds())
}

func isValidURL(s string) bool {
	if s == "" {
		return false
	}
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

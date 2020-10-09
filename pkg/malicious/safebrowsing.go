package malicious

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/safebrowsing"
	"github.com/joho/godotenv"
)

type sbAPI struct {
	*safebrowsing.SafeBrowser
}

func createSafeBrowsingAPI() (API, error) {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if os.Getenv("API_KEY") == "" {
		return nil, fmt.Errorf("no API_KEY provided in env vars")
	}

	config := safebrowsing.Config{
		APIKey: os.Getenv("API_KEY"),
		DBPath: "malicious.db.gz", // So the total number of requests is limited
		Logger: os.Stdout,         // Logging automatically handled by the library
	}

	sb, err := safebrowsing.NewSafeBrowser(config)
	if err != nil {
		return nil, &CreateAPIError{err.Error()}
	}

	return &sbAPI{sb}, nil
}

// Check is the implementation of the API interface
func (sb *sbAPI) Check(ctx context.Context, done chan<- APIResp, url string) {
	log.Printf("(SafeBrowsing): checking url: %s\n", url)

	threats, err := sb.LookupURLsContext(ctx, []string{url})
	if err != nil {
		done <- APIResp{Valid: false, Err: err}
		return
	}

	// As LookupURLs return a slice of slice of len == URLs -> there is only one elem
	if len(threats[0]) != 0 {
		done <- APIResp{Valid: false, Err: nil}
		return
	}

	done <- APIResp{Valid: true, Err: nil}
}

func (sb *sbAPI) Close() error {
	return sb.Close()
}

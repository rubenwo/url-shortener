package malicious

import (
	"context"
	"io"
)

// APIResp allows to communicate API response through channel
type APIResp struct {
	Valid bool
	Err   error
}

// API is an interface that defines the actions a database should be able to execute.
type API interface {
	// So we can terminates the communication when the program ends
	io.Closer
	// Check takes a key (string) and value (interface{}) and stores it in the database.
	// If an error occurred this will be returned.
	Check(ctx context.Context, done chan<- APIResp, url string)
}

// Factory is a factory function that creates an api based on the apiName param.
// apiName wether "sb", "dummy" or "". Otherwise, it will return an NotImplementedAPIError
func Factory(apiName string) (API, error) {
	switch apiName {
	case "sb":
		return createSafeBrowsingAPI()
	case "", "dummy":
		return createDummyAPI()
	default:
		return nil, &NotImplementedAPIError{api: apiName}
	}
}

package malicious

import "context"

// This is just a dummy interface implementation that does nothing
// Just built for compatibility of malicious.API if we don't want any
// verification for any reasons (USE AT YOUR OWN RISK !!)
type dummyAPI struct{}

func createDummyAPI() (API, error) {
	return &dummyAPI{}, nil
}

// Check is the implementation of the MaliciousAPI interface
func (sb *dummyAPI) Check(ctx context.Context, done chan<- APIResp, url string) {
	done <- APIResp{Valid: true, Err: nil}
}

func (sb *dummyAPI) Close() error {
	return nil
}

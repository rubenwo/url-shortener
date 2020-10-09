package malicious

import "fmt"

// OperationError indicates the operation could not be performed.
type OperationError struct {
	operation string
}

func (err *OperationError) Error() string {
	return fmt.Sprintf("could not perform the: %s operation", err.operation)
}

// DownError indicates that the api is not reachable.
type DownError struct{}

func (err *DownError) Error() string {
	return "could not connect to the api"
}

// CreateAPIError contains a reason as to why the api couldn't be created.
type CreateAPIError struct {
	reason string
}

func (err *CreateAPIError) Error() string {
	return fmt.Sprintf("could not create api: %s", err.reason)
}

// NotImplementedAPIError is an error returned when this api mode does not exist.
type NotImplementedAPIError struct {
	api string
}

func (err *NotImplementedAPIError) Error() string {
	return fmt.Sprintf("%s not implemented", err.api)
}

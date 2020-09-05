package database

import "fmt"

// OperationError indicates the operation could not be performed.
type OperationError struct {
	operation string
}

func (err *OperationError) Error() string {
	return fmt.Sprintf("could not perform the: %s operation", err.operation)
}

// DownError indicates that the database is not reachable.
type DownError struct{}

func (err *DownError) Error() string {
	return "could not connect to the database"
}

// CreateDatabaseError contains a reason as to why the database couldn't be created.
type CreateDatabaseError struct {
	reason string
}

func (err *CreateDatabaseError) Error() string {
	return fmt.Sprintf("could not create database: %s", err.reason)
}

// NotImplementedDatabaseError is an error returned when this database mode does not exist.
// eg. Trying to use a MySQL database when no implementation is created for MySQL.
type NotImplementedDatabaseError struct {
	database string
}

func (err *NotImplementedDatabaseError) Error() string {
	return fmt.Sprintf("%s not implemented", err.database)
}


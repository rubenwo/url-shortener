package database

// Database is an interface that defines the actions a database should be able to execute.
type Database interface {
	// Set takes a key (string) and value (interface{}) and stores it in the database. If an error occurred this will be returned.
	// error is nil if all went well.
	Set(key string, value interface{}) error
	// Get takes a key (string) and returns a value (interface{}) and an error. Error is nil if everything went okay.
	Get(key string) (interface{}, error)
	// Delete takes a key (string) and deletes the value associated with that key. error is nil when all went well.
	Delete(key string) error
}

// Factory is a factory function that creates a database based on the databaseName param.
// The databaseName must be "redis".
func Factory(databaseName string) (Database, error) {
	switch databaseName {
	case "redis":
		return createRedisDatabase()
	default:
		return nil, &NotImplementedDatabaseError{database: databaseName}
	}
}

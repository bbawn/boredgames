package errors

import (
	"fmt"
)

// NotFoundError indicates the object with the given key was not found
type NotFoundError struct {
	Key string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("Key %s not found in datastore", e.Key)
}

// AlreadyExistsError indicates the object with the given key already
// exists in the datastore
type AlreadyExistsError struct {
	Key string
}

func (e AlreadyExistsError) Error() string {
	return fmt.Sprintf("Key %s already exists in datastore", e.Key)
}

// InternalError indicates an internal datastore error occurred
type InternalError struct {
	Details string
}

func (e InternalError) Error() string {
	return fmt.Sprintf("Internal datastore error occurred: %s", e.Details)
}

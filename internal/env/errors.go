// Package env provides utilities for environment variable handling.
package env

import "fmt"

// Common error definitions
var (
	// Environment variable errors
	ErrInvalidEnvFormat = fmt.Errorf("invalid environment variable format")

	// Secret management errors
	ErrInvalidSecretURI   = fmt.Errorf("invalid secret URI format")
	ErrInvalidPathSyntax  = fmt.Errorf("invalid JSON path syntax")
	ErrPathNotFound       = fmt.Errorf("JSON path not found in secret value")
	ErrEmptySecretValue   = fmt.Errorf("secret value is empty")
	ErrInvalidSecretValue = fmt.Errorf("invalid secret value format")

	// Key lookup errors
	ErrKeyNotFound      = fmt.Errorf("key not found in secret JSON")
	ErrInvalidKeyPath   = fmt.Errorf("invalid key path in secret")
	ErrNotAnObject      = fmt.Errorf("intermediate path element is not a JSON object")
	ErrNotAnArray       = fmt.Errorf("intermediate path element is not a JSON array")
	ErrIndexOutOfBounds = fmt.Errorf("array index out of bounds")
	ErrUnexpectedError  = fmt.Errorf("unexpected error while processing secret")
)

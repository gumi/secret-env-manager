// Package cmd implements command-line commands for the secret-env-manager
package cmd

import (
	"fmt"

	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"github.com/gumi-tsd/secret-env-manager/internal/logging"
)

// logger is a package-level logger instance used for command logging
var logger = logging.DefaultLogger()

// Helper functions for creating Result monads
// withSuccess wraps a value in a successful Result
// Pure function: Always returns the same output for the same input
func withSuccess[T any](value T) functional.Result[T] {
	return functional.Success(value)
}

// withFailure creates a failure Result with an error message
// Pure function: Always returns the same output for the same input
func withFailure[T any](message string) functional.Result[T] {
	return functional.Failure[T](fmt.Errorf("%s", message))
}

// logInfoMsg logs information message (side effect)
func logInfoMsg(message string) {
	logger.Info("%s", message)
}

// logSuccessInfo logs success information (side effect)
func logSuccessInfo(message string) {
	logger.Success("%s", message)
}

// logWarning logs a warning message (side effect)
func logWarning(message string) {
	logger.Warn("%s", message)
}

// logDebugInfo logs debug information (side effect)
func logDebugInfo(message string) {
	logger.Info("%s", message)
}

// logErrorMsg logs an error message (side effect)
func logErrorMsg(message string) {
	logger.Error("%s", message)
}

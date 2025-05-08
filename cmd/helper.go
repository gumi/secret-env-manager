// Package cmd implements command-line commands for the secret-env-manager
package cmd

import (
	"fmt"

	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"github.com/gumi-tsd/secret-env-manager/internal/logging"
)

// logger is a package-level logger instance used for command logging
var logger = logging.DefaultLogger()

// withSuccess wraps a value in a successful Result
func withSuccess[T any](value T) functional.Result[T] {
	return functional.Success(value)
}

// withFailure creates a failure Result with an error message
func withFailure[T any](message string) functional.Result[T] {
	return functional.Failure[T](fmt.Errorf("%s", message))
}

// logInfoMsg logs information message
func logInfoMsg(message string) {
	logger.Info("%s", message)
}

// logSuccessInfo logs success information
func logSuccessInfo(message string) {
	logger.Success("%s", message)
}

// logWarning logs a warning message
func logWarning(message string) {
	logger.Warn("%s", message)
}

// logDebugInfo logs debug information
func logDebugInfo(message string) {
	logger.Info("%s", message)
}

// logErrorMsg logs an error message
func logErrorMsg(message string) {
	logger.Error("%s", message)
}

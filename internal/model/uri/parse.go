// Package uri provides URI parsing and manipulation for secret references
package uri

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/functional"
)

// URI constants
const (
	URIPrefix           = "sem://"         // Standard prefix for secret URIs
	AwsPlatform         = "aws"            // AWS platform identifier
	GoogleCloudPlatform = "googlecloud"    // Google Cloud platform identifier (to be used consistently instead of GCP)
	AwsDefaultRegion    = "ap-northeast-1" // Default AWS region (Tokyo)
)

// cleanURI removes whitespace and trailing newlines from a URI string.
func cleanURI(uri string) string {
	return strings.TrimRight(strings.TrimSpace(uri), "\r\n")
}

// validatePrefix checks if the URI has the required prefix.
func validatePrefix(uri string) functional.Result[string] {
	if !strings.HasPrefix(uri, URIPrefix) {
		return functional.Failure[string](
			fmt.Errorf("invalid URI: must start with '%s'", URIPrefix))
	}
	return functional.Success(strings.TrimPrefix(uri, URIPrefix))
}

// ParsedPath represents the components extracted from the path part of a URI
type parsedPath struct {
	Platform   string
	Service    string
	Account    string
	SecretName string
}

// ParsedQuery represents the components extracted from the query part of a URI
type parsedQuery struct {
	Version string
	Region  string
	Key     string
}

// splitPathAndQuery separates the path and query parts of a URI.
func splitPathAndQuery(trimmedURI string) (string, string) {
	path, query, _ := strings.Cut(trimmedURI, "?")
	return path, query
}

// parsePathComponents extracts components from the path part of a URI.
func parsePathComponents(path string) functional.Result[parsedPath] {
	parts := strings.SplitN(path, "/", 3)
	if len(parts) != 3 {
		return functional.Failure[parsedPath](
			fmt.Errorf("invalid URI path: expected '<Platform>:<Service>/<Account>/<SecretName>'"))
	}

	platform, service, ok := strings.Cut(parts[0], ":")
	if !ok || platform == "" || service == "" || parts[1] == "" || parts[2] == "" {
		return functional.Failure[parsedPath](
			fmt.Errorf("invalid URI: missing required fields"))
	}

	return functional.Success(parsedPath{
		Platform:   platform,
		Service:    service,
		Account:    parts[1],
		SecretName: parts[2],
	})
}

// parseQueryParams extracts parameters from the query part of a URI.
func parseQueryParams(query string) functional.Result[parsedQuery] {
	if query == "" {
		return functional.Success(parsedQuery{})
	}

	q, err := url.ParseQuery(query)
	if err != nil {
		return functional.Failure[parsedQuery](
			fmt.Errorf("invalid query: %w", err))
	}

	return functional.Success(parsedQuery{
		Version: q.Get("version"),
		Region:  q.Get("region"),
		Key:     q.Get("key"),
	})
}

// determineVersion sets the version based on platform or query parameter.
func determineVersion(platform, queryVersion string) string {
	if queryVersion != "" {
		return queryVersion
	}

	switch platform {
	case AwsPlatform:
		return AwsDefaultVersion
	case GoogleCloudPlatform:
		return GoogleCloudDefaultVersion
	default:
		return ""
	}
}

// determineRegion sets the region based on platform or query parameter.
func determineRegion(platform, queryRegion string) string {
	if queryRegion != "" {
		return queryRegion
	}

	if platform == AwsPlatform {
		return AwsDefaultRegion
	}

	return ""
}

// createSecretURI combines the parsed path and query into a SecretURI
func createSecretURI(path parsedPath, query parsedQuery) SecretURI {
	version := determineVersion(path.Platform, query.Version)
	region := determineRegion(path.Platform, query.Region)

	secretURI := NewSecretURI(path.Platform, path.Service, path.Account, path.SecretName)

	if region != "" {
		secretURI = secretURI.WithRegion(region)
	}

	if version != "" {
		secretURI = secretURI.WithVersion(version)
	}

	if query.Key != "" {
		secretURI = secretURI.WithKey(query.Key)
	}

	return secretURI
}

// splitAndProcess takes a validated URI string and handles path/query splitting and processing
func splitAndProcess(trimmed string) functional.Result[SecretURI] {
	// Extract path and query - a pure operation with no side effects
	path, query := splitPathAndQuery(trimmed)

	// Parse the path first
	return functional.Chain(
		parsePathComponents(path),
		func(parsedPath parsedPath) functional.Result[SecretURI] {
			// Then parse the query and combine results
			return combinePathAndQuery(parsedPath, query)
		},
	)
}

// combinePathAndQuery combines path and query parsing results into a SecretURI
func combinePathAndQuery(path parsedPath, query string) functional.Result[SecretURI] {
	return functional.MapResultTo(
		parseQueryParams(query),
		func(parsedQuery parsedQuery) SecretURI {
			return createSecretURI(path, parsedQuery)
		},
	)
}

// ParseResult is a monadic version of Parse for SecretURI
func ParseResult(uriString string) functional.Result[SecretURI] {
	// Create a clean pipeline of transformations
	// 1. Clean the URI string
	cleanedUri := cleanURI(uriString)

	// 2. Validate the prefix and process if valid
	return functional.Chain(
		validatePrefix(cleanedUri),
		splitAndProcess,
	)
}

// Parse parses a secret URI string and returns its components.
func Parse(uriString string) (SecretURI, error) {
	result := ParseResult(uriString)
	if result.IsFailure() {
		return SecretURI{}, result.GetError()
	}
	return result.Unwrap(), nil
}

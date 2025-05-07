// Package provider supplies interfaces and implementations for retrieving secrets
package provider

import (
	"fmt"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"github.com/gumi-tsd/secret-env-manager/internal/model/env"
	"github.com/gumi-tsd/secret-env-manager/internal/model/uri"
	"github.com/gumi-tsd/secret-env-manager/internal/provider/aws"
	"github.com/gumi-tsd/secret-env-manager/internal/provider/googlecloud"
	"github.com/gumi-tsd/secret-env-manager/internal/secret"
	"github.com/gumi-tsd/secret-env-manager/internal/text"
)

// ProviderConfig contains configuration for creating providers
type ProviderConfig struct {
	EndpointURL  string
	NoExpandJson bool
}

// NewProviderConfig creates a new provider configuration
func NewProviderConfig(endpointURL string) ProviderConfig {
	return ProviderConfig{
		EndpointURL:  endpointURL,
		NoExpandJson: false,
	}
}

// SecretResult represents the result of retrieving a secret
type SecretResult struct {
	Values map[string]string
	Keys   []string
	Error  error
}

// NewSecretResult creates a successful secret result
func NewSecretResult(values map[string]string, keys []string) SecretResult {
	return SecretResult{
		Values: values,
		Keys:   keys,
		Error:  nil,
	}
}

// WithError creates a new SecretResult with an error
func (r SecretResult) WithError(err error) SecretResult {
	return SecretResult{
		Values: r.Values,
		Keys:   r.Keys,
		Error:  err,
	}
}

// IsSuccess checks if the result is successful
func (r SecretResult) IsSuccess() bool {
	return r.Error == nil
}

// Merge combines two SecretResults
func (r SecretResult) Merge(other SecretResult) SecretResult {
	if !r.IsSuccess() {
		return r
	}
	if !other.IsSuccess() {
		return other
	}

	result := make(map[string]string, len(r.Values)+len(other.Values))
	for k, v := range r.Values {
		result[k] = v
	}
	for k, v := range other.Values {
		result[k] = v
	}

	return NewSecretResult(
		result,
		append(r.Keys, other.Keys...),
	)
}

// CreateProviderMap constructs a map of platform identifiers to SecretProviders
func CreateProviderMap(config ProviderConfig) map[string]SecretProvider {
	return map[string]SecretProvider{
		"aws": &awsSecretProvider{
			provider:    aws.NewAwsProvider(),
			endpointURL: config.EndpointURL,
			config:      config,
		},
		"googlecloud": &googleCloudSecretProvider{
			provider: googlecloud.NewGoogleCloudProvider(),
			cache:    make(map[string]string),
			config:   config,
		},
	}
}

// NewAwsSecretProvider creates a new AWS secret provider
func NewAwsSecretProvider(provider *aws.AwsProvider, endpointURL string) SecretProvider {
	config := NewProviderConfig(endpointURL)
	return &awsSecretProvider{
		provider:    provider,
		endpointURL: endpointURL,
		config:      config,
	}
}

// NewGoogleCloudSecretProvider creates a new Google Cloud secret provider
func NewGoogleCloudSecretProvider(provider *googlecloud.GoogleCloudProvider) SecretProvider {
	config := NewProviderConfig("")
	return &googleCloudSecretProvider{
		provider: provider,
		cache:    make(map[string]string),
		config:   config,
	}
}

// AcquireSecretsMapping retrieves secrets for a list of environment entries
func AcquireSecretsMapping(entries []env.Entry, endpointURL string) (map[string]string, []string, error) {
	// この関数はNoExpandJsonのオプションに対応していない
	// クライアントでNoExpandJsonを使用する場合は、直接CreateProviderMapとProcessEntriesResultを使用する必要がある
	config := NewProviderConfig(endpointURL)
	providers := CreateProviderMap(config)
	result := ProcessEntriesResult(entries, providers)

	if !result.IsSuccess() {
		return nil, nil, result.Error
	}

	return result.Values, result.Keys, nil
}

type awsSecretProvider struct {
	provider    *aws.AwsProvider
	endpointURL string
	config      ProviderConfig
}

// GetSecrets retrieves secrets from AWS with optional custom endpoint
func (p *awsSecretProvider) GetSecrets(uri uri.SecretURI) (string, error) {
	if p.endpointURL != "" {
		return p.provider.GetSecretsWithEndpoint(uri, p.endpointURL)
	}
	return p.provider.GetSecrets(uri)
}

// GetSecretsResult retrieves secrets from AWS with Result monad
func (p *awsSecretProvider) GetSecretsResult(uri uri.SecretURI) functional.Result[string] {
	secretValue, err := p.GetSecrets(uri)
	if err != nil {
		return functional.Failure[string](err)
	}
	return functional.Success(secretValue)
}

// GetConfig returns the provider configuration
func (p *awsSecretProvider) GetConfig() ProviderConfig {
	return p.config
}

type googleCloudSecretProvider struct {
	provider *googlecloud.GoogleCloudProvider
	cache    map[string]string
	config   ProviderConfig
}

// GetSecrets retrieves secrets from Google Cloud
func (p *googleCloudSecretProvider) GetSecrets(uri uri.SecretURI) (string, error) {
	return p.provider.GetSecrets(uri)
}

// GetSecretsResult retrieves secrets from Google Cloud with Result monad
func (p *googleCloudSecretProvider) GetSecretsResult(uri uri.SecretURI) functional.Result[string] {
	value, err := p.GetSecrets(uri)
	if err != nil {
		return functional.Failure[string](err)
	}
	return functional.Success(value)
}

// GetConfig returns the provider configuration
func (p *googleCloudSecretProvider) GetConfig() ProviderConfig {
	return p.config
}

// ProcessEntriesResult processes entries and returns a SecretResult
func ProcessEntriesResult(entries []env.Entry, providers map[string]SecretProvider) SecretResult {
	result := NewSecretResult(make(map[string]string), []string{})

	// Get config from first provider to access NoExpandJson setting
	var config ProviderConfig
	for _, p := range providers {
		// 各プロバイダーから設定を取得
		config = p.GetConfig()
		// 設定が見つかったらループを抜ける
		break
	}

	for i, entry := range entries {
		entryResult := ProcessEntryResultWithOptions(i, entry, providers, config.NoExpandJson)
		if !entryResult.IsSuccess() {
			return entryResult
		}

		result = result.Merge(entryResult)
	}

	return result
}

// ProcessEntries processes a list of environment entries and retrieves associated secrets
func ProcessEntries(entries []env.Entry, providers map[string]SecretProvider) (map[string]string, []string, error) {
	result := ProcessEntriesResult(entries, providers)
	if !result.IsSuccess() {
		return nil, nil, result.Error
	}
	return result.Values, result.Keys, nil
}

// ProcessEntryResult processes a single environment entry and returns a SecretResult
// 部分的に純粋な関数: ログ出力以外の副作用はありません
func ProcessEntryResult(idx int, entry env.Entry, providers map[string]SecretProvider) SecretResult {
	// Try to parse the entry as a secret URI
	uriResult := ParseEntryAsSecretURI(entry)

	// If it's not a valid secret URI, handle as a regular entry
	if uriResult.IsFailure() {
		err := uriResult.GetError()
		// ログ出力は副作用なのでこの関数は厳密には純粋関数ではありません
		logSkippedEntry(idx+1, entry.Key, "not a valid secret URI: "+err.Error())
		return NewSecretResult(handleRegularEntry(entry), []string{entry.Key})
	}

	uri := uriResult.Unwrap()

	// Retrieve the secret using the appropriate provider
	secretResult := RetrieveSecretResult(uri, providers)
	if secretResult.IsFailure() {
		err := secretResult.GetError()

		// For unsupported platforms, log and handle as a regular entry
		if strings.HasPrefix(err.Error(), "unsupported platform") {
			logSkippedEntry(idx+1, entry.Key, err.Error())
			return NewSecretResult(handleRegularEntry(entry), []string{entry.Key})
		}

		// Otherwise, return the error
		return NewSecretResult(nil, nil).WithError(
			fmt.Errorf("failed to retrieve secret for line %d: %w", idx+1, err))
	}

	// Process the secret value
	secretValue := secretResult.Unwrap()
	key := DetermineEntryKey(entry)
	vals := ProcessSecret(key, uri, secretValue)

	return NewSecretResult(vals, ExtractKeys(vals))
}

// ProcessEntryResultWithOptions processes a single environment entry with options and returns a SecretResult
func ProcessEntryResultWithOptions(idx int, entry env.Entry, providers map[string]SecretProvider, noExpandJson bool) SecretResult {
	// Try to parse the entry as a secret URI
	uriResult := ParseEntryAsSecretURI(entry)

	// If it's not a valid secret URI, handle as a regular entry
	if uriResult.IsFailure() {
		err := uriResult.GetError()
		// ログ出力は副作用なのでこの関数は厳密には純粋関数ではありません
		logSkippedEntry(idx+1, entry.Key, "not a valid secret URI: "+err.Error())
		return NewSecretResult(handleRegularEntry(entry), []string{entry.Key})
	}

	uri := uriResult.Unwrap()

	// Retrieve the secret using the appropriate provider
	secretResult := RetrieveSecretResult(uri, providers)
	if secretResult.IsFailure() {
		err := secretResult.GetError()

		// For unsupported platforms, log and handle as a regular entry
		if strings.HasPrefix(err.Error(), "unsupported platform") {
			logSkippedEntry(idx+1, entry.Key, err.Error())
			return NewSecretResult(handleRegularEntry(entry), []string{entry.Key})
		}

		// Otherwise, return the error
		return NewSecretResult(nil, nil).WithError(
			fmt.Errorf("failed to retrieve secret for line %d: %w", idx+1, err))
	}

	// Process the secret value with options
	secretValue := secretResult.Unwrap()
	key := DetermineEntryKey(entry)
	vals := ProcessSecretWithOptions(key, uri, secretValue, noExpandJson)

	return NewSecretResult(vals, ExtractKeys(vals))
}

// ProcessEntry processes a single environment entry
func ProcessEntry(idx int, entry env.Entry, providers map[string]SecretProvider) (map[string]string, []string, error) {
	result := ProcessEntryResult(idx, entry, providers)
	if !result.IsSuccess() {
		return nil, nil, result.Error
	}
	return result.Values, result.Keys, nil
}

// ParseEntryAsSecretURI parses an environment entry into a SecretURI
func ParseEntryAsSecretURI(entry env.Entry) functional.Result[uri.SecretURI] {
	v := entry.Key
	if entry.Key != "" && entry.Value != "" {
		v = entry.Value
	}
	return uri.ParseResult(v)
}

// DetermineEntryKey determines the key to use for the environment entry
func DetermineEntryKey(entry env.Entry) string {
	if entry.Key == "" || entry.Value == "" {
		return ""
	}
	return entry.Key
}

// ExtractKeys extracts keys from a map
func ExtractKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// handleRegularEntry processes a regular (non-secret) environment entry
func handleRegularEntry(entry env.Entry) map[string]string {
	if entry.Key != "" && entry.Value != "" {
		return map[string]string{entry.Key: entry.Value}
	}
	return map[string]string{entry.Key: ""}
}

// RetrieveSecretResult retrieves a secret using the appropriate provider with Result monad
func RetrieveSecretResult(uri uri.SecretURI, providers map[string]SecretProvider) functional.Result[string] {
	p, ok := providers[uri.Platform]
	if !ok {
		return functional.Failure[string](fmt.Errorf("unsupported platform '%s'", uri.Platform))
	}

	secretValue, err := p.GetSecrets(uri)
	if err != nil {
		return functional.Failure[string](err)
	}

	return functional.Success(secretValue)
}

// RetrieveSecret retrieves a secret using the appropriate provider
func RetrieveSecret(uri uri.SecretURI, providers map[string]SecretProvider) (string, error) {
	result := RetrieveSecretResult(uri, providers)
	if result.IsFailure() {
		return "", result.GetError()
	}
	return result.Unwrap(), nil
}

// ProcessSecret processes a secret value and returns key-value pairs
func ProcessSecret(entryKey string, uri uri.SecretURI, secretValue string) map[string]string {
	return ProcessSecretWithOptions(entryKey, uri, secretValue, false)
}

// ProcessSecretWithOptions processes a secret value with additional options
func ProcessSecretWithOptions(entryKey string, uri uri.SecretURI, secretValue string, noExpandJson bool) map[string]string {
	// Try to parse as JSON first
	jsonResult := text.ParseJSONMapResult(secretValue)

	if jsonResult.IsSuccess() && !noExpandJson {
		// Only expand JSON if not disabled
		return processJSONSecret(entryKey, uri, jsonResult.Unwrap())
	}

	// If not JSON or JSON expansion is disabled, process as plain text
	return processPlainTextSecret(entryKey, uri, secretValue)
}

// processJSONSecret processes a JSON secret value
func processJSONSecret(entryKey string, uri uri.SecretURI, jsonValues map[string]interface{}) map[string]string {
	result := make(map[string]string, len(jsonValues))
	for jsonKey, jsonValue := range jsonValues {
		finalKey := jsonKey
		if entryKey != "" {
			finalKey = entryKey + "_" + jsonKey
		}

		// Convert the interface{} value to string
		var stringValue string
		switch v := jsonValue.(type) {
		case string:
			stringValue = v
		case nil:
			stringValue = ""
		default:
			// For non-string types, convert to JSON string
			jsonBytes, err := text.MarshalToJSONString(v)
			if err == nil {
				stringValue = jsonBytes
			} else {
				stringValue = fmt.Sprintf("%v", v)
			}
		}

		// Remove any surrounding quotes that might have been added
		result[finalKey] = strings.Trim(stringValue, "'")
	}
	return result
}

// processPlainTextSecret processes a plain text secret value
func processPlainTextSecret(entryKey string, uri uri.SecretURI, secretValue string) map[string]string {
	finalValue := secretValue

	// If a specific key was requested, try to extract it
	if uri.Key != "" {
		valueResult := secret.ParseValueResult(secretValue, uri.Key)
		if valueResult.IsSuccess() {
			finalValue = valueResult.Unwrap()
		}
	}

	// Remove any surrounding single quotes that might have been added
	finalValue = strings.Trim(finalValue, "'")

	finalKey := DetermineFinalKey(entryKey, uri)

	return map[string]string{finalKey: finalValue}
}

// DetermineFinalKey determines the final key to use for the secret
func DetermineFinalKey(entryKey string, uri uri.SecretURI) string {
	if entryKey != "" {
		return entryKey
	}

	if uri.Key != "" {
		return uri.Key
	}

	return uri.SecretName
}

// logSkippedEntry logs information about skipped entries
func logSkippedEntry(lineNum int, key, reason string) {
	fmt.Printf("Line %d skipped: %s (reason: %s)\n", lineNum, key, reason)
}

// Package aws provides functionality for interacting with AWS Secrets Manager.
package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"github.com/gumi-tsd/secret-env-manager/internal/model/uri"
	"github.com/gumi-tsd/secret-env-manager/internal/secret"
)

// SecretRequest encapsulates all parameters needed to retrieve a secret
type SecretRequest struct {
	URI      uri.SecretURI
	Endpoint string
	Ctx      context.Context
}

// NewSecretRequest creates a new SecretRequest with default context
func NewSecretRequest(uri uri.SecretURI) SecretRequest {
	return SecretRequest{
		URI:      uri,
		Endpoint: "",
		Ctx:      context.Background(),
	}
}

// WithEndpoint returns a new SecretRequest with the specified endpoint
func (r SecretRequest) WithEndpoint(endpoint string) SecretRequest {
	return SecretRequest{
		URI:      r.URI,
		Endpoint: endpoint,
		Ctx:      r.Ctx,
	}
}

// WithContext returns a new SecretRequest with the specified context
func (r SecretRequest) WithContext(ctx context.Context) SecretRequest {
	return SecretRequest{
		URI:      r.URI,
		Endpoint: r.Endpoint,
		Ctx:      ctx,
	}
}

// GetCacheKey returns a unique identifier for caching
func (r SecretRequest) GetCacheKey() string {
	return uri.BuildCacheKey(r.URI.Account, r.URI.Service, r.URI.SecretName, r.URI.Version, r.URI.Region)
}

// GetSecrets retrieves secrets using a background context
func (p *AwsProvider) GetSecrets(uri uri.SecretURI) (string, error) {
	req := NewSecretRequest(uri)
	result := p.GetSecretValue(req)

	if result.IsFailure() {
		return "", result.GetError()
	}

	return result.Unwrap(), nil
}

// GetSecretsWithEndpoint retrieves secrets using a background context and a custom endpoint
func (p *AwsProvider) GetSecretsWithEndpoint(uri uri.SecretURI, endpoint string) (string, error) {
	req := NewSecretRequest(uri).WithEndpoint(endpoint)
	result := p.GetSecretValue(req)

	if result.IsFailure() {
		return "", result.GetError()
	}

	return result.Unwrap(), nil
}

// GetSecretValue is a functional implementation that fetches and processes secrets
// これは複合的な関数で、キャッシュチェック、シークレット取得、パース処理を組み合わせています
func (p *AwsProvider) GetSecretValue(req SecretRequest) functional.Result[string] {
	// Check cache first
	cacheKey := req.GetCacheKey()
	cachedSecret := p.GetCachedSecret(cacheKey)

	// If found in cache, parse and return
	if cachedSecret.IsSome() {
		valueResult := secret.ParseValueResult(cachedSecret.Unwrap(), req.URI.Key)
		if valueResult.IsFailure() {
			// キーが見つからない場合のエラーをより詳細なメッセージに変換
			err := valueResult.GetError()
			if strings.Contains(err.Error(), "key not found") {
				return functional.Failure[string](
					fmt.Errorf("指定されたキー '%s' がシークレット '%s' に存在しません: %w",
						req.URI.Key, req.URI.SecretName, err))
			}
			return valueResult
		}
		return valueResult
	}

	// Fetch the secret if not in cache
	secretResult := p.RetrieveSecret(req)
	if secretResult.IsFailure() {
		return secretResult
	}

	// Store in cache
	secretVal := secretResult.Unwrap()
	p.CacheSecret(cacheKey, secretVal)

	// Parse and return the value
	valueResult := secret.ParseValueResult(secretVal, req.URI.Key)
	if valueResult.IsFailure() {
		// キーが見つからない場合のエラーをより詳細なメッセージに変換
		err := valueResult.GetError()
		if strings.Contains(err.Error(), "key not found") {
			return functional.Failure[string](
				fmt.Errorf("指定されたキー '%s' がシークレット '%s' に存在しません: %w",
					req.URI.Key, req.URI.SecretName, err))
		}
		return valueResult
	}
	return valueResult
}

// RetrieveSecret fetches a secret from AWS Secrets Manager
// 副作用のある関数: AWS APIを呼び出します
func (p *AwsProvider) RetrieveSecret(req SecretRequest) functional.Result[string] {
	// Get or create client
	client, err := p.GetClient(req.Ctx, req.URI.Account, req.URI.Region, req.Endpoint)
	if err != nil {
		return functional.Failure[string](
			fmt.Errorf("failed to get AWS client - account: %s, region: %s: %w",
				req.URI.Account, req.URI.Region, err))
	}

	// Fetch secret
	secretResult := FetchSecret(req.Ctx, client, req.URI)
	if secretResult.IsFailure() {
		return functional.Failure[string](
			fmt.Errorf("failed to retrieve secret [%s/%s] - account: %s, region: %s: %w",
				req.URI.Service, req.URI.SecretName, req.URI.Account, req.URI.Region, secretResult.GetError()))
	}

	return secretResult
}

// FetchSecret calls AWS Secrets Manager API to get a secret value
func FetchSecret(ctx context.Context, client *secretsmanager.Client, uri uri.SecretURI) functional.Result[string] {
	// Create input for GetSecretValue
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(uri.SecretName),
		VersionStage: aws.String(uri.Version),
	}

	// Use proper logging format instead of fmt.Println
	secret.LogInfoMsg(fmt.Sprintf("Accessing secret: %s", uri.GetUri()))

	// Call AWS API
	result, err := client.GetSecretValue(ctx, input)
	if err != nil {
		return functional.Failure[string](
			fmt.Errorf("AWS Secrets Manager API error [%s] - version: %s, region: %s: %w",
				uri.SecretName, uri.Version, uri.Region, err))
	}

	// Validate response
	if result.SecretString == nil {
		return functional.Failure[string](
			fmt.Errorf("empty secret value [%s] - version: %s, region: %s",
				uri.SecretName, uri.Version, uri.Region))
	}

	// Return the secret string
	return functional.Success(*result.SecretString)
}

// Package env provides environment variable related models and utilities
package env

import (
	"github.com/gumi-tsd/secret-env-manager/internal/functional"
	"github.com/gumi-tsd/secret-env-manager/internal/model/uri"
)

// Variable represents a resolved environment variable that can be either a plain value
// or a reference to a secret managed by a cloud provider
type Variable struct {
	Key   string                           // Environment variable name
	Value string                           // Either a plain value or resolved secret value
	URI   functional.Option[uri.SecretURI] // If this variable was resolved from a secret, contains the URI
}

// NewVariable creates a new environment Variable with the specified key and value
func NewVariable(key, value string) Variable {
	return Variable{
		Key:   key,
		Value: value,
		URI:   functional.None[uri.SecretURI](),
	}
}

// NewSecretVariable creates a new environment Variable resolved from a secret URI
func NewSecretVariable(key, value string, secretURI uri.SecretURI) Variable {
	return Variable{
		Key:   key,
		Value: value,
		URI:   functional.Some(secretURI),
	}
}

// WithValue returns a new Variable with the specified value
func (v Variable) WithValue(value string) Variable {
	return Variable{
		Key:   v.Key,
		Value: value,
		URI:   v.URI,
	}
}

// WithKey returns a new Variable with the specified key
func (v Variable) WithKey(key string) Variable {
	return Variable{
		Key:   key,
		Value: v.Value,
		URI:   v.URI,
	}
}

// WithURI returns a new Variable with the specified URI
func (v Variable) WithURI(uri uri.SecretURI) Variable {
	return Variable{
		Key:   v.Key,
		Value: v.Value,
		URI:   functional.Some(uri),
	}
}

// IsSecret returns true if this variable is resolved from a secret
func (v Variable) IsSecret() bool {
	return v.URI.IsSome()
}

// IsPlainValue returns true if this variable is a plain text value
func (v Variable) IsPlainValue() bool {
	return !v.IsSecret()
}

// AsOption converts a Variable to an Option
func (v Variable) AsOption() functional.Option[Variable] {
	if v.Key == "" {
		return functional.None[Variable]()
	}
	return functional.Some(v)
}

// Map applies a function to transform a Variable
func (v Variable) Map(f func(Variable) Variable) Variable {
	return f(v)
}

// FilterVariables returns a new slice containing only variables that satisfy the predicate
func FilterVariables(variables []Variable, predicate func(Variable) bool) []Variable {
	var result []Variable
	for _, variable := range variables {
		if predicate(variable) {
			result = append(result, variable)
		}
	}
	return result
}

// MapVariables applies a function to each variable in a slice and returns a new slice
func MapVariables(variables []Variable, f func(Variable) Variable) []Variable {
	result := make([]Variable, len(variables))
	for i, variable := range variables {
		result[i] = f(variable)
	}
	return result
}

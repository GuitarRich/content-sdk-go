package graphql

import (
	"net/http"
	"time"
)

// ClientFactory creates GraphQL clients
type ClientFactory interface {
	Create(config ServiceConfig) (Client, error)
}

// ServiceConfig contains configuration for creating a GraphQL client
type ServiceConfig struct {
	// Endpoint is the GraphQL endpoint URL
	Endpoint string

	// APIKey is the Sitecore API key
	APIKey string

	// Retries is the number of retry attempts
	Retries int

	// Timeout is the request timeout
	Timeout time.Duration

	// HTTPClient is an optional custom HTTP client
	HTTPClient *http.Client

	// Headers are custom headers to include in requests
	Headers map[string]string
}

// DefaultClientFactory is the default implementation of ClientFactory
type DefaultClientFactory struct{}

// NewClientFactory creates a new client factory
func NewClientFactory() ClientFactory {
	return &DefaultClientFactory{}
}

// Create creates a new GraphQL client from the service configuration
func (f *DefaultClientFactory) Create(config ServiceConfig) (Client, error) {
	clientConfig := &ClientConfig{
		Retries:    config.Retries,
		Timeout:    config.Timeout,
		RetryDelay: 1 * time.Second,
		Headers:    config.Headers,
	}

	// Apply defaults if not specified
	if clientConfig.Retries == 0 {
		clientConfig.Retries = 3
	}
	if clientConfig.Timeout == 0 {
		clientConfig.Timeout = 30 * time.Second
	}
	if clientConfig.Headers == nil {
		clientConfig.Headers = make(map[string]string)
	}

	return NewClient(config.Endpoint, config.APIKey, config.HTTPClient, clientConfig), nil
}

// CreateGraphQLClient is a convenience function to create a GraphQL client
func CreateGraphQLClient(endpoint, apiKey string, httpClient *http.Client) Client {
	factory := NewClientFactory()
	client, _ := factory.Create(ServiceConfig{
		Endpoint:   endpoint,
		APIKey:     apiKey,
		HTTPClient: httpClient,
	})
	return client
}

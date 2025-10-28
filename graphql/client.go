package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/content-sdk-go/debug"
	"github.com/content-sdk-go/models"
)

// Client is the interface for GraphQL operations
type Client interface {
	Request(ctx context.Context, query string, variables map[string]interface{}) (map[string]interface{}, error)
}

// ClientImpl is the default implementation of the GraphQL client
type ClientImpl struct {
	endpoint   string
	apiKey     string
	httpClient *http.Client
	config     *ClientConfig
}

// ClientConfig contains configuration for the GraphQL client
type ClientConfig struct {
	// Retries is the number of retry attempts for failed requests
	Retries int

	// Timeout is the request timeout duration
	Timeout time.Duration

	// RetryDelay is the base delay between retries (exponential backoff)
	RetryDelay time.Duration

	// Headers are custom headers to include in requests
	Headers map[string]string
}

// DefaultClientConfig returns the default client configuration
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Retries:    3,
		Timeout:    30 * time.Second,
		RetryDelay: 1 * time.Second,
		Headers:    make(map[string]string),
	}
}

// NewClient creates a new GraphQL client
func NewClient(endpoint, apiKey string, httpClient *http.Client, config *ClientConfig) Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		}
	}

	if config == nil {
		config = DefaultClientConfig()
	}

	// Override HTTP client timeout if config specifies it
	if config.Timeout > 0 {
		httpClient.Timeout = config.Timeout
	}

	return &ClientImpl{
		endpoint:   endpoint,
		apiKey:     apiKey,
		httpClient: httpClient,
		config:     config,
	}
}

// Request executes a GraphQL query with retry logic
func (c *ClientImpl) Request(
	ctx context.Context,
	query string,
	variables map[string]interface{},
) (map[string]interface{}, error) {
	var lastErr error

	debug.Common("Requesting GraphQL query: %s", query)
	debug.Common("Variables: %+v", variables)
	debug.Common("Timeout: %v", c.config.Timeout)
	debug.Common("Retries: %d", c.config.Retries)
	debug.Common("RetryDelay: %v", c.config.RetryDelay)
	debug.Common("Headers: %+v", c.config.Headers)
	debug.Common("Endpoint: %s", c.endpoint)
	debug.Common("API Key: %s", c.apiKey)

	// Add context timeout if not already set
	if _, hasDeadline := ctx.Deadline(); !hasDeadline && c.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.config.Timeout)
		defer cancel()
	}

	// Retry loop with exponential backoff
	for attempt := 0; attempt <= c.config.Retries; attempt++ {
		if attempt > 0 {
			// Calculate exponential backoff delay
			delay := time.Duration(math.Pow(2, float64(attempt-1))) * c.config.RetryDelay
			debug.Http("retrying GraphQL request (attempt %d/%d) after %v", attempt, c.config.Retries, delay)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		result, err := c.doRequest(ctx, query, variables)
		if err == nil {
			return result, nil
		}

		lastErr = err

		// Don't retry on context cancellation or validation errors
		if ctx.Err() != nil {
			break
		}
		if _, ok := err.(*models.ValidationError); ok {
			break
		}
	}

	return nil, fmt.Errorf("GraphQL request failed after %d retries: %w", c.config.Retries, lastErr)
}

// doRequest performs a single GraphQL request
func (c *ClientImpl) doRequest(
	ctx context.Context,
	query string,
	variables map[string]interface{},
) (map[string]interface{}, error) {
	// Prepare request body
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GraphQL request: %w", err)
	}

	debug.Http("Request body: %s", string(jsonData))

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Only set sc_apikey header for local API (not Edge API)
	// Edge API uses sitecoreContextId as a query parameter in the URL
	if c.apiKey != "" && !isEdgeAPI(c.endpoint) {
		req.Header.Set("sc_apikey", c.apiKey)
	}

	// Apply custom headers
	for key, value := range c.config.Headers {
		req.Header.Set(key, value)
	}

	// Execute request
	debug.Http("GraphQL request to %s", c.endpoint)
	debug.Http("Request headers: %+v", req.Header)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GraphQL request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("GraphQL request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result struct {
		Data   map[string]interface{} `json:"data"`
		Errors []models.GraphQLError  `json:"errors,omitempty"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GraphQL response: %w", err)
	}

	// Check for GraphQL errors
	if len(result.Errors) > 0 {
		// Return the first error (could be enhanced to return all)
		return nil, &result.Errors[0]
	}

	return result.Data, nil
}

// isEdgeAPI checks if the endpoint is using Edge API (contains sitecoreContextId query parameter)
func isEdgeAPI(endpoint string) bool {
	return strings.Contains(endpoint, "sitecoreContextId=")
}

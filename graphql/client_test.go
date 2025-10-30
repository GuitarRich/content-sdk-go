package graphql

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClient_Request_Success(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type: application/json")
		}
		if r.Header.Get("sc_apikey") != "test-key" {
			t.Errorf("expected sc_apikey header")
		}

		// Return mock response
		response := map[string]any{
			"data": map[string]any{
				"layout": map[string]any{
					"item": map[string]any{
						"name": "Home",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client := NewClient(server.URL, "test-key", nil, DefaultClientConfig())

	// Execute request
	result, err := client.Request(context.Background(), "query { layout { item { name } } }", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	// Verify result
	layout := result["layout"].(map[string]any)
	item := layout["item"].(map[string]any)
	if item["name"] != "Home" {
		t.Errorf("expected name=Home, got %v", item["name"])
	}
}

func TestClient_Request_GraphQLError(t *testing.T) {
	// Mock server that returns GraphQL error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"data": nil,
			"errors": []map[string]any{
				{
					"message": "Field 'nonexistent' not found",
					"path":    []any{"layout", "nonexistent"},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key", nil, DefaultClientConfig())

	result, err := client.Request(context.Background(), "query { layout { nonexistent } }", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result on error, got %v", result)
	}

	// Check that we got an error (the exact message may vary due to retries)
	if !strings.Contains(err.Error(), "Field 'nonexistent' not found") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestClient_Request_HTTPError(t *testing.T) {
	// Mock server that returns HTTP error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key", nil, DefaultClientConfig())

	result, err := client.Request(context.Background(), "query { layout }", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result on error, got %v", result)
	}
}

func TestClient_Request_Retry(t *testing.T) {
	attempts := 0
	// Mock server that fails first 2 times, succeeds on 3rd
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		response := map[string]any{
			"data": map[string]any{
				"success": true,
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &ClientConfig{
		Retries:    3,
		Timeout:    5 * time.Second,
		RetryDelay: 10 * time.Millisecond, // Short delay for testing
		Headers:    make(map[string]string),
	}

	client := NewClient(server.URL, "test-key", nil, config)

	result, err := client.Request(context.Background(), "query { test }", nil)

	if err != nil {
		t.Fatalf("unexpected error after retries: %v", err)
	}

	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}

	if result == nil {
		t.Fatal("expected result after retry, got nil")
	}
}

func TestClient_Request_ContextTimeout(t *testing.T) {
	// Mock server with slow response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		response := map[string]any{
			"data": map[string]any{},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &ClientConfig{
		Retries:    0,
		Timeout:    50 * time.Millisecond,
		RetryDelay: 10 * time.Millisecond,
		Headers:    make(map[string]string),
	}

	client := NewClient(server.URL, "test-key", nil, config)

	ctx := context.Background()
	result, err := client.Request(ctx, "query { test }", nil)

	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result on timeout, got %v", result)
	}
}

func TestClientFactory_Create(t *testing.T) {
	factory := NewClientFactory()

	config := ServiceConfig{
		Endpoint: "https://test.example.com/graphql",
		APIKey:   "test-key",
		Retries:  5,
		Timeout:  10 * time.Second,
	}

	client, err := factory.Create(config)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client == nil {
		t.Fatal("expected client, got nil")
	}

	// Verify client has correct configuration
	impl, ok := client.(*ClientImpl)
	if !ok {
		t.Fatal("expected ClientImpl type")
	}

	if impl.endpoint != config.Endpoint {
		t.Errorf("expected endpoint %s, got %s", config.Endpoint, impl.endpoint)
	}

	if impl.apiKey != config.APIKey {
		t.Errorf("expected apiKey %s, got %s", config.APIKey, impl.apiKey)
	}

	if impl.config.Retries != config.Retries {
		t.Errorf("expected %d retries, got %d", config.Retries, impl.config.Retries)
	}

	if impl.config.Timeout != config.Timeout {
		t.Errorf("expected timeout %v, got %v", config.Timeout, impl.config.Timeout)
	}
}

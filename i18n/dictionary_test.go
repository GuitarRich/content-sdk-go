package i18n

import (
	"context"
	"testing"
)

// Mock GraphQL client for testing
type mockGraphQLClient struct {
	response map[string]any
	err      error
}

func (m *mockGraphQLClient) Request(ctx context.Context, query string, variables map[string]any) (map[string]any, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

func TestDictionaryService_FetchDictionaryData_Success(t *testing.T) {
	mockResponse := map[string]any{
		"site": map[string]any{
			"siteInfo": map[string]any{
				"dictionary": []any{
					map[string]any{
						"key":   "welcome",
						"value": "Welcome",
					},
					map[string]any{
						"key":   "goodbye",
						"value": "Goodbye",
					},
					map[string]any{
						"key":   "hello",
						"value": "Hello",
					},
				},
			},
		},
	}

	mockClient := &mockGraphQLClient{
		response: mockResponse,
	}

	service := NewDictionaryService(DictionaryServiceConfig{
		GraphQLClient: mockClient,
		SiteName:      "testsite",
	})

	phrases, err := service.FetchDictionaryData(context.Background(), "en", "testsite")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(phrases) != 3 {
		t.Errorf("expected 3 phrases, got %d", len(phrases))
	}

	if phrases["welcome"] != "Welcome" {
		t.Errorf("expected 'Welcome', got '%s'", phrases["welcome"])
	}

	if phrases["goodbye"] != "Goodbye" {
		t.Errorf("expected 'Goodbye', got '%s'", phrases["goodbye"])
	}

	if phrases["hello"] != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", phrases["hello"])
	}
}

func TestDictionaryService_FetchDictionaryData_EmptyResponse(t *testing.T) {
	mockResponse := map[string]any{
		"site": map[string]any{
			"siteInfo": map[string]any{
				"dictionary": []any{},
			},
		},
	}

	mockClient := &mockGraphQLClient{
		response: mockResponse,
	}

	service := NewDictionaryService(DictionaryServiceConfig{
		GraphQLClient: mockClient,
		SiteName:      "testsite",
	})

	phrases, err := service.FetchDictionaryData(context.Background(), "en", "testsite")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(phrases) != 0 {
		t.Errorf("expected 0 phrases, got %d", len(phrases))
	}
}

func TestDictionaryService_FetchDictionaryData_NoSite(t *testing.T) {
	mockResponse := map[string]any{}

	mockClient := &mockGraphQLClient{
		response: mockResponse,
	}

	service := NewDictionaryService(DictionaryServiceConfig{
		GraphQLClient: mockClient,
		SiteName:      "testsite",
	})

	phrases, err := service.FetchDictionaryData(context.Background(), "en", "testsite")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(phrases) != 0 {
		t.Errorf("expected 0 phrases, got %d", len(phrases))
	}
}

func TestDictionaryService_GetDictionaryQuery(t *testing.T) {
	service := &dictionaryServiceImpl{
		siteName: "testsite",
	}

	query := service.getDictionaryQuery("mysite", "fr")

	if query == "" {
		t.Fatal("expected non-empty query")
	}

	// Verify query contains site and locale
	expectedSite := `site: "mysite"`
	expectedLocale := `language: "fr"`

	if !contains(query, expectedSite) {
		t.Errorf("query should contain '%s'", expectedSite)
	}

	if !contains(query, expectedLocale) {
		t.Errorf("query should contain '%s'", expectedLocale)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

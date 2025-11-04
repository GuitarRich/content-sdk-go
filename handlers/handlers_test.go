package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/guitarrich/content-sdk-go/middleware"
)

// MockContext for testing handlers
type MockContext struct {
	request  *http.Request
	response *httptest.ResponseRecorder
	path     string
	values   map[string]any
}

func NewMockContext(method, path string, body []byte) *MockContext {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewBuffer(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	resp := httptest.NewRecorder()

	return &MockContext{
		request:  req,
		response: resp,
		path:     path,
		values:   make(map[string]any),
	}
}

func (m *MockContext) Request() *http.Request        { return m.request }
func (m *MockContext) Response() http.ResponseWriter { return m.response }
func (m *MockContext) Path() string                  { return m.path }
func (m *MockContext) SetPath(path string)           { m.path = path }
func (m *MockContext) Get(key string) any            { return m.values[key] }
func (m *MockContext) Set(key string, val any)       { m.values[key] = val }
func (m *MockContext) Cookie(name string) (*http.Cookie, error) {
	return m.request.Cookie(name)
}
func (m *MockContext) SetCookie(cookie *http.Cookie) {
	m.response.Header().Add("Set-Cookie", cookie.String())
}
func (m *MockContext) Header(key string) string { return m.request.Header.Get(key) }
func (m *MockContext) SetHeader(key, value string) {
	m.response.Header().Set(key, value)
}
func (m *MockContext) Redirect(code int, url string) error {
	http.Redirect(m.response, m.request, url, code)
	return nil
}
func (m *MockContext) String(code int, s string) error {
	m.response.WriteHeader(code)
	m.response.WriteString(s)
	return nil
}
func (m *MockContext) JSON(code int, i any) error {
	m.response.WriteHeader(code)
	m.response.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(m.response).Encode(i)
}
func (m *MockContext) NoContent(code int) error {
	m.response.WriteHeader(code)
	return nil
}

// MockComponentRegistry is a mock implementation for testing
type MockComponentRegistry struct {
	components []string
}

func (m *MockComponentRegistry) List() []string {
	return m.components
}

func TestEditingConfigHandler(t *testing.T) {
	// Create mock registry with sample components
	mockRegistry := &MockComponentRegistry{
		components: []string{
			"Hero",
			"ProductListing",
			"RichTextBlock",
			"PromoBlock",
			"Unknown", // Should be filtered out
			"Container",
		},
	}

	handler := NewEditingConfigHandler(mockRegistry)
	ctx := NewMockContext("GET", "/api/editing/config", nil)

	err := handler.Handle(ctx)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ctx.response.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", ctx.response.Code)
	}

	// Parse response
	var response EditingConfigResponse
	if err := json.Unmarshal(ctx.response.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}

	// Verify components are present and "Unknown" is filtered out
	foundHero := false
	foundUnknown := false
	for _, component := range response.Components {
		if component == "Hero" {
			foundHero = true
		}
		if component == "Unknown" {
			foundUnknown = true
		}
	}

	if !foundHero {
		t.Errorf("expected Hero component in list")
	}

	if foundUnknown {
		t.Errorf("Unknown component should be filtered out")
	}

	// Verify editMode
	if response.EditMode != "metadata" {
		t.Errorf("expected editMode 'metadata', got '%s'", response.EditMode)
	}

	// Verify packages are present
	if len(response.Packages) == 0 {
		t.Errorf("expected packages to be present")
	}

	// Verify specific package versions
	if response.Packages["@sitecore-content-sdk/core"] != "1.1.0" {
		t.Errorf("expected @sitecore-content-sdk/core version 1.1.0, got %s",
			response.Packages["@sitecore-content-sdk/core"])
	}

	if response.Packages["@sitecore/components"] != "2.1.0" {
		t.Errorf("expected @sitecore/components version 2.1.0, got %s",
			response.Packages["@sitecore/components"])
	}
}

func TestCatchAllHandler_GetSiteFromContext(t *testing.T) {
	handler := &CatchAllHandler{}
	ctx := NewMockContext("GET", "/test", nil)

	// Set site in context
	ctx.Set(middleware.SiteKey, "mysite")

	site := handler.getSiteFromContext(ctx)

	if site != "mysite" {
		t.Errorf("expected 'mysite', got '%s'", site)
	}
}

func TestCatchAllHandler_GetLocaleFromContext(t *testing.T) {
	handler := &CatchAllHandler{}
	ctx := NewMockContext("GET", "/test", nil)

	// Set locale in context
	ctx.Set(middleware.LocaleKey, "fr")

	locale := handler.getLocaleFromContext(ctx)

	if locale != "fr" {
		t.Errorf("expected 'fr', got '%s'", locale)
	}
}

func TestCatchAllHandler_GetLocaleFromContext_Default(t *testing.T) {
	handler := &CatchAllHandler{}
	ctx := NewMockContext("GET", "/test", nil)

	// Don't set locale in context
	locale := handler.getLocaleFromContext(ctx)

	if locale != "en" {
		t.Errorf("expected default 'en', got '%s'", locale)
	}
}

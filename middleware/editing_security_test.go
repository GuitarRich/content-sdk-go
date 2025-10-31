package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestEditingSecurityMiddleware_ValidSecret(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/editing/config?secret=test-secret", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware with test secret
	middleware := EditingSecurityMiddleware(EditingSecurityConfig{
		Secret:         "test-secret",
		AllowedOrigins: []string{"https://example.com"},
	})

	// Create a test handler that the middleware will call
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "success", rec.Body.String())
	// Verify iframe headers are set
	assert.Equal(t, "frame-ancestors https://example.com", rec.Header().Get("Content-Security-Policy"))
}

func TestEditingSecurityMiddleware_InvalidSecret(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/editing/config?secret=wrong-secret", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware with test secret
	middleware := EditingSecurityMiddleware(EditingSecurityConfig{
		Secret:         "test-secret",
		AllowedOrigins: []string{"https://example.com"},
	})

	// Create a test handler
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert - middleware should return JSON error
	assert.NoError(t, err) // The middleware returns the error as JSON, not as error
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid editing secret")
}

func TestEditingSecurityMiddleware_MissingSecret(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/editing/config", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware with test secret
	middleware := EditingSecurityMiddleware(EditingSecurityConfig{
		Secret:         "test-secret",
		AllowedOrigins: []string{"https://example.com"},
	})

	// Create a test handler
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "editing secret is required")
}

func TestEditingSecurityMiddleware_SkipSecretValidation(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/editing/config", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware with skip validation
	middleware := EditingSecurityMiddleware(EditingSecurityConfig{
		Secret:               "test-secret",
		AllowedOrigins:       []string{"https://example.com"},
		SkipSecretValidation: true,
	})

	// Create a test handler
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert - should succeed without secret
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "success", rec.Body.String())
}

func TestEditingSecurityMiddleware_CORS_AllowedOrigin(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/editing/config?secret=test-secret", nil)
	req.Header.Set("Origin", "https://example.com")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware
	middleware := EditingSecurityMiddleware(EditingSecurityConfig{
		Secret:         "test-secret",
		AllowedOrigins: []string{"https://example.com"},
	})

	// Create a test handler
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "https://example.com", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
}

func TestEditingSecurityMiddleware_CORS_DisallowedOrigin(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/editing/config?secret=test-secret", nil)
	req.Header.Set("Origin", "https://malicious.com")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware
	middleware := EditingSecurityMiddleware(EditingSecurityConfig{
		Secret:         "test-secret",
		AllowedOrigins: []string{"https://example.com"},
	})

	// Create a test handler
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert - request should succeed but CORS headers should not be set
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Header().Get("Access-Control-Allow-Origin"))
}

func TestEditingSecurityMiddleware_CORS_Preflight_Allowed(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodOptions, "/api/editing/config", nil)
	req.Header.Set("Origin", "https://example.com")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware
	middleware := EditingSecurityMiddleware(EditingSecurityConfig{
		Secret:         "test-secret",
		AllowedOrigins: []string{"https://example.com"},
	})

	// Create a test handler
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "https://example.com", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), "POST")
}

func TestEditingSecurityMiddleware_CORS_Preflight_Disallowed(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodOptions, "/api/editing/config", nil)
	req.Header.Set("Origin", "https://malicious.com")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware
	middleware := EditingSecurityMiddleware(EditingSecurityConfig{
		Secret:         "test-secret",
		AllowedOrigins: []string{"https://example.com"},
	})

	// Create a test handler
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestEditingSecurityMiddleware_CORS_Wildcard(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/editing/config?secret=test-secret", nil)
	req.Header.Set("Origin", "https://any-origin.com")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware with wildcard
	middleware := EditingSecurityMiddleware(EditingSecurityConfig{
		Secret:         "test-secret",
		AllowedOrigins: []string{"*"},
	})

	// Create a test handler
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "https://any-origin.com", rec.Header().Get("Access-Control-Allow-Origin"))
}

func TestEditingSecurityMiddleware_CORS_EmptyAllowedOrigins(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/editing/config?secret=test-secret", nil)
	req.Header.Set("Origin", "https://any-origin.com")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware with empty allowed origins (development mode)
	middleware := EditingSecurityMiddleware(EditingSecurityConfig{
		Secret:         "test-secret",
		AllowedOrigins: []string{},
	})

	// Create a test handler
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert - should allow all origins in development mode
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "https://any-origin.com", rec.Header().Get("Access-Control-Allow-Origin"))
	// Verify iframe wildcard in development mode
	assert.Equal(t, "frame-ancestors *", rec.Header().Get("Content-Security-Policy"))
}

func TestEditingSecurityMiddleware_Iframe_MultipleOrigins(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/editing/config?secret=test-secret", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware with multiple allowed origins
	middleware := EditingSecurityMiddleware(EditingSecurityConfig{
		Secret:         "test-secret",
		AllowedOrigins: []string{"https://pages.sitecorecloud.io", "https://pages-eu.sitecorecloud.io"},
	})

	// Create a test handler
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	// Verify iframe headers contain all origins
	csp := rec.Header().Get("Content-Security-Policy")
	assert.Contains(t, csp, "frame-ancestors")
	assert.Contains(t, csp, "https://pages.sitecorecloud.io")
	assert.Contains(t, csp, "https://pages-eu.sitecorecloud.io")
}

func TestEditingSecurityMiddleware_Iframe_Wildcard(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/editing/config?secret=test-secret", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware with wildcard
	middleware := EditingSecurityMiddleware(EditingSecurityConfig{
		Secret:         "test-secret",
		AllowedOrigins: []string{"*"},
	})

	// Create a test handler
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	// Verify iframe wildcard
	assert.Equal(t, "frame-ancestors *", rec.Header().Get("Content-Security-Policy"))
}

func TestEditingSecurityMiddleware_Iframe_SingleOrigin(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/editing/config?secret=test-secret", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create middleware with single origin
	middleware := EditingSecurityMiddleware(EditingSecurityConfig{
		Secret:         "test-secret",
		AllowedOrigins: []string{"https://pages.sitecorecloud.io"},
	})

	// Create a test handler
	handler := middleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	// Execute
	err := handler(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	// Verify iframe header with single origin
	assert.Equal(t, "frame-ancestors https://pages.sitecorecloud.io", rec.Header().Get("Content-Security-Policy"))
}

func TestIsOriginAllowed(t *testing.T) {
	tests := []struct {
		name           string
		origin         string
		allowedOrigins []string
		expected       bool
	}{
		{
			name:           "Exact match",
			origin:         "https://example.com",
			allowedOrigins: []string{"https://example.com", "https://other.com"},
			expected:       true,
		},
		{
			name:           "Not in list",
			origin:         "https://malicious.com",
			allowedOrigins: []string{"https://example.com", "https://other.com"},
			expected:       false,
		},
		{
			name:           "Wildcard",
			origin:         "https://any-origin.com",
			allowedOrigins: []string{"*"},
			expected:       true,
		},
		{
			name:           "Empty list allows all",
			origin:         "https://any-origin.com",
			allowedOrigins: []string{},
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOriginAllowed(tt.origin, tt.allowedOrigins)
			assert.Equal(t, tt.expected, result)
		})
	}
}


package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockContext is a simple mock implementation of Context for testing
type MockContext struct {
	request  *http.Request
	response *httptest.ResponseRecorder
	path     string
	values   map[string]any
	headers  map[string]string
}

func NewMockContext(method, path string) *MockContext {
	req := httptest.NewRequest(method, path, nil)
	resp := httptest.NewRecorder()

	return &MockContext{
		request:  req,
		response: resp,
		path:     path,
		values:   make(map[string]any),
		headers:  make(map[string]string),
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
	m.headers[key] = value
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
	return nil
}
func (m *MockContext) NoContent(code int) error {
	m.response.WriteHeader(code)
	return nil
}

func TestMiddlewareFunc_Handle(t *testing.T) {
	called := false

	mw := MiddlewareFunc(func(ctx Context, next HandlerFunc) error {
		called = true
		return next(ctx)
	})

	ctx := NewMockContext("GET", "/test")

	err := mw.Handle(ctx, func(ctx Context) error {
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !called {
		t.Error("middleware was not called")
	}
}

func TestChain(t *testing.T) {
	var order []int

	mw1 := MiddlewareFunc(func(ctx Context, next HandlerFunc) error {
		order = append(order, 1)
		return next(ctx)
	})

	mw2 := MiddlewareFunc(func(ctx Context, next HandlerFunc) error {
		order = append(order, 2)
		return next(ctx)
	})

	mw3 := MiddlewareFunc(func(ctx Context, next HandlerFunc) error {
		order = append(order, 3)
		return next(ctx)
	})

	chained := Chain(mw1, mw2, mw3)

	ctx := NewMockContext("GET", "/test")

	err := chained.Handle(ctx, func(ctx Context) error {
		order = append(order, 4)
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify execution order
	expected := []int{1, 2, 3, 4}
	if len(order) != len(expected) {
		t.Errorf("expected %d calls, got %d", len(expected), len(order))
	}

	for i, v := range expected {
		if order[i] != v {
			t.Errorf("expected order[%d]=%d, got %d", i, v, order[i])
		}
	}
}

func TestContext_SetAndGet(t *testing.T) {
	ctx := NewMockContext("GET", "/test")

	// Set a value
	ctx.Set("key", "value")

	// Get the value
	val := ctx.Get("key")

	if val != "value" {
		t.Errorf("expected 'value', got %v", val)
	}
}

func TestContext_PathOperations(t *testing.T) {
	ctx := NewMockContext("GET", "/original")

	// Check initial path
	if ctx.Path() != "/original" {
		t.Errorf("expected '/original', got %s", ctx.Path())
	}

	// Set new path
	ctx.SetPath("/rewritten")

	if ctx.Path() != "/rewritten" {
		t.Errorf("expected '/rewritten', got %s", ctx.Path())
	}
}

func TestHealthcheckMiddleware(t *testing.T) {
	config := HealthcheckConfig{
		Path: "/health",
		Response: map[string]string{
			"status": "healthy",
		},
	}

	mw := NewHealthcheckMiddleware(config)
	mockCtx := NewMockContext("GET", "/health")

	err := mw.Handle(mockCtx, func(c Context) error {
		t.Error("next handler should not be called for healthcheck path")
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify response code
	if mockCtx.response.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", mockCtx.response.Code)
	}
}

func TestHealthcheckMiddleware_PassThrough(t *testing.T) {
	config := HealthcheckConfig{
		Path: "/health",
	}

	mw := NewHealthcheckMiddleware(config)
	mockCtx := NewMockContext("GET", "/other-path")

	nextCalled := false

	err := mw.Handle(mockCtx, func(c Context) error {
		nextCalled = true
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !nextCalled {
		t.Error("next handler should be called for non-healthcheck paths")
	}
}

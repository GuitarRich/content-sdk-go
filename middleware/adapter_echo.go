package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// EchoContext wraps an Echo context to implement our Context interface
type EchoContext struct {
	echo.Context
	values map[string]any
}

// NewEchoContext creates a new Echo context wrapper
func NewEchoContext(c echo.Context) *EchoContext {
	return &EchoContext{
		Context: c,
		values:  make(map[string]any),
	}
}

// Request returns the HTTP request
func (c *EchoContext) Request() *http.Request {
	return c.Context.Request()
}

// Response returns the response writer
func (c *EchoContext) Response() http.ResponseWriter {
	return c.Context.Response().Writer
}

// Path returns the request path
func (c *EchoContext) Path() string {
	return c.Context.Request().URL.Path
}

// SetPath sets the request path
func (c *EchoContext) SetPath(path string) {
	c.Context.Request().URL.Path = path
}

// Get retrieves a value from the context
func (c *EchoContext) Get(key string) any {
	// Try our local values first
	if val, exists := c.values[key]; exists {
		return val
	}
	// Fall back to Echo's context
	return c.Context.Get(key)
}

// Set stores a value in the context
func (c *EchoContext) Set(key string, val any) {
	c.values[key] = val
	// Also set in Echo's context for compatibility
	c.Context.Set(key, val)
}

// Cookie retrieves a cookie by name
func (c *EchoContext) Cookie(name string) (*http.Cookie, error) {
	return c.Context.Cookie(name)
}

// SetCookie sets a cookie
func (c *EchoContext) SetCookie(cookie *http.Cookie) {
	c.Context.SetCookie(cookie)
}

// Header retrieves a header value
func (c *EchoContext) Header(key string) string {
	return c.Context.Request().Header.Get(key)
}

// SetHeader sets a header value
func (c *EchoContext) SetHeader(key, value string) {
	c.Context.Response().Header().Set(key, value)
}

// Redirect performs an HTTP redirect
func (c *EchoContext) Redirect(code int, url string) error {
	return c.Context.Redirect(code, url)
}

// String sends a string response
func (c *EchoContext) String(code int, s string) error {
	return c.Context.String(code, s)
}

// JSON sends a JSON response
func (c *EchoContext) JSON(code int, i any) error {
	return c.Context.JSON(code, i)
}

// NoContent sends a no content response
func (c *EchoContext) NoContent(code int) error {
	return c.Context.NoContent(code)
}

// AdaptMiddlewareToEcho adapts our middleware to Echo middleware
func AdaptMiddlewareToEcho(mw Middleware) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := NewEchoContext(c)
			return mw.Handle(ctx, func(ctx Context) error {
				return next(c)
			})
		}
	}
}

// AdaptHandlerToEcho adapts our handler to Echo handler
func AdaptHandlerToEcho(handler HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := NewEchoContext(c)
		return handler(ctx)
	}
}

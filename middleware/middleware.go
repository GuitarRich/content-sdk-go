package middleware

import (
	"net/http"
)

// Context is a framework-agnostic HTTP context abstraction
type Context interface {
	// Request returns the HTTP request
	Request() *http.Request

	// Response returns the response writer
	Response() http.ResponseWriter

	// Path returns the request path
	Path() string

	// SetPath sets the request path (for rewrites)
	SetPath(path string)

	// Get retrieves a value from the context
	Get(key string) any

	// Set stores a value in the context
	Set(key string, val any)

	// Cookie retrieves a cookie by name
	Cookie(name string) (*http.Cookie, error)

	// SetCookie sets a cookie
	SetCookie(cookie *http.Cookie)

	// Header retrieves a header value
	Header(key string) string

	// SetHeader sets a header value
	SetHeader(key, value string)

	// Redirect performs an HTTP redirect
	Redirect(code int, url string) error

	// String sends a string response
	String(code int, s string) error

	// JSON sends a JSON response
	JSON(code int, i any) error

	// NoContent sends a no content response
	NoContent(code int) error
}

// HandlerFunc is a framework-agnostic handler function
type HandlerFunc func(ctx Context) error

// Middleware is the base middleware interface
type Middleware interface {
	Handle(ctx Context, next HandlerFunc) error
}

// MiddlewareFunc is a function type that implements Middleware
type MiddlewareFunc func(ctx Context, next HandlerFunc) error

// Handle implements the Middleware interface
func (f MiddlewareFunc) Handle(ctx Context, next HandlerFunc) error {
	return f(ctx, next)
}

// Chain chains multiple middleware together
func Chain(middlewares ...Middleware) Middleware {
	return MiddlewareFunc(func(ctx Context, next HandlerFunc) error {
		// Build the chain from right to left
		handler := next
		for i := len(middlewares) - 1; i >= 0; i-- {
			mw := middlewares[i]
			currentHandler := handler
			handler = func(c Context) error {
				return mw.Handle(c, currentHandler)
			}
		}
		return handler(ctx)
	})
}

// Constants for common context keys
const (
	// SiteKey is the context key for site name
	SiteKey = "site"

	// LocaleKey is the context key for locale/language
	LocaleKey = "locale"

	// OriginalPathKey is the context key for the original path before rewrites
	OriginalPathKey = "originalPath"

	// RewritePathKey is the context key for the rewritten path
	RewritePathKey = "rewritePath"

	// PersonalizeVariantKey is the context key for personalization variant ID
	PersonalizeVariantKey = "personalizeVariant"
)

package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/guitarrich/content-sdk-go/debug"
	"github.com/labstack/echo/v4"
)

// EditingSecurityConfig contains configuration for editing security middleware
type EditingSecurityConfig struct {
	// Secret is the editing secret for validation
	Secret string

	// AllowedOrigins is the list of origins allowed to access editing APIs
	AllowedOrigins []string

	// SkipSecretValidation skips secret validation (for testing)
	SkipSecretValidation bool
}

// EditingSecurityMiddleware validates the editing secret and enforces CORS for editing endpoints
func EditingSecurityMiddleware(config EditingSecurityConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get origin from request
			origin := c.Request().Header.Get("Origin")

			// Handle CORS preflight requests
			if c.Request().Method == http.MethodOptions {
				return handleCORSPreflight(c, origin, config.AllowedOrigins)
			}

			// Validate editing secret from query parameter
			if !config.SkipSecretValidation {
				secret := c.QueryParam("secret")
				if secret == "" {
					debug.Editing("editing secret missing in request")
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"error": "Unauthorized: editing secret is required",
					})
				}

				if secret != config.Secret {
					debug.Editing("invalid editing secret provided")
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"error": "Unauthorized: invalid editing secret",
					})
				}

				debug.Editing("editing secret validated successfully")
			}

			// Set CORS headers for actual requests
			setCORSHeaders(c, origin, config.AllowedOrigins)

			// Set iframe/embedding headers for allowed origins
			setIframeHeaders(c, config.AllowedOrigins)

			return next(c)
		}
	}
}

// handleCORSPreflight handles OPTIONS preflight requests
func handleCORSPreflight(c echo.Context, origin string, allowedOrigins []string) error {
	// Check if origin is allowed
	if origin != "" && isOriginAllowed(origin, allowedOrigins) {
		c.Response().Header().Set("Access-Control-Allow-Origin", origin)
		c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Response().Header().Set("Access-Control-Max-Age", "3600")

		debug.Editing("CORS preflight request handled for origin: %s", origin)
		return c.NoContent(http.StatusNoContent)
	}

	debug.Editing("CORS preflight request rejected for origin: %s", origin)
	return c.NoContent(http.StatusForbidden)
}

// setCORSHeaders sets CORS headers for actual requests
func setCORSHeaders(c echo.Context, origin string, allowedOrigins []string) {
	if origin != "" && isOriginAllowed(origin, allowedOrigins) {
		c.Response().Header().Set("Access-Control-Allow-Origin", origin)
		c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		c.Response().Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		debug.Editing("CORS headers set for origin: %s", origin)
	}
}

// setIframeHeaders sets headers to allow iframe embedding from allowed origins
func setIframeHeaders(c echo.Context, allowedOrigins []string) {
	// If no origins specified, allow all (for development)
	if len(allowedOrigins) == 0 {
		debug.Editing("no allowed origins configured, allowing iframe from all origins (development mode)")
		c.Response().Header().Set("Content-Security-Policy", "frame-ancestors *")
		// Don't set X-Frame-Options when using CSP frame-ancestors
		return
	}

	// Check for wildcard
	if slices.Contains(allowedOrigins, "*") {
		debug.Editing("wildcard configured, allowing iframe from all origins")
		c.Response().Header().Set("Content-Security-Policy", "frame-ancestors *")
		return
	}

	// Build frame-ancestors directive with specific origins
	// The CSP frame-ancestors directive takes a space-separated list of origins
	frameAncestors := strings.Join(allowedOrigins, " ")
	cspHeader := "frame-ancestors " + frameAncestors

	c.Response().Header().Set("Content-Security-Policy", cspHeader)
	debug.Editing("iframe headers set for origins: %v", allowedOrigins)

	// Note: X-Frame-Options is deprecated in favor of CSP frame-ancestors
	// If you need backwards compatibility with very old browsers, you could also set:
	// c.Response().Header().Set("X-Frame-Options", "ALLOW-FROM "+allowedOrigins[0])
	// However, X-Frame-Options only supports a single origin, so CSP is preferred
}

// isOriginAllowed checks if the origin is in the allowed list
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	// If no origins specified, allow all (for development)
	if len(allowedOrigins) == 0 {
		debug.Editing("no allowed origins configured, allowing all origins (development mode)")
		return true
	}

	// Check if origin is in allowed list
	if slices.Contains(allowedOrigins, origin) {
		return true
	}

	// Check for wildcard
	if slices.Contains(allowedOrigins, "*") {
		return true
	}

	return false
}

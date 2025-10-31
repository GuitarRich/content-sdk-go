# Editing API Security

This document describes the security features for Sitecore editing APIs, including secret validation and CORS protection.

## Overview

All editing API endpoints (`/api/editing/*`) are protected with:

1. **Secret-based authentication** - Validates a secret token passed as a query parameter
2. **CORS protection** - Restricts access to specific origins

## Configuration

### Environment Variables

```bash
# Required: Secret key for editing API authentication
EDITING_SECRET=your-secret-key-here

# Required: Comma-separated list of allowed origins
ALLOWED_ORIGINS=https://pages.sitecorecloud.io,https://staging-pages.sitecorecloud.io

# Enable editing mode
EDITING_ENABLED=true
```

### Allowed Origins Format

The `ALLOWED_ORIGINS` environment variable supports:

**Single origin:**

```bash
ALLOWED_ORIGINS=https://pages.sitecorecloud.io
```

**Multiple origins (comma-separated):**

```bash
ALLOWED_ORIGINS=https://pages.sitecorecloud.io,https://staging-pages.sitecorecloud.io,https://preview.sitecorecloud.io
```

**Wildcard (allow all origins):**

```bash
ALLOWED_ORIGINS=*
```

**Empty (development mode - allows all origins):**

```bash
# No ALLOWED_ORIGINS set, or:
ALLOWED_ORIGINS=
```

## Usage

### Making Authenticated Requests

All requests to editing endpoints must include the `secret` query parameter:

```bash
# Config endpoint
GET /api/editing/config?secret=your-secret-key-here

# Render endpoint
POST /api/editing/render?secret=your-secret-key-here
```

### CORS Headers

The middleware automatically sets the following CORS headers for allowed origins:

**For actual requests:**

- `Access-Control-Allow-Origin`: The requesting origin (if allowed)
- `Access-Control-Allow-Credentials`: `true`
- `Access-Control-Expose-Headers`: `Content-Length, Content-Type`

**For preflight requests (OPTIONS):**

- `Access-Control-Allow-Origin`: The requesting origin (if allowed)
- `Access-Control-Allow-Credentials`: `true`
- `Access-Control-Allow-Methods`: `GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers`: `Content-Type, Authorization, X-Requested-With`
- `Access-Control-Max-Age`: `3600`

### Iframe/Embedding Headers

The middleware automatically sets `Content-Security-Policy` headers to allow iframe embedding from allowed origins:

**Single origin:**

```
Content-Security-Policy: frame-ancestors https://pages.sitecorecloud.io
```

**Multiple origins:**

```
Content-Security-Policy: frame-ancestors https://pages.sitecorecloud.io https://pages-eu.sitecorecloud.io
```

**Wildcard (allow all):**

```
Content-Security-Policy: frame-ancestors *
```

**Development mode (empty origins):**

```
Content-Security-Policy: frame-ancestors *
```

This allows the editing endpoints to be embedded in iframes from Sitecore Pages, enabling the in-context editing experience.

## API Responses

### Successful Authentication

```json
HTTP/1.1 200 OK
Access-Control-Allow-Origin: https://pages.sitecorecloud.io
Access-Control-Allow-Credentials: true
Content-Security-Policy: frame-ancestors https://pages.sitecorecloud.io
Content-Type: application/json

{
  "components": [...],
  "packages": {...},
  "editMode": "metadata"
}
```

### Missing Secret

```json
HTTP/1.1 401 Unauthorized
Content-Type: application/json

{
  "error": "Unauthorized: editing secret is required"
}
```

### Invalid Secret

```json
HTTP/1.1 401 Unauthorized
Content-Type: application/json

{
  "error": "Unauthorized: invalid editing secret"
}
```

### CORS Preflight Rejection

```
HTTP/1.1 403 Forbidden
(No response body)
```

## Code Examples

### Programmatic Configuration

Using the config builder:

```go
import (
    "github.com/content-sdk-go/config"
)

cfg := config.NewConfigBuilder().
    WithEdgeAPI("context-id", "client-context-id", "").
    WithEditing(
        true,                                    // enabled
        "my-secret-key",                        // secret
        "http://localhost:8080",                // internal host URL
        "https://pages.sitecorecloud.io",       // allowed origins (variadic)
        "https://staging-pages.sitecorecloud.io",
    ).
    BuildOrPanic()
```

### Setting Up Middleware

```go
import (
    "github.com/content-sdk-go/middleware"
    "github.com/labstack/echo/v4"
)

// Create editing security middleware
editingSecurity := middleware.EditingSecurityMiddleware(middleware.EditingSecurityConfig{
    Secret:         cfg.Editing.Secret,
    AllowedOrigins: cfg.Editing.AllowedOrigins,
})

// Apply to editing route group
editingGroup := e.Group("/api/editing")
editingGroup.Use(editingSecurity)

// Add routes
editingGroup.GET("/config", editingConfigHandler)
editingGroup.POST("/render", editingRenderHandler)
```

### Testing Configuration

For testing purposes, you can skip secret validation:

```go
middleware := middleware.EditingSecurityMiddleware(middleware.EditingSecurityConfig{
    Secret:               "test-secret",
    AllowedOrigins:       []string{"https://test.com"},
    SkipSecretValidation: true,  // Skip validation for tests
})
```

## Security Best Practices

### Production

1. **Use strong secrets**: Generate a long, random secret key
2. **Restrict origins**: Only allow specific Sitecore Cloud origins
3. **Enable HTTPS**: Ensure all origins use HTTPS protocol
4. **Rotate secrets**: Periodically update the editing secret
5. **Monitor access**: Log authentication failures

Example production configuration:

```bash
EDITING_SECRET=a9f4c8e2d1b3a5c7e9f8d2b4a6c8e0f2
ALLOWED_ORIGINS=https://pages.sitecorecloud.io,https://pages-eu.sitecorecloud.io
```

### Development

For local development, you can use relaxed settings:

```bash
EDITING_SECRET=local-dev-secret
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
```

Or allow all origins:

```bash
EDITING_SECRET=local-dev-secret
ALLOWED_ORIGINS=*
```

### Staging

Use environment-specific secrets and restrict to staging origins:

```bash
EDITING_SECRET=staging-secret-key
ALLOWED_ORIGINS=https://staging-pages.sitecorecloud.io
```

## Troubleshooting

### CORS Errors in Browser Console

**Symptom:**

```
Access to fetch at 'http://localhost:8080/api/editing/config' from origin 'https://pages.sitecorecloud.io'
has been blocked by CORS policy
```

**Solutions:**

1. Verify the origin is in `ALLOWED_ORIGINS`
2. Check that the origin matches exactly (including protocol and port)
3. Ensure the server is returning proper CORS headers

### 401 Unauthorized Errors

**Symptom:**

```json
{ "error": "Unauthorized: invalid editing secret" }
```

**Solutions:**

1. Verify the `secret` query parameter is included in the request
2. Check that the secret matches `EDITING_SECRET` exactly
3. Ensure the secret hasn't been changed without updating the client

### Preflight Request Failures

**Symptom:**
OPTIONS request returns 403 Forbidden

**Solutions:**

1. Verify the origin is in `ALLOWED_ORIGINS`
2. Check that the middleware is applied to the editing routes
3. Ensure the preflight request includes the `Origin` header

## Debug Logging

Enable debug mode to see detailed security logs:

```bash
DEBUG=true
```

This will output:

- Secret validation attempts
- Origin checking results
- CORS header operations
- Authentication successes/failures

Example debug output:

```
[Editing] editing secret validated successfully
[Editing] CORS headers set for origin: https://pages.sitecorecloud.io
[Editing] CORS preflight request handled for origin: https://pages.sitecorecloud.io
```

## Related Documentation

- [Editing Configuration](./EDITING_CONFIG.md)
- [Middleware Documentation](./MIDDLEWARE.md)
- [Security Best Practices](./SECURITY.md)

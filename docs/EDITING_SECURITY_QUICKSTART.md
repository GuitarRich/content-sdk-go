# Editing Security - Quick Start Guide

## 5-Minute Setup

### 1. Set Environment Variables

```bash
# Required
EDITING_SECRET=your-secret-key-here
ALLOWED_ORIGINS=https://pages.sitecorecloud.io

# Enable editing
EDITING_ENABLED=true
```

### 2. Apply Middleware (Already Done in Kit Examples)

The kit examples already have this configured in `main.go`:

```go
editingSecurity := middleware.EditingSecurityMiddleware(middleware.EditingSecurityConfig{
    Secret:         cfg.Editing.Secret,
    AllowedOrigins: cfg.Editing.AllowedOrigins,
})

editingGroup := e.Group("/api/editing")
editingGroup.Use(editingSecurity)
```

### 3. Make Requests

All editing API calls must include the `secret` query parameter:

```bash
# Get editing config
curl "http://localhost:8080/api/editing/config?secret=your-secret-key-here"

# Render component
curl -X POST "http://localhost:8080/api/editing/render?secret=your-secret-key-here" \
  -H "Content-Type: application/json" \
  -H "Origin: https://pages.sitecorecloud.io" \
  -d '{"componentName": "Hero", "fields": {...}}'
```

## Common Configurations

### Production

```bash
EDITING_SECRET=prod-strong-secret-key-here
ALLOWED_ORIGINS=https://pages.sitecorecloud.io,https://pages-eu.sitecorecloud.io
```

### Staging

```bash
EDITING_SECRET=staging-secret-key-here
ALLOWED_ORIGINS=https://staging-pages.sitecorecloud.io
```

### Development (Relaxed)

```bash
EDITING_SECRET=dev-secret
ALLOWED_ORIGINS=*
# or
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
```

## API Endpoints Protected

✅ `/api/editing/config` - Component registry and package info  
✅ `/api/editing/render` - Server-side component rendering  
✅ `/api/editing/feaas/render` - FEaaS component rendering

## Iframe Support

The editing endpoints can be embedded in iframes from allowed origins. The same `ALLOWED_ORIGINS` configuration controls both CORS and iframe embedding:

```bash
ALLOWED_ORIGINS=https://pages.sitecorecloud.io
```

This sets the `Content-Security-Policy` header:

```
Content-Security-Policy: frame-ancestors https://pages.sitecorecloud.io
```

This allows Sitecore Pages to embed your application for in-context editing.

## Error Responses

### Missing Secret

```json
{
  "error": "Unauthorized: editing secret is required"
}
```

**Fix:** Add `?secret=your-secret` to the URL

### Invalid Secret

```json
{
  "error": "Unauthorized: invalid editing secret"
}
```

**Fix:** Check that `EDITING_SECRET` matches the secret in your request

### CORS Error

Browser console shows:

```
Access to fetch ... has been blocked by CORS policy
```

**Fix:** Add your origin to `ALLOWED_ORIGINS`

## Testing

To test your setup:

```bash
# Test with valid secret
curl -i "http://localhost:8080/api/editing/config?secret=your-secret-key"

# Should return 200 OK with JSON response

# Test with invalid secret
curl -i "http://localhost:8080/api/editing/config?secret=wrong"

# Should return 401 Unauthorized

# Test CORS preflight
curl -i -X OPTIONS "http://localhost:8080/api/editing/config" \
  -H "Origin: https://pages.sitecorecloud.io"

# Should return 204 No Content with CORS headers
```

## Debug Mode

Enable debug logging to troubleshoot:

```bash
DEBUG=true
```

You'll see:

```
[Editing] editing secret validated successfully
[Editing] CORS headers set for origin: https://pages.sitecorecloud.io
```

## Need Help?

- Full documentation: [EDITING_SECURITY.md](./EDITING_SECURITY.md)
- Test examples: [editing_security_test.go](../middleware/editing_security_test.go)
- Implementation: [editing_security.go](../middleware/editing_security.go)

# Basic Echo Example

A complete example of using the Sitecore Content SDK for Go with the Echo web framework.

## Features

- ✅ Sitecore page rendering
- ✅ Multisite support
- ✅ Locale detection
- ✅ Redirects middleware
- ✅ SEO (robots.txt, sitemap.xml)
- ✅ Editing support
- ✅ Health check endpoint

## Setup

### 1. Install Dependencies

```bash
go get github.com/labstack/echo/v4
go get github.com/content-sdk-go
```

### 2. Configure Environment

Create a `.env` file:

```bash
USE_EDGE_API=true
SITECORE_EDGE_CONTEXT_ID=your-context-id
SITECORE_EDGE_CLIENT_CONTEXT_ID=your-client-context-id
DEFAULT_SITE_NAME=mysite
DEFAULT_LANGUAGE=en
MULTISITE_ENABLED=true
EDITING_ENABLED=true
EDITING_SECRET=your-secret
PORT=3000
```

### 3. Run

```bash
go run main.go
```

## Endpoints

- `GET /` - Catch-all for Sitecore pages
- `GET /healthz` - Health check
- `GET /robots.txt` - Robots.txt
- `GET /sitemap.xml` - Sitemap
- `GET /api/editing/config` - Editing configuration
- `POST /api/editing/render` - Editing render
- `GET /api/preview` - Preview endpoint

## Architecture

```
Request Flow:
  → Logger Middleware
  → Recover Middleware
  → Healthcheck Middleware
  → Multisite Middleware (resolves site)
  → Locale Middleware (detects language)
  → Redirects Middleware (applies redirects)
  → Route Handler (renders page)
```

## Customization

### Add Custom Middleware

```go
e.Use(middleware.AdaptMiddlewareToEcho(yourCustomMiddleware))
```

### Add Custom Routes

```go
e.GET("/custom", func(c echo.Context) error {
    return c.String(http.StatusOK, "Custom route")
})
```

### Modify Configuration

Edit the configuration loading in `main.go`:

```go
cfg, _ := config.NewConfigBuilder().
    WithEdgeAPI(contextID, "", "").
    WithDefaultSite("mysite").
    WithMultisite(true, sites, true).
    Build()
```

## Production Considerations

1. **Environment Variables**: Use a secrets manager (AWS Secrets Manager, HashiCorp Vault)
2. **Logging**: Add structured logging (e.g., zap, logrus)
3. **Metrics**: Add Prometheus metrics
4. **Caching**: Add Redis for layout caching
5. **Rate Limiting**: Add rate limiting middleware
6. **CORS**: Configure CORS if needed
7. **TLS**: Use TLS/HTTPS in production

## Troubleshooting

### Server won't start
- Check if port 3000 is available
- Verify environment variables
- Check logs for configuration errors

### Pages not rendering
- Verify `SITECORE_EDGE_CONTEXT_ID`
- Check site name matches Sitecore
- Enable debug logging: `DEBUG=1`

### Multisite not working
- Verify hostname configuration
- Check DNS/hosts file
- Test with cookie-based resolution

## Next Steps

1. Add rendering logic (HTML templates)
2. Implement component rendering
3. Add caching layer
4. Set up CI/CD pipeline
5. Deploy to cloud (AWS, GCP, Azure)


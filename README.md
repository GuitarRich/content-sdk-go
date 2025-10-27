# Sitecore Content SDK for Go

A comprehensive Go SDK for building Sitecore XM Cloud applications with multisite, personalization, and editing support.

## 🚀 Features

- **Framework Agnostic** - Works with Echo, Gin, net/http, and other Go web frameworks
- **Multisite Support** - Automatic site resolution by hostname, path, or cookie
- **Personalization** - Sitecore Personalize (CDP) integration
- **Editing Support** - Full Sitecore Pages editor integration
- **SEO Services** - Sitemap and robots.txt generation
- **i18n** - Dictionary service for internationalization
- **Media API** - Image URL generation with transformations
- **Type Safe** - Comprehensive Go structs for all data models
- **Testable** - Interface-based design for easy testing
- **Performant** - Context support, timeouts, and retries built-in

## 📦 Installation

```bash
go get github.com/content-sdk-go
```

## 🏗️ Quick Start

### 1. Configuration

Create a `.env` file (see `env.example`):

```bash
USE_EDGE_API=true
SITECORE_EDGE_CONTEXT_ID=your-context-id
DEFAULT_SITE_NAME=mysite
DEFAULT_LANGUAGE=en
```

### 2. Basic Usage

```go
package main

import (
    "github.com/content-sdk-go/config"
    "github.com/content-sdk-go/client"
    "github.com/content-sdk-go/handlers"
    "github.com/content-sdk-go/middleware"
    "github.com/labstack/echo/v4"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()
    if err := cfg.Validate(); err != nil {
        panic(err)
    }

    // Create Sitecore client
    sitecoreClient := client.NewSitecoreClient(client.ClientConfig{
        APIEndpoint: cfg.GetGraphQLEndpoint(),
        APIKey:      cfg.GetAPIKey(),
        SiteName:    cfg.DefaultSite,
    })

    // Create Echo app
    e := echo.New()

    // Add middleware chain
    e.Use(middleware.AdaptMiddlewareToEcho(
        middleware.NewMultisiteMiddleware(middleware.MultisiteConfig{
            DefaultSite: cfg.DefaultSite,
        }),
    ))

    // Add catch-all handler
    catchAll := handlers.NewCatchAllHandler(sitecoreClient)
    e.GET("/*", func(c echo.Context) error {
        return catchAll.Handle(middleware.NewEchoContext(c))
    })

    // Start server
    e.Start(":3000")
}
```

### 3. With Multisite

```go
// Configure multisite
sites := []models.SiteInfo{
    {Name: "site1", Language: "en", HostName: "www.site1.com"},
    {Name: "site2", Language: "fr", HostName: "www.site2.com"},
}

cfg, _ := config.NewConfigBuilder().
    WithEdgeAPI(contextID, "", "").
    WithDefaultSite("site1").
    WithMultisite(true, sites, true).
    Build()
```

### 4. With Personalization

```go
cfg, _ := config.NewConfigBuilder().
    WithEdgeAPI(contextID, "", "").
    WithDefaultSite("mysite").
    WithPersonalization(true, "your-scope", "").
    Build()

// Add personalization middleware
e.Use(middleware.AdaptMiddlewareToEcho(
    middleware.NewPersonalizeMiddleware(middleware.PersonalizeConfig{
        Enabled: true,
        Scope:   cfg.Personalize.Scope,
    }),
))
```

## 📚 Package Structure

```
content-sdk-go/
├── client/           # Core Sitecore client (GetPage, GetPreview, etc.)
├── config/           # Configuration management
├── graphql/          # GraphQL client with retries
├── handlers/         # HTTP handlers (catch-all, robots, sitemap, editing)
├── i18n/             # Dictionary service
├── layoutService/    # Layout service for fetching page data
├── media/            # Media API for image URLs
├── middleware/       # Framework-agnostic middleware
├── models/           # Data models
├── seo/              # SEO services (sitemap, robots, error pages)
├── site/             # Site resolution and redirects
└── utils/            # Utilities (env, http)
```

## 🔧 Core Concepts

### Client

The `SitecoreClient` is the main entry point for fetching page data:

```go
client := client.NewSitecoreClient(client.ClientConfig{
    APIEndpoint: "https://edge.sitecorecloud.io/api/graphql/v1",
    APIKey:      "your-context-id",
    SiteName:    "mysite",
})

page, err := client.GetPage("/products", models.PageOptions{
    Site:   "mysite",
    Locale: stringPtr("en"),
})
```

### Middleware

Middleware processes requests before they reach handlers:

```go
chain := middleware.Chain(
    middleware.NewHealthcheckMiddleware(...),
    middleware.NewMultisiteMiddleware(...),
    middleware.NewLocaleMiddleware(...),
    middleware.NewRedirectsMiddleware(...),
)
```

### Handlers

Handlers respond to specific routes:

```go
// Catch-all for dynamic pages
catchAll := handlers.NewCatchAllHandler(sitecoreClient)

// Robots.txt
robotsHandler := handlers.NewRobotsHandler(...)

// Sitemap.xml
sitemapHandler := handlers.NewSitemapHandler(...)

// Editing config
editingConfig := handlers.NewEditingConfigHandler(...)

// Editing render
editingRender := handlers.NewEditingRenderHandler(sitecoreClient)
```

## 🔌 Framework Integration

### Echo

```go
import (
    "github.com/content-sdk-go/middleware"
    "github.com/labstack/echo/v4"
)

e := echo.New()

// Adapt middleware
e.Use(middleware.AdaptMiddlewareToEcho(yourMiddleware))

// Use handlers
e.GET("/*", func(c echo.Context) error {
    ctx := middleware.NewEchoContext(c)
    return handler.Handle(ctx)
})
```

### Gin (create adapter)

```go
func AdaptMiddlewareToGin(mw middleware.Middleware) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := NewGinContext(c)
        mw.Handle(ctx, func(ctx middleware.Context) error {
            c.Next()
            return nil
        })
    }
}
```

### net/http

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    ctx := middleware.NewNetHTTPContext(w, r)
    handler.Handle(ctx)
})
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./client/...
```

## 📖 API Documentation

### Client Methods

- `GetPage(path string, options PageOptions) (*Page, error)` - Fetch a page
- `GetPreview(data PreviewData) (*Page, error)` - Fetch preview data
- `GetDesignLibraryData(data DesignLibraryRenderPreviewData) (*Page, error)` - Design library
- `GetStaticPaths(site string) ([]string, error)` - Get all static paths
- `GetSiteNameFromPath(path string) string` - Extract site from path
- `ParsePath(path string) string` - Parse and normalize path

### Services

- **DictionaryService** - Fetch i18n phrases
- **SiteInfoService** - Fetch site configuration
- **SiteResolver** - Resolve sites by hostname or name
- **RedirectsService** - Fetch and match redirects
- **MediaAPI** - Generate image URLs with transformations
- **SitemapXmlService** - Generate sitemaps
- **RobotsService** - Generate robots.txt
- **ErrorPagesService** - Fetch custom error pages

## 🌐 Environment Variables

See `env.example` for a complete list of environment variables.

Key variables:

- `USE_EDGE_API` - Use Edge API (true) or Local API (false)
- `SITECORE_EDGE_CONTEXT_ID` - Edge context ID
- `SITECORE_API_KEY` - Local API key
- `SITECORE_API_HOST` - Local API host
- `DEFAULT_SITE_NAME` - Default site name
- `DEFAULT_LANGUAGE` - Default language
- `MULTISITE_ENABLED` - Enable multisite
- `PERSONALIZE_ENABLED` - Enable personalization

## 🤝 Contributing

Contributions are welcome! Please ensure:

- All tests pass: `go test ./...`
- Code is formatted: `go fmt ./...`
- Follow Go best practices
- Add tests for new functionality

## 📄 License

This project is licensed under the MIT License.

## 🔗 Resources

- [Sitecore XM Cloud Documentation](https://doc.sitecore.com/xmc)
- [Sitecore Content SDK (TypeScript/JS)](https://github.com/Sitecore/content-sdk)
- [Migration Guide](./MIGRATION_PLAN.md)

## 💡 Examples

See the `examples/` directory for complete working examples:

- Basic Echo application
- Multisite configuration
- Personalization integration
- Editing support

## 🐛 Troubleshooting

### "Page not found" errors

- Verify `DEFAULT_SITE_NAME` matches your Sitecore site
- Check `SITECORE_EDGE_CONTEXT_ID` is correct
- Ensure the path exists in Sitecore

### GraphQL errors

- Verify API key/context ID
- Check network connectivity
- Enable debug logging: `DEBUG=1`

### Multisite issues

- Verify hostname configuration
- Check site resolver configuration
- Test cookie-based resolution

## 🚀 Next Steps

1. Configure your Sitecore connection
2. Set up middleware chain
3. Add handlers for your routes
4. Implement rendering logic
5. Deploy to production

Happy coding! 🎉

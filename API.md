# API Reference

Complete API documentation for the Sitecore Content SDK for Go.

## Table of Contents

- [Client](#client)
- [Configuration](#configuration)
- [Services](#services)
- [Middleware](#middleware)
- [Handlers](#handlers)
- [Models](#models)

---

## Client

### SitecoreClient

The main client for interacting with Sitecore.

#### Constructor

```go
func NewSitecoreClient(config ClientConfig) *SitecoreClient
```

**ClientConfig:**

```go
type ClientConfig struct {
    APIEndpoint string              // GraphQL endpoint
    APIKey      string              // API key or context ID
    SiteName    string              // Default site name
    Timeout     time.Duration       // Request timeout (optional)
}
```

#### Methods

##### GetPage

Fetches a page from Sitecore.

```go
func (c *SitecoreClient) GetPage(path string, options PageOptions) (*Page, error)
```

**Parameters:**

- `path` (string): The page path (e.g., "/products")
- `options` (PageOptions): Fetch options

**Returns:**

- `*Page`: The page data
- `error`: Error if any

**Example:**

```go
locale := "en"
page, err := client.GetPage("/products", models.PageOptions{
    Site:   "mysite",
    Locale: &locale,
})
```

##### GetPreview

Fetches preview data for editing.

```go
func (c *SitecoreClient) GetPreview(data PreviewData) (*Page, error)
```

##### GetDesignLibraryData

Fetches design library data.

```go
func (c *SitecoreClient) GetDesignLibraryData(data DesignLibraryRenderPreviewData) (*Page, error)
```

##### GetStaticPaths

Gets all static paths for a site.

```go
func (c *SitecoreClient) GetStaticPaths(site string) ([]string, error)
```

##### GetSiteNameFromPath

Extracts the site name from a path.

```go
func (c *SitecoreClient) GetSiteNameFromPath(path string) string
```

##### ParsePath

Parses and normalizes a path.

```go
func (c *SitecoreClient) ParsePath(path string) string
```

---

## Configuration

### Config

Central configuration for the SDK.

#### Loading Configuration

```go
// From environment variables
cfg := config.LoadConfig()

// Using builder pattern
cfg, err := config.NewConfigBuilder().
    WithEdgeAPI(contextID, clientContextID, "").
    WithDefaultSite("mysite").
    Build()
```

#### Methods

##### Validate

Validates the configuration.

```go
func (c *Config) Validate() error
```

##### GetGraphQLEndpoint

Returns the GraphQL endpoint based on configuration.

```go
func (c *Config) GetGraphQLEndpoint() string
```

##### GetAPIKey

Returns the API key based on configuration.

```go
func (c *Config) GetAPIKey() string
```

#### Builder Methods

```go
func (b *ConfigBuilder) WithEdgeAPI(contextID, clientContextID, edgeURL string) *ConfigBuilder
func (b *ConfigBuilder) WithLocalAPI(apiKey, apiHost string) *ConfigBuilder
func (b *ConfigBuilder) WithDefaultSite(siteName string) *ConfigBuilder
func (b *ConfigBuilder) WithDefaultLanguage(language string) *ConfigBuilder
func (b *ConfigBuilder) WithMultisite(enabled bool, sites []SiteInfo, useCookie bool) *ConfigBuilder
func (b *ConfigBuilder) WithPersonalization(enabled bool, scope, cdpEndpoint string) *ConfigBuilder
func (b *ConfigBuilder) WithEditing(enabled bool, secret, internalHostURL string) *ConfigBuilder
func (b *ConfigBuilder) WithTimeouts(edgeTimeout, cdpTimeout time.Duration) *ConfigBuilder
func (b *ConfigBuilder) WithDebug(enabled bool) *ConfigBuilder
func (b *ConfigBuilder) Build() (*Config, error)
```

---

## Services

### DictionaryService

Fetches i18n dictionary phrases.

#### Constructor

```go
func NewDictionaryService(graphQLEndpoint, apiKey string) DictionaryService
```

#### Methods

```go
func (s *DictionaryService) FetchDictionaryData(ctx context.Context, locale, site string) (DictionaryPhrases, error)
```

**Example:**

```go
service := i18n.NewDictionaryService(endpoint, apiKey)
phrases, err := service.FetchDictionaryData(ctx, "en", "mysite")
```

---

### SiteInfoService

Fetches site configuration.

#### Constructor

```go
func NewSiteInfoService(graphQLEndpoint, apiKey string) SiteInfoService
```

#### Methods

```go
func (s *SiteInfoService) FetchSiteInfo(ctx context.Context, siteName string) (*SiteInfo, error)
func (s *SiteInfoService) FetchSites(ctx context.Context) ([]SiteInfo, error)
```

---

### SiteResolver

Resolves sites by hostname or name.

#### Constructor

```go
func NewSiteResolver(siteService SiteInfoService, cacheTimeout time.Duration) *SiteResolver
```

#### Methods

```go
func (r *SiteResolver) GetByHost(hostname string) (*SiteInfo, error)
func (r *SiteResolver) GetByName(name string) (*SiteInfo, error)
```

---

### RedirectsService

Manages URL redirects.

#### Constructor

```go
func NewRedirectsService(graphQLEndpoint, apiKey string) RedirectsService
```

#### Methods

```go
func (s *RedirectsService) FetchRedirects(ctx context.Context, site string) ([]RedirectInfo, error)
func (s *RedirectsService) GetRedirect(path string, redirects []RedirectInfo) (*RedirectInfo, error)
```

---

### MediaAPI

Generates image URLs with transformations.

#### Constructor

```go
func NewMediaAPI(mediaHost string) *MediaAPI
```

#### Methods

```go
func (m *MediaAPI) GetImageURL(imageField interface{}, params *ImageParams) string
func (m *MediaAPI) GetResponsiveImageURL(imageField interface{}, widths []int) map[int]string
```

**ImageParams:**

```go
type ImageParams struct {
    Width         int
    Height        int
    Quality       int
    Format        string
    FocalPoint    *FocalPoint
    CropMode      string
}
```

**Example:**

```go
api := media.NewMediaAPI("https://cdn.example.com")
url := api.GetImageURL(imageField, &media.ImageParams{
    Width:   800,
    Height:  600,
    Quality: 90,
})
```

---

### SitemapXmlService

Generates XML sitemaps.

#### Constructor

```go
func NewSitemapXmlService(graphQLEndpoint, apiKey string) SitemapXmlService
```

#### Methods

```go
func (s *SitemapXmlService) FetchSitemap(ctx context.Context, sites, languages []string) ([]SitemapEntry, error)
func (s *SitemapXmlService) GenerateSitemapXML(entries []SitemapEntry) (string, error)
```

---

### RobotsService

Generates robots.txt.

#### Constructor

```go
func NewRobotsService(graphQLEndpoint, apiKey string) RobotsService
```

#### Methods

```go
func (s *RobotsService) FetchRobotsDirectives(ctx context.Context, site string) (*RobotsQueryResult, error)
func (s *RobotsService) GenerateRobotsTxt(directive *RobotsQueryResult, sitemapURLs []string) string
```

---

### ErrorPagesService

Fetches custom error pages.

#### Constructor

```go
func NewErrorPagesService(graphQLEndpoint, apiKey string) ErrorPagesService
```

#### Methods

```go
func (s *ErrorPagesService) FetchErrorPages(ctx context.Context, site string) (*ErrorPages, error)
```

---

## Middleware

### Base Middleware Interface

```go
type Middleware interface {
    Handle(ctx Context, next HandlerFunc) error
}

type HandlerFunc func(ctx Context) error
```

### MultisiteMiddleware

Resolves the current site.

#### Constructor

```go
func NewMultisiteMiddleware(config MultisiteConfig) Middleware
```

**MultisiteConfig:**

```go
type MultisiteConfig struct {
    DefaultSite         string
    UseCookieResolution bool
    SiteResolver        *SiteResolver // optional
}
```

#### Context Keys

- `middleware.SiteKey` - The resolved site name

---

### LocaleMiddleware

Detects and sets the language/locale.

#### Constructor

```go
func NewLocaleMiddleware(config LocaleConfig) Middleware
```

**LocaleConfig:**

```go
type LocaleConfig struct {
    DefaultLanguage    string
    SupportedLanguages []string
    UseURLPrefix       bool
}
```

#### Context Keys

- `middleware.LocaleKey` - The resolved locale

---

### RedirectsMiddleware

Applies URL redirects.

#### Constructor

```go
func NewRedirectsMiddleware(config RedirectsConfig) Middleware
```

**RedirectsConfig:**

```go
type RedirectsConfig struct {
    RedirectsService RedirectsService
    CacheDuration    time.Duration
}
```

---

### PersonalizeMiddleware

Handles personalization.

#### Constructor

```go
func NewPersonalizeMiddleware(config PersonalizeConfig) Middleware
```

**PersonalizeConfig:**

```go
type PersonalizeConfig struct {
    Enabled     bool
    Scope       string
    CDPEndpoint string
}
```

---

### HealthcheckMiddleware

Health check endpoint.

#### Constructor

```go
func NewHealthcheckMiddleware(config HealthcheckConfig) Middleware
```

**HealthcheckConfig:**

```go
type HealthcheckConfig struct {
    Path     string
    Response map[string]string
}
```

---

## Handlers

### CatchAllHandler

Handles dynamic Sitecore pages.

#### Constructor

```go
func NewCatchAllHandler(sitecoreClient *SitecoreClient) *CatchAllHandler
```

#### Method

```go
func (h *CatchAllHandler) Handle(ctx Context) error
```

---

### RobotsHandler

Handles robots.txt requests.

#### Constructor

```go
func NewRobotsHandler(config RobotsHandlerConfig) *RobotsHandler
```

---

### SitemapHandler

Handles sitemap.xml requests.

#### Constructor

```go
func NewSitemapHandler(config SitemapHandlerConfig) *SitemapHandler
```

---

### EditingConfigHandler

Provides editing configuration.

#### Constructor

```go
func NewEditingConfigHandler(config EditingConfigHandlerConfig) *EditingConfigHandler
```

---

### EditingRenderHandler

Handles editing render requests.

#### Constructor

```go
func NewEditingRenderHandler(sitecoreClient *SitecoreClient) *EditingRenderHandler
```

---

## Models

### Page

```go
type Page struct {
    LayoutData interface{}       `json:"layoutData"`
    Dictionary DictionaryPhrases `json:"dictionary,omitempty"`
    ErrorPages *ErrorPages        `json:"errorPages,omitempty"`
    HeadLinks  []HTMLLink         `json:"headLinks,omitempty"`
}
```

### PageOptions

```go
type PageOptions struct {
    Site           string
    Locale         *string
    IncludeDict    bool
    IncludeErrors  bool
}
```

### SiteInfo

```go
type SiteInfo struct {
    Name     string `json:"name"`
    HostName string `json:"hostName"`
    Language string `json:"language"`
    RootPath string `json:"rootPath,omitempty"`
    Database string `json:"database,omitempty"`
}
```

### RedirectInfo

```go
type RedirectInfo struct {
    Pattern     string `json:"pattern"`
    Target      string `json:"target"`
    RedirectType string `json:"redirectType"`
    IsRegex     bool   `json:"isRegex"`
    Locale      string `json:"locale,omitempty"`
}
```

### PreviewData

```go
type PreviewData struct {
    ItemID     string `json:"itemId"`
    Language   string `json:"language"`
    Site       string `json:"site"`
    Version    string `json:"version,omitempty"`
    Mode       string `json:"mode,omitempty"`
}
```

---

## Error Types

```go
type NotFoundError struct {
    Path string
}

type PreviewError struct {
    Message string
}

type GraphQLError struct {
    Message string
    Errors  []GraphQLErrorDetail
}

type ValidationError struct {
    Field   string
    Message string
}
```

---

## Constants

### Context Keys

```go
const (
    SiteKey   = "sitecore:site"
    LocaleKey = "sitecore:locale"
)
```

### Preview Modes

```go
const (
    PreviewModeEdit        = "edit"
    PreviewModePreview     = "preview"
    PreviewModeDesignLibrary = "design-library"
)
```

---

## Complete Example

```go
package main

import (
    "github.com/content-sdk-go/client"
    "github.com/content-sdk-go/config"
    "github.com/content-sdk-go/handlers"
    "github.com/content-sdk-go/middleware"
    "github.com/labstack/echo/v4"
)

func main() {
    // Configuration
    cfg := config.LoadConfig()
    cfg.Validate()

    // Client
    sitecoreClient := client.NewSitecoreClient(client.ClientConfig{
        APIEndpoint: cfg.GetGraphQLEndpoint(),
        APIKey:      cfg.GetAPIKey(),
        SiteName:    cfg.DefaultSite,
    })

    // Echo app
    e := echo.New()

    // Middleware
    e.Use(middleware.AdaptMiddlewareToEcho(
        middleware.NewMultisiteMiddleware(middleware.MultisiteConfig{
            DefaultSite: cfg.DefaultSite,
        }),
    ))

    // Handlers
    catchAll := handlers.NewCatchAllHandler(sitecoreClient)
    e.GET("/*", func(c echo.Context) error {
        return catchAll.Handle(middleware.NewEchoContext(c))
    })

    e.Start(":3000")
}
```

---

## See Also

- [README](./README.md) - Getting started guide
- [Migration Guide](./MIGRATION_GUIDE.md) - Migrating from TypeScript
- [Architecture](./ARCHITECTURE.md) - Architecture overview
- [Examples](./examples/) - Complete examples

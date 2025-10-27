# Migration Guide: TypeScript/Next.js to Go

This guide helps you migrate from `@sitecore-content-sdk/nextjs` (TypeScript) to `content-sdk-go`.

## Overview

The Go SDK maintains similar concepts and patterns to the TypeScript SDK, adapted for Go's idioms and type system.

## Key Differences

### 1. Type System

**TypeScript:**
```typescript
interface Page {
  layoutData: LayoutServiceData;
  dictionary?: DictionaryPhrases;
}
```

**Go:**
```go
type Page struct {
    LayoutData interface{} `json:"layoutData"`
    Dictionary DictionaryPhrases `json:"dictionary,omitempty"`
}
```

### 2. Error Handling

**TypeScript:**
```typescript
try {
  const page = await client.getPage('/products');
} catch (error) {
  console.error(error);
}
```

**Go:**
```go
page, err := client.GetPage("/products", models.PageOptions{})
if err != nil {
    log.Printf("Error: %v", err)
}
```

### 3. Configuration

**TypeScript:**
```typescript
export default defineConfig({
  sitecoreEdgeContextId: process.env.SITECORE_EDGE_CONTEXT_ID,
  defaultSiteName: 'mysite',
});
```

**Go:**
```go
cfg := config.NewConfigBuilder().
    WithEdgeAPI(contextID, "", "").
    WithDefaultSite("mysite").
    Build()
```

### 4. Middleware

**TypeScript (Next.js):**
```typescript
export function middleware(request: NextRequest) {
  const multisite = new MultisiteMiddleware();
  return multisite.handler(request);
}
```

**Go:**
```go
mw := middleware.NewMultisiteMiddleware(config)
e.Use(middleware.AdaptMiddlewareToEcho(mw))
```

## Feature Mapping

### Client Methods

| TypeScript | Go | Notes |
|------------|-----|-------|
| `client.getPage()` | `client.GetPage()` | Pascal case |
| `client.getPreview()` | `client.GetPreview()` | Pascal case |
| `client.getStaticPaths()` | `client.GetStaticPaths()` | Pascal case |
| `client.getSiteNameFromPath()` | `client.GetSiteNameFromPath()` | Pascal case |

### Services

| TypeScript Package | Go Package | Location |
|-------------------|-----------|----------|
| `@sitecore-content-sdk/core/dictionary` | `i18n` | `i18n/dictionary.go` |
| `@sitecore-content-sdk/core/site` | `site` | `site/siteinfo.go` |
| `@sitecore-content-sdk/core/redirects` | `site` | `site/redirects.go` |
| `@sitecore-content-sdk/core/sitemap` | `seo` | `seo/sitemap.go` |
| `@sitecore-content-sdk/core/robots` | `seo` | `seo/robots.go` |
| `@sitecore-content-sdk/core/media` | `media` | `media/media.go` |

### Middleware

| TypeScript | Go | Config |
|-----------|-----|--------|
| `MultisiteMiddleware` | `MultisiteMiddleware` | `MultisiteConfig` |
| `LocaleMiddleware` | `LocaleMiddleware` | `LocaleConfig` |
| `RedirectsMiddleware` | `RedirectsMiddleware` | `RedirectsConfig` |
| `PersonalizeMiddleware` | `PersonalizeMiddleware` | `PersonalizeConfig` |

### Handlers

| TypeScript Route | Go Handler | Function |
|-----------------|------------|----------|
| `app/[[...path]]/page.tsx` | `CatchAllHandler` | Dynamic pages |
| `app/robots.txt/route.ts` | `RobotsHandler` | Robots.txt |
| `app/sitemap.xml/route.ts` | `SitemapHandler` | Sitemap |
| `app/api/editing/config/route.ts` | `EditingConfigHandler` | Editing config |
| `app/api/editing/render/route.ts` | `EditingRenderHandler` | Editing render |

## Step-by-Step Migration

### 1. Setup

**Install Go SDK:**
```bash
go get github.com/content-sdk-go
```

**Create configuration:**
```go
// config/config.go
cfg := config.LoadConfig()
```

### 2. Create Client

**TypeScript:**
```typescript
import { createSitecoreNextJSClient } from '@sitecore-content-sdk/nextjs';

const client = createSitecoreNextJSClient({
  sitecoreEdgeContextId: process.env.SITECORE_EDGE_CONTEXT_ID,
});
```

**Go:**
```go
import "github.com/content-sdk-go/client"

client := client.NewSitecoreClient(client.ClientConfig{
    APIEndpoint: cfg.GetGraphQLEndpoint(),
    APIKey:      cfg.GetAPIKey(),
    SiteName:    cfg.DefaultSite,
})
```

### 3. Fetch Page Data

**TypeScript:**
```typescript
const page = await client.getPage('/products', {
  site: 'mysite',
  locale: 'en',
});
```

**Go:**
```go
locale := "en"
page, err := client.GetPage("/products", models.PageOptions{
    Site:   "mysite",
    Locale: &locale,
})
if err != nil {
    // Handle error
}
```

### 4. Middleware Setup

**TypeScript (Next.js middleware.ts):**
```typescript
export function middleware(request: NextRequest) {
  return multisiteMiddleware(request);
}
```

**Go (Echo):**
```go
e.Use(middleware.AdaptMiddlewareToEcho(
    middleware.NewMultisiteMiddleware(config),
))
```

### 5. Route Handlers

**TypeScript (app/[[...path]]/page.tsx):**
```typescript
export default async function Page({ params }) {
  const page = await client.getPage(params.path.join('/'));
  return <PageRenderer page={page} />;
}
```

**Go:**
```go
catchAll := handlers.NewCatchAllHandler(sitecoreClient)
e.GET("/*", func(c echo.Context) error {
    ctx := middleware.NewEchoContext(c)
    return catchAll.Handle(ctx)
})
```

## Common Patterns

### 1. Multisite

**TypeScript:**
```typescript
const config = defineConfig({
  multisite: {
    sites: [
      { name: 'site1', hostName: 'www.site1.com' },
      { name: 'site2', hostName: 'www.site2.com' },
    ],
  },
});
```

**Go:**
```go
sites := []models.SiteInfo{
    {Name: "site1", HostName: "www.site1.com"},
    {Name: "site2", HostName: "www.site2.com"},
}

cfg, _ := config.NewConfigBuilder().
    WithMultisite(true, sites, true).
    Build()
```

### 2. Personalization

**TypeScript:**
```typescript
const config = defineConfig({
  personalize: {
    enabled: true,
    scope: 'your-scope',
  },
});
```

**Go:**
```go
cfg, _ := config.NewConfigBuilder().
    WithPersonalization(true, "your-scope", "").
    Build()
```

### 3. Editing/Preview

**TypeScript:**
```typescript
const config = defineConfig({
  editing: {
    enabled: true,
    secret: process.env.EDITING_SECRET,
  },
});
```

**Go:**
```go
cfg, _ := config.NewConfigBuilder().
    WithEditing(true, secret, internalHostURL).
    Build()
```

## Async/Await vs Error Handling

### TypeScript Pattern
```typescript
async function fetchData() {
  try {
    const data = await service.fetch();
    return data;
  } catch (error) {
    console.error(error);
    throw error;
  }
}
```

### Go Pattern
```go
func fetchData() (*Data, error) {
    data, err := service.Fetch()
    if err != nil {
        log.Printf("Error: %v", err)
        return nil, err
    }
    return data, nil
}
```

## Context & Cancellation

Go provides context for request cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

data, err := service.FetchWithContext(ctx)
```

## Testing

### TypeScript
```typescript
import { describe, it, expect } from 'vitest';

describe('Client', () => {
  it('fetches page', async () => {
    const page = await client.getPage('/test');
    expect(page).toBeDefined();
  });
});
```

### Go
```go
import "testing"

func TestClient_GetPage(t *testing.T) {
    page, err := client.GetPage("/test", models.PageOptions{})
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if page == nil {
        t.Error("expected page to be defined")
    }
}
```

## Performance Considerations

1. **Concurrency**: Use goroutines for parallel requests
2. **Context**: Always pass context for timeouts/cancellation
3. **Caching**: Implement caching for layout data
4. **Connection Pooling**: Use HTTP client with connection pooling

## Best Practices

1. **Error Handling**: Always check errors explicitly
2. **Interfaces**: Use interfaces for testability
3. **Context**: Pass context through the call chain
4. **Defer**: Use defer for cleanup (e.g., `defer cancel()`)
5. **Pointers**: Use pointers for optional fields in structs
6. **JSON Tags**: Always add JSON tags to struct fields

## Troubleshooting

### Type Assertions

When using `interface{}` types (e.g., `Page.LayoutData`):

```go
if layoutData, ok := page.LayoutData.(*layoutService.LayoutServiceData); ok {
    // Use layoutData
}
```

### Nil Checks

Always check for nil before dereferencing:

```go
if page != nil && page.Dictionary != nil {
    // Use dictionary
}
```

### Import Cycles

If you encounter import cycles:
- Use interfaces
- Move shared types to a common package
- Use `interface{}` and type assertions

## Next Steps

1. Review the [API Documentation](./API.md)
2. Explore the [Examples](./examples/)
3. Read the [Architecture Guide](./ARCHITECTURE.md)
4. Check the [Migration Plan](./MIGRATION_PLAN.md)

## Support

- GitHub Issues: [Report issues](https://github.com/your-repo/issues)
- Documentation: [Read the docs](./README.md)
- Examples: [See examples](./examples/)


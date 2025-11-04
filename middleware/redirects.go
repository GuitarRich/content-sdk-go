package middleware

import (
	"context"
	"net/http"

	"github.com/guitarrich/content-sdk-go/debug"
	"github.com/guitarrich/content-sdk-go/models"
	"github.com/guitarrich/content-sdk-go/site"
)

// RedirectsConfig contains configuration for redirects middleware
type RedirectsConfig struct {
	// RedirectsService is the service to fetch redirects
	RedirectsService site.RedirectsService

	// Site is the site name to fetch redirects for
	Site string

	// RefreshInterval is how often to refresh redirects (in seconds)
	// Set to 0 to disable auto-refresh
	RefreshInterval int
}

// RedirectsMiddleware handles URL redirects
type RedirectsMiddleware struct {
	config    RedirectsConfig
	redirects []models.RedirectInfo
}

// NewRedirectsMiddleware creates a new redirects middleware
func NewRedirectsMiddleware(config RedirectsConfig) *RedirectsMiddleware {
	return &RedirectsMiddleware{
		config:    config,
		redirects: nil, // Will be loaded on first request
	}
}

// Handle processes the redirects middleware
func (m *RedirectsMiddleware) Handle(ctx Context, next HandlerFunc) error {
	path := ctx.Path()

	debug.Redirects("checking redirects for path=%s", path)

	// Load redirects if not already loaded
	if m.redirects == nil {
		if err := m.loadRedirects(ctx); err != nil {
			debug.Redirects("failed to load redirects: %v", err)
			// Continue without redirects
			return next(ctx)
		}
	}

	// Check for matching redirect
	redirect, err := m.config.RedirectsService.GetRedirect(path, m.redirects)
	if err != nil {
		debug.Redirects("error checking redirect: %v", err)
		return next(ctx)
	}

	// No redirect found, continue normally
	if redirect == nil {
		return next(ctx)
	}

	debug.Redirects("redirect found: %s -> %s (type=%s)", path, redirect.Target, redirect.RedirectType)

	// Apply redirect based on type
	switch redirect.RedirectType {
	case models.Redirect301:
		return ctx.Redirect(http.StatusMovedPermanently, redirect.Target)

	case models.Redirect302:
		return ctx.Redirect(http.StatusFound, redirect.Target)

	case models.RedirectServerTransfer:
		// Server transfer: rewrite the path and continue
		ctx.SetPath(redirect.Target)
		ctx.Set(OriginalPathKey, path)
		ctx.Set(RewritePathKey, redirect.Target)
		return next(ctx)

	default:
		// Unknown redirect type, use 302
		return ctx.Redirect(http.StatusFound, redirect.Target)
	}
}

// loadRedirects loads redirects from the service
func (m *RedirectsMiddleware) loadRedirects(ctx Context) error {
	// Get site from context if available
	site := m.config.Site
	if siteFromCtx := ctx.Get(SiteKey); siteFromCtx != nil {
		if siteStr, ok := siteFromCtx.(string); ok {
			site = siteStr
		}
	}

	// Fetch redirects
	redirects, err := m.config.RedirectsService.FetchRedirects(context.Background(), site)
	if err != nil {
		return err
	}

	m.redirects = redirects
	debug.Redirects("loaded %d redirects for site %s", len(redirects), site)
	return nil
}

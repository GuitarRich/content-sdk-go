package middleware

import (
	"net/http"
	"strings"

	"github.com/guitarrich/content-sdk-go/client"
	"github.com/guitarrich/content-sdk-go/debug"
	"github.com/guitarrich/content-sdk-go/models"
	"github.com/guitarrich/content-sdk-go/site"
)

// MultisiteConfig contains configuration for multisite middleware
type MultisiteConfig struct {
	// Sites is the list of available sites
	Sites []models.SiteInfo

	// DefaultSite is the default site to use if no match is found
	DefaultSite models.SiteInfo

	// Enabled determines if multisite is enabled
	Enabled bool

	// UseCookieResolution enables cookie-based site resolution
	UseCookieResolution bool

	// CookieName is the name of the site cookie
	CookieName string

	// CookieSecure sets the Secure attribute
	CookieSecure bool

	// CookieHTTPOnly sets the HttpOnly attribute
	CookieHTTPOnly bool

	// CookieSameSite sets the SameSite attribute
	CookieSameSite http.SameSite
}

// MultisiteMiddleware handles multi-site resolution
type MultisiteMiddleware struct {
	config   MultisiteConfig
	resolver site.SiteResolver
}

// NewMultisiteMiddleware creates a new multisite middleware
func NewMultisiteMiddleware(config MultisiteConfig) *MultisiteMiddleware {
	// Set defaults
	if config.CookieName == "" {
		config.CookieName = "sc_site"
	}
	if !config.CookieSecure {
		config.CookieSecure = true
	}
	if !config.CookieHTTPOnly {
		config.CookieHTTPOnly = true
	}
	if config.CookieSameSite == 0 {
		config.CookieSameSite = http.SameSiteNoneMode
	}

	resolver := site.NewSiteResolver(config.Sites, config.DefaultSite)

	return &MultisiteMiddleware{
		config:   config,
		resolver: resolver,
	}
}

// Handle processes the multisite middleware
func (m *MultisiteMiddleware) Handle(ctx Context, next HandlerFunc) error {
	if !m.config.Enabled {
		debug.Multisite("multisite disabled, skipping")
		return next(ctx)
	}

	// Get the hostname from the request
	hostname := m.getHostname(ctx)
	path := ctx.Path()

	debug.Multisite("processing multisite for hostname=%s, path=%s", hostname, path)

	// Determine site name
	var siteName string
	var siteInfo *models.SiteInfo

	// Check for site query parameter first (for preview mode)
	siteParam := ctx.Request().URL.Query().Get("site")
	if siteParam != "" {
		debug.Multisite("site from query param: %s", siteParam)
		siteName = siteParam
		siteInfo, _ = m.resolver.GetByName(siteName)
	}

	// Check for site cookie (if enabled)
	if siteName == "" && m.config.UseCookieResolution {
		if cookie, err := ctx.Cookie(m.config.CookieName); err == nil && cookie != nil {
			debug.Multisite("site from cookie: %s", cookie.Value)
			siteName = cookie.Value
			siteInfo, _ = m.resolver.GetByName(siteName)
		}
	}

	// Resolve by hostname
	if siteName == "" {
		debug.Multisite("resolving site by hostname: %s", hostname)
		siteInfo, _ = m.resolver.GetByHost(hostname)
		if siteInfo != nil {
			siteName = siteInfo.Name
		}
	}

	// Fallback to default site
	if siteName == "" {
		siteName = m.config.DefaultSite.Name
		siteInfo = &m.config.DefaultSite
		debug.Multisite("using default site: %s", siteName)
	}

	// Store site in context
	ctx.Set(SiteKey, siteName)

	// Set site cookie
	ctx.SetCookie(&http.Cookie{
		Name:     m.config.CookieName,
		Value:    siteName,
		Path:     "/",
		Secure:   m.config.CookieSecure,
		HttpOnly: m.config.CookieHTTPOnly,
		SameSite: m.config.CookieSameSite,
	})

	// Rewrite the path to include site prefix
	rewritePath := client.GetSiteRewrite(path, siteName)
	ctx.Set(RewritePathKey, rewritePath)
	ctx.Set(OriginalPathKey, path)

	debug.Multisite("site resolved: %s, rewrite path: %s", siteName, rewritePath)

	return next(ctx)
}

// getHostname extracts hostname from request
func (m *MultisiteMiddleware) getHostname(ctx Context) string {
	req := ctx.Request()

	// Try X-Forwarded-Host first (for proxies)
	if host := req.Header.Get("X-Forwarded-Host"); host != "" {
		return m.normalizeHostname(host)
	}

	// Use Host header
	return m.normalizeHostname(req.Host)
}

// normalizeHostname removes port from hostname
func (m *MultisiteMiddleware) normalizeHostname(hostname string) string {
	// Remove port if present
	if idx := strings.Index(hostname, ":"); idx > 0 {
		return hostname[:idx]
	}
	return hostname
}

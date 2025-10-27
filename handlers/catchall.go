package handlers

import (
	"net/http"

	"github.com/content-sdk-go/client"
	"github.com/content-sdk-go/debug"
	"github.com/content-sdk-go/middleware"
	"github.com/content-sdk-go/models"
)

// CatchAllHandler handles all dynamic Sitecore routes
type CatchAllHandler struct {
	client *client.SitecoreClient
}

// NewCatchAllHandler creates a new catch-all handler
func NewCatchAllHandler(sitecoreClient *client.SitecoreClient) *CatchAllHandler {
	return &CatchAllHandler{
		client: sitecoreClient,
	}
}

// Handle processes the catch-all route
func (h *CatchAllHandler) Handle(ctx middleware.Context) error {
	// Get path from context (may have been rewritten by middleware)
	path := ctx.Path()

	// Get site and locale from context (set by middleware)
	site := h.getSiteFromContext(ctx)
	locale := h.getLocaleFromContext(ctx)

	debug.Layout("handling catch-all for path=%s, site=%s, locale=%s", path, site, locale)

	// Fetch page data
	page, err := h.client.GetPage(path, models.PageOptions{
		Site:   site,
		Locale: &locale,
	})

	if err != nil {
		// Check if it's a not found error
		if _, ok := err.(*models.NotFoundError); ok {
			debug.Layout("page not found: %s", path)
			return ctx.String(http.StatusNotFound, "Page not found")
		}

		// Other errors
		debug.Layout("error fetching page: %v", err)
		return ctx.String(http.StatusInternalServerError, "Internal server error")
	}

	// Return page as JSON
	// In a real application, you would render HTML here
	return ctx.JSON(http.StatusOK, page)
}

// getSiteFromContext gets the site name from context
func (h *CatchAllHandler) getSiteFromContext(ctx middleware.Context) string {
	if site := ctx.Get(middleware.SiteKey); site != nil {
		if siteStr, ok := site.(string); ok {
			return siteStr
		}
	}
	return ""
}

// getLocaleFromContext gets the locale from context
func (h *CatchAllHandler) getLocaleFromContext(ctx middleware.Context) string {
	if locale := ctx.Get(middleware.LocaleKey); locale != nil {
		if localeStr, ok := locale.(string); ok {
			return localeStr
		}
	}
	return "en"
}

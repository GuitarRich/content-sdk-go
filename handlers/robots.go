package handlers

import (
	"context"
	"net/http"

	"github.com/content-sdk-go/debug"
	"github.com/content-sdk-go/middleware"
	"github.com/content-sdk-go/seo"
)

// RobotsHandler handles robots.txt requests
type RobotsHandler struct {
	robotsService seo.RobotsService
	sitemapURLs   []string
}

// RobotsHandlerConfig contains configuration for the robots handler
type RobotsHandlerConfig struct {
	RobotsService seo.RobotsService
	SitemapURLs   []string
}

// NewRobotsHandler creates a new robots.txt handler
func NewRobotsHandler(config RobotsHandlerConfig) *RobotsHandler {
	return &RobotsHandler{
		robotsService: config.RobotsService,
		sitemapURLs:   config.SitemapURLs,
	}
}

// Handle processes robots.txt requests
func (h *RobotsHandler) Handle(ctx middleware.Context) error {
	debug.Robots("handling robots.txt request")

	// Get site from context
	site := ""
	if siteVal := ctx.Get(middleware.SiteKey); siteVal != nil {
		if siteStr, ok := siteVal.(string); ok {
			site = siteStr
		}
	}

	// Fetch robots directives from Sitecore
	directive, err := h.robotsService.FetchRobotsDirectives(context.Background(), site)
	if err != nil {
		debug.Robots("error fetching robots directives: %v", err)
		// Use default directives on error
		directive = nil
	}

	// Generate robots.txt content
	content := h.robotsService.GenerateRobotsTxt(directive, h.sitemapURLs)

	// Set content type and return
	ctx.SetHeader("Content-Type", "text/plain")
	return ctx.String(http.StatusOK, content)
}

package handlers

import (
	"context"
	"net/http"

	"github.com/guitarrich/content-sdk-go/debug"
	"github.com/guitarrich/content-sdk-go/middleware"
	"github.com/guitarrich/content-sdk-go/seo"
)

// SitemapHandler handles sitemap.xml requests
type SitemapHandler struct {
	sitemapService seo.SitemapXmlService
	sites          []string
	languages      []string
}

// SitemapHandlerConfig contains configuration for the sitemap handler
type SitemapHandlerConfig struct {
	SitemapService seo.SitemapXmlService
	Sites          []string
	Languages      []string
}

// NewSitemapHandler creates a new sitemap.xml handler
func NewSitemapHandler(config SitemapHandlerConfig) *SitemapHandler {
	return &SitemapHandler{
		sitemapService: config.SitemapService,
		sites:          config.Sites,
		languages:      config.Languages,
	}
}

// Handle processes sitemap.xml requests
func (h *SitemapHandler) Handle(ctx middleware.Context) error {
	debug.Sitemap("handling sitemap.xml request")

	// Fetch sitemap entries
	entries, err := h.sitemapService.FetchSitemap(context.Background(), h.sites, h.languages)
	if err != nil {
		debug.Sitemap("error fetching sitemap: %v", err)
		return ctx.String(http.StatusInternalServerError, "Error generating sitemap")
	}

	// Generate XML
	xml, err := h.sitemapService.GenerateSitemapXML(entries)
	if err != nil {
		debug.Sitemap("error generating sitemap XML: %v", err)
		return ctx.String(http.StatusInternalServerError, "Error generating sitemap XML")
	}

	// Set content type and return
	ctx.SetHeader("Content-Type", "application/xml")
	return ctx.String(http.StatusOK, xml)
}

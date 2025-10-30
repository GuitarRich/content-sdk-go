package seo

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/content-sdk-go/debug"
	"github.com/content-sdk-go/graphql"
	"github.com/content-sdk-go/models"
)

// SitemapXmlService generates XML sitemaps
type SitemapXmlService interface {
	FetchSitemap(ctx context.Context, sites []string, languages []string) ([]models.SitemapEntry, error)
	GenerateSitemapXML(entries []models.SitemapEntry) (string, error)
}

// SitemapXmlServiceConfig contains configuration for sitemap service
type SitemapXmlServiceConfig struct {
	GraphQLClient graphql.Client
	BaseURL       string
}

// sitemapXmlServiceImpl is the default implementation
type sitemapXmlServiceImpl struct {
	graphQLClient graphql.Client
	baseURL       string
}

// NewSitemapXmlService creates a new sitemap service
func NewSitemapXmlService(config SitemapXmlServiceConfig) SitemapXmlService {
	return &sitemapXmlServiceImpl{
		graphQLClient: config.GraphQLClient,
		baseURL:       config.BaseURL,
	}
}

// FetchSitemap fetches sitemap entries for the specified sites and languages
func (s *sitemapXmlServiceImpl) FetchSitemap(
	ctx context.Context,
	sites []string,
	languages []string,
) ([]models.SitemapEntry, error) {
	debug.Sitemap("fetching sitemap for sites=%v, languages=%v", sites, languages)

	allEntries := []models.SitemapEntry{}

	// Fetch routes for each site/language combination
	for _, site := range sites {
		for _, language := range languages {
			query := s.getSitemapQuery(site, language)

			result, err := s.graphQLClient.Request(ctx, query, nil)
			if err != nil {
				debug.Sitemap("error fetching sitemap for site=%s, language=%s: %v", site, language, err)
				continue
			}

			entries, err := s.parseSitemapResponse(result, site, language)
			if err != nil {
				debug.Sitemap("error parsing sitemap for site=%s, language=%s: %v", site, language, err)
				continue
			}

			allEntries = append(allEntries, entries...)
		}
	}

	debug.Sitemap("fetched %d sitemap entries", len(allEntries))
	return allEntries, nil
}

// GenerateSitemapXML generates XML sitemap from entries
func (s *sitemapXmlServiceImpl) GenerateSitemapXML(entries []models.SitemapEntry) (string, error) {
	// Create URL set
	urlset := &URLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]URL, 0, len(entries)),
	}

	// Convert entries to XML URLs
	for _, entry := range entries {
		url := URL{
			Loc:        entry.Loc,
			LastMod:    entry.LastMod,
			ChangeFreq: entry.ChangeFreq,
			Priority:   entry.Priority,
		}
		urlset.URLs = append(urlset.URLs, url)
	}

	// Marshal to XML
	output, err := xml.MarshalIndent(urlset, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to generate sitemap XML: %w", err)
	}

	return xml.Header + string(output), nil
}

// getSitemapQuery builds the GraphQL query for sitemap
func (s *sitemapXmlServiceImpl) getSitemapQuery(siteName, language string) string {
	return fmt.Sprintf(`
		query SitemapQuery {
			site {
				siteInfo(site: "%s") {
					routes(language: "%s") {
						path
						template
						lastModified
					}
				}
			}
		}
	`, siteName, language)
}

// parseSitemapResponse parses the sitemap response
func (s *sitemapXmlServiceImpl) parseSitemapResponse(
	data map[string]any,
	siteName, language string,
) ([]models.SitemapEntry, error) {
	entries := []models.SitemapEntry{}

	site, ok := data["site"].(map[string]any)
	if !ok {
		return entries, nil
	}

	siteInfo, ok := site["siteInfo"].(map[string]any)
	if !ok {
		return entries, nil
	}

	routes, ok := siteInfo["routes"].([]any)
	if !ok {
		return entries, nil
	}

	for _, item := range routes {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}

		path, _ := itemMap["path"].(string)
		if path == "" {
			continue
		}

		// Build full URL
		loc := s.baseURL + path

		// Get last modified date
		lastMod := ""
		if lastModified, ok := itemMap["lastModified"].(string); ok {
			lastMod = lastModified
		} else {
			// Use current date if not provided
			lastMod = time.Now().Format("2006-01-02")
		}

		entry := models.SitemapEntry{
			Loc:        loc,
			LastMod:    lastMod,
			ChangeFreq: "daily",
			Priority:   "0.5",
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// URLSet represents the XML sitemap urlset element
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

// URL represents a single sitemap URL entry
type URL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

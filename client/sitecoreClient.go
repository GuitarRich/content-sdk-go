package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/content-sdk-go/debug"
	"github.com/content-sdk-go/graphql"
	layoutservice "github.com/content-sdk-go/layoutService"
	"github.com/content-sdk-go/models"
)

const (
	// SITE_PREFIX is the prefix used to identify site names in paths
	SITE_PREFIX = "_site_"

	// PERSONALIZE_PREFIX is the prefix used for personalization variants
	PERSONALIZE_PREFIX = "_variantId_"
)

// SitecoreClient provides access to Sitecore content and services
type SitecoreClient struct {
	layoutService   *layoutservice.LayoutService
	httpClient      *http.Client
	defaultSite     string
	defaultLang     string
	graphQLEndpoint string
	graphQLAPIKey   string
}

// ClientConfig contains configuration for the Sitecore client
type ClientConfig struct {
	LayoutService   *layoutservice.LayoutService
	HTTPClient      *http.Client
	DefaultSite     string
	DefaultLanguage string
	GraphQLEndpoint string
	GraphQLAPIKey   string
}

// NewSitecoreClient creates a new Sitecore client
func NewSitecoreClient(config ClientConfig) *SitecoreClient {
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	defaultSite := config.DefaultSite
	if defaultSite == "" {
		defaultSite = "default"
	}

	defaultLang := config.DefaultLanguage
	if defaultLang == "" {
		defaultLang = "en"
	}

	return &SitecoreClient{
		layoutService:   config.LayoutService,
		httpClient:      httpClient,
		defaultSite:     defaultSite,
		defaultLang:     defaultLang,
		graphQLEndpoint: config.GraphQLEndpoint,
		graphQLAPIKey:   config.GraphQLAPIKey,
	}
}

// GetPage fetches a page from Sitecore
func (c *SitecoreClient) GetPage(path string, options models.PageOptions) (*models.Page, error) {
	// Parse and normalize the path
	normalizedPath := c.ParsePath(path)

	// Get site name (use provided or extract from path or use default)
	site := options.Site
	if site == "" {
		siteData := GetSiteRewriteData(normalizedPath, c.defaultSite)
		site = siteData.SiteName
		normalizedPath = siteData.NormalizedPath
	}

	// Get locale
	locale := options.Locale
	if locale == nil {
		locale = &c.defaultLang
	}

	// Fetch layout data
	layoutData, err := c.layoutService.FetchLayoutData(normalizedPath, layoutservice.RouteOptions{
		Site:   site,
		Locale: locale,
	}, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch layout data: %w", err)
	}

	// Check if route exists (404 handling)
	if layoutData.Sitecore.Route == nil {
		return nil, &models.NotFoundError{
			Path: path,
			Site: site,
		}
	}

	// Build page response
	page := &models.Page{
		LayoutData: layoutData,
		Dictionary: make(models.DictionaryPhrases), // TODO: Fetch dictionary
		ErrorPages: nil,                            // TODO: Fetch error pages
		HeadLinks:  []models.HTMLLink{},            // TODO: Build head links
	}

	return page, nil
}

// GetPreview fetches preview/editing data for Sitecore Pages editor
func (c *SitecoreClient) GetPreview(previewData models.PreviewData) (*models.Page, error) {
	debug.Editing("fetching preview data for item %s, language %s, site %s, mode %s",
		previewData.ItemID, previewData.Language, previewData.Site, previewData.Mode)

	// Build GraphQL query for editing data
	query := `
		query EditingQuery($itemId: String!, $language: String!, $version: String) {
			item(path: $itemId, language: $language, version: $version) {
				rendered
			}
		}
	`

	// Build variables
	variables := map[string]any{
		"itemId":   previewData.ItemID,
		"language": previewData.Language,
	}

	if previewData.Version != "" {
		variables["version"] = previewData.Version
	}

	// Create custom headers for editing mode
	headers := make(map[string]string)
	headers["sc_layoutKind"] = string(previewData.LayoutKind)

	// Set edit mode header
	if previewData.Mode == models.PreviewModeEdit || previewData.Mode == models.PreviewModeMetadata {
		headers["sc_editMode"] = "true"
	} else {
		headers["sc_editMode"] = "false"
	}

	debug.Editing("GraphQL headers: %+v", headers)

	// Create a new GraphQL client config with custom headers
	clientConfig := &graphql.ClientConfig{
		Retries:    3,
		Timeout:    30 * time.Second,
		RetryDelay: 1 * time.Second,
		Headers:    headers,
	}

	// Create a temporary GraphQL client with custom headers
	graphQLClient := graphql.NewClient(
		c.graphQLEndpoint,
		c.graphQLAPIKey,
		c.httpClient,
		clientConfig,
	)

	// Execute GraphQL request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := graphQLClient.Request(ctx, query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch editing data: %w", err)
	}

	// Extract the rendered layout data
	item, ok := result["item"].(map[string]any)
	if !ok || item == nil {
		return nil, fmt.Errorf("item not found in GraphQL response")
	}

	rendered, ok := item["rendered"].(map[string]any)
	if !ok || rendered == nil {
		return nil, fmt.Errorf("rendered data not found in item")
	}

	// Parse rendered data into LayoutServiceData struct
	// This is necessary for the renderer to work correctly
	layoutData, err := c.parseRenderedData(rendered)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rendered data: %w", err)
	}

	// Create page response
	page := &models.Page{
		LayoutData: layoutData,
		Language:   previewData.Language,
		Site:       previewData.Site,
		ItemID:     previewData.ItemID,
		Path:       previewData.Route,
		EditingContext: &models.EditingContext{
			IsEditing: previewData.Mode == models.PreviewModeEdit || previewData.Mode == models.PreviewModeMetadata,
			IsPreview: previewData.Mode == models.PreviewModePreview,
			Mode:      models.PageModeEdit,
			QueryParams: map[string]string{
				"sc_itemid":     previewData.ItemID,
				"sc_lang":       previewData.Language,
				"sc_site":       previewData.Site,
				"sc_layoutKind": string(previewData.LayoutKind),
				"mode":          string(previewData.Mode),
			},
		},
	}

	debug.Editing("preview page fetched successfully")

	return page, nil
}

// GetDesignLibraryData fetches design library component data
func (c *SitecoreClient) GetDesignLibraryData(data models.DesignLibraryRenderPreviewData) (*models.Page, error) {
	// TODO: Implement design library data fetching
	return nil, fmt.Errorf("design library mode not yet implemented")
}

// GetStaticPaths generates static paths for all pages in given sites and languages
func (c *SitecoreClient) GetStaticPaths(sites []string, languages []string) ([]models.StaticPath, error) {
	// TODO: Implement static path generation
	// This requires querying all routes from Sitecore via GraphQL
	return nil, fmt.Errorf("static path generation not yet implemented")
}

// GetSiteNameFromPath extracts the site name from a path
func (c *SitecoreClient) GetSiteNameFromPath(path string) string {
	normalizedPath := c.ParsePath(path)
	siteData := GetSiteRewriteData(normalizedPath, c.defaultSite)
	return siteData.SiteName
}

// ParsePath normalizes a path (string or []string)
func (c *SitecoreClient) ParsePath(path any) string {
	normalized := parsePath(path)

	// Remove site rewrite prefix
	normalized = NormalizeSiteRewrite(normalized)

	// Remove personalization prefix
	normalized = NormalizePersonalizedRewrite(normalized)

	return normalized
}

// parsePath normalizes path regardless of type
// Accepts either string or []string and returns a normalized string path
func parsePath(path any) string {
	switch p := path.(type) {
	case string:
		// If string starts with '/', return as-is, otherwise prepend '/'
		if strings.HasPrefix(p, "/") {
			return p
		}
		return "/" + p
	case []string:
		// Filter out '/' parts and trim slashes from each part
		var parts []string
		for _, part := range p {
			if part == "/" {
				continue
			}
			// Remove leading and trailing slashes
			trimmed := strings.Trim(part, "/")
			if trimmed != "" {
				parts = append(parts, trimmed)
			}
		}
		// Join with '/' and prepend '/'
		return "/" + strings.Join(parts, "/")
	default:
		return "/"
	}
}

// GetSiteRewrite adds site prefix to a path
func GetSiteRewrite(path string, siteName string) string {
	// Normalize path first
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Add site prefix
	return fmt.Sprintf("/%s%s%s", SITE_PREFIX, siteName, path)
}

// GetSiteRewriteData extracts site information from a path
func GetSiteRewriteData(path string, defaultSite string) models.SiteRewriteData {
	// Check if path contains site prefix
	if strings.Contains(path, SITE_PREFIX) {
		parts := strings.Split(path, "/")
		for i, part := range parts {
			if after, ok := strings.CutPrefix(part, SITE_PREFIX); ok {
				siteName := after
				// Remove the site part from path
				remainingParts := append(parts[:i], parts[i+1:]...)
				normalizedPath := "/" + strings.Join(remainingParts, "/")
				normalizedPath = strings.ReplaceAll(normalizedPath, "//", "/")

				return models.SiteRewriteData{
					SiteName:       siteName,
					NormalizedPath: normalizedPath,
				}
			}
		}
	}

	// No site prefix found, use default
	return models.SiteRewriteData{
		SiteName:       defaultSite,
		NormalizedPath: path,
	}
}

// NormalizeSiteRewrite removes site prefix from a path
func NormalizeSiteRewrite(path string) string {
	if !strings.Contains(path, SITE_PREFIX) {
		return path
	}

	parts := strings.Split(path, "/")
	var normalizedParts []string

	for _, part := range parts {
		if !strings.HasPrefix(part, SITE_PREFIX) {
			normalizedParts = append(normalizedParts, part)
		}
	}

	result := "/" + strings.Join(normalizedParts, "/")
	return strings.ReplaceAll(result, "//", "/")
}

// GetPersonalizedRewrite adds personalization variant prefix to a path
func GetPersonalizedRewrite(path string, variantId string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return fmt.Sprintf("/%s%s%s", PERSONALIZE_PREFIX, variantId, path)
}

// GetPersonalizedRewriteData extracts personalization variant from a path
func GetPersonalizedRewriteData(path string) models.PersonalizeRewriteData {
	if !strings.Contains(path, PERSONALIZE_PREFIX) {
		return models.PersonalizeRewriteData{
			VariantId:      "",
			NormalizedPath: path,
		}
	}

	parts := strings.Split(path, "/")
	for i, part := range parts {
		if after, ok := strings.CutPrefix(part, PERSONALIZE_PREFIX); ok {
			variantId := after
			// Remove the variant part from path
			remainingParts := append(parts[:i], parts[i+1:]...)
			normalizedPath := "/" + strings.Join(remainingParts, "/")
			normalizedPath = strings.ReplaceAll(normalizedPath, "//", "/")

			return models.PersonalizeRewriteData{
				VariantId:      variantId,
				NormalizedPath: normalizedPath,
			}
		}
	}

	return models.PersonalizeRewriteData{
		VariantId:      "",
		NormalizedPath: path,
	}
}

// NormalizePersonalizedRewrite removes personalization variant prefix from a path
func NormalizePersonalizedRewrite(path string) string {
	if !strings.Contains(path, PERSONALIZE_PREFIX) {
		return path
	}

	parts := strings.Split(path, "/")
	var normalizedParts []string

	for _, part := range parts {
		if !strings.HasPrefix(part, PERSONALIZE_PREFIX) {
			normalizedParts = append(normalizedParts, part)
		}
	}

	result := "/" + strings.Join(normalizedParts, "/")
	return strings.ReplaceAll(result, "//", "/")
}

// parseRenderedData converts a map[string]any rendered response into LayoutServiceData
func (c *SitecoreClient) parseRenderedData(rendered map[string]any) (*layoutservice.LayoutServiceData, error) {
	// Marshal the map to JSON
	jsonBytes, err := json.Marshal(rendered)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rendered data: %w", err)
	}

	// Unmarshal into LayoutServiceData struct
	var layoutData layoutservice.LayoutServiceData
	if err := json.Unmarshal(jsonBytes, &layoutData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal layout data: %w", err)
	}

	return &layoutData, nil
}

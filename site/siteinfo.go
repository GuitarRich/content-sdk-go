package site

import (
	"context"
	"fmt"

	"github.com/guitarrich/content-sdk-go/debug"
	"github.com/guitarrich/content-sdk-go/graphql"
	"github.com/guitarrich/content-sdk-go/models"
)

// SiteInfoService fetches site configuration from Sitecore
type SiteInfoService interface {
	FetchSiteInfo(ctx context.Context, siteName string) (*models.SiteInfo, error)
	FetchSites(ctx context.Context) ([]models.SiteInfo, error)
}

// SiteInfoServiceConfig contains configuration for the site info service
type SiteInfoServiceConfig struct {
	GraphQLClient graphql.Client
}

// siteInfoServiceImpl is the default implementation
type siteInfoServiceImpl struct {
	graphQLClient graphql.Client
}

// NewSiteInfoService creates a new site info service
func NewSiteInfoService(config SiteInfoServiceConfig) SiteInfoService {
	return &siteInfoServiceImpl{
		graphQLClient: config.GraphQLClient,
	}
}

// FetchSiteInfo fetches information for a specific site
func (s *siteInfoServiceImpl) FetchSiteInfo(ctx context.Context, siteName string) (*models.SiteInfo, error) {
	debug.Multisite("fetching site info for %s", siteName)

	query := s.getSiteInfoQuery(siteName)

	result, err := s.graphQLClient.Request(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch site info: %w", err)
	}

	siteInfo, err := s.parseSiteInfoResponse(result, siteName)
	if err != nil {
		return nil, fmt.Errorf("failed to parse site info response: %w", err)
	}

	return siteInfo, nil
}

// FetchSites fetches all available sites
func (s *siteInfoServiceImpl) FetchSites(ctx context.Context) ([]models.SiteInfo, error) {
	debug.Multisite("fetching all sites")

	query := s.getAllSitesQuery()

	result, err := s.graphQLClient.Request(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sites: %w", err)
	}

	sites, err := s.parseAllSitesResponse(result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sites response: %w", err)
	}

	debug.Multisite("fetched %d sites", len(sites))
	return sites, nil
}

// getSiteInfoQuery builds the GraphQL query for a single site
func (s *siteInfoServiceImpl) getSiteInfoQuery(siteName string) string {
	return fmt.Sprintf(`
		query SiteInfoQuery {
			site {
				siteInfo(site: "%s") {
					name
					hostName
					language
					rootPath
					database
				}
			}
		}
	`, siteName)
}

// getAllSitesQuery builds the GraphQL query for all sites
func (s *siteInfoServiceImpl) getAllSitesQuery() string {
	return `
		query AllSitesQuery {
			site {
				siteInfoCollection {
					name
					hostName
					language
					rootPath
					database
				}
			}
		}
	`
}

// parseSiteInfoResponse parses a single site info response
func (s *siteInfoServiceImpl) parseSiteInfoResponse(data map[string]any, siteName string) (*models.SiteInfo, error) {
	site, ok := data["site"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("site not found in response")
	}

	siteInfo, ok := site["siteInfo"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("siteInfo not found for site %s", siteName)
	}

	return s.mapToSiteInfo(siteInfo), nil
}

// parseAllSitesResponse parses all sites response
func (s *siteInfoServiceImpl) parseAllSitesResponse(data map[string]any) ([]models.SiteInfo, error) {
	site, ok := data["site"].(map[string]any)
	if !ok {
		return []models.SiteInfo{}, nil
	}

	siteInfoCollection, ok := site["siteInfoCollection"].([]any)
	if !ok {
		return []models.SiteInfo{}, nil
	}

	sites := make([]models.SiteInfo, 0, len(siteInfoCollection))
	for _, item := range siteInfoCollection {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		sites = append(sites, *s.mapToSiteInfo(itemMap))
	}

	return sites, nil
}

// mapToSiteInfo converts GraphQL response to SiteInfo model
func (s *siteInfoServiceImpl) mapToSiteInfo(data map[string]any) *models.SiteInfo {
	siteInfo := &models.SiteInfo{}

	if name, ok := data["name"].(string); ok {
		siteInfo.Name = name
	}

	if hostName, ok := data["hostName"].(string); ok {
		siteInfo.HostName = hostName
	}

	if language, ok := data["language"].(string); ok {
		siteInfo.Language = language
	}

	if rootPath, ok := data["rootPath"].(string); ok {
		siteInfo.RootPath = rootPath
	}

	if database, ok := data["database"].(string); ok {
		siteInfo.Database = database
	}

	return siteInfo
}

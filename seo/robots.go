package seo

import (
	"context"
	"fmt"
	"strings"

	"github.com/guitarrich/content-sdk-go/debug"
	"github.com/guitarrich/content-sdk-go/graphql"
	"github.com/guitarrich/content-sdk-go/models"
)

// RobotsService generates robots.txt content
type RobotsService interface {
	FetchRobotsDirectives(ctx context.Context, siteName string) (*models.RobotsDirective, error)
	GenerateRobotsTxt(directive *models.RobotsDirective, sitemapURLs []string) string
}

// RobotsServiceConfig contains configuration for robots service
type RobotsServiceConfig struct {
	GraphQLClient graphql.Client
}

// robotsServiceImpl is the default implementation
type robotsServiceImpl struct {
	graphQLClient graphql.Client
}

// NewRobotsService creates a new robots service
func NewRobotsService(config RobotsServiceConfig) RobotsService {
	return &robotsServiceImpl{
		graphQLClient: config.GraphQLClient,
	}
}

// FetchRobotsDirectives fetches robots directives from Sitecore
func (s *robotsServiceImpl) FetchRobotsDirectives(
	ctx context.Context,
	siteName string,
) (*models.RobotsDirective, error) {
	debug.Robots("fetching robots directives for site %s", siteName)

	query := s.getRobotsQuery(siteName)

	result, err := s.graphQLClient.Request(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch robots directives: %w", err)
	}

	directive, err := s.parseRobotsResponse(result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse robots response: %w", err)
	}

	return directive, nil
}

// GenerateRobotsTxt generates robots.txt content
func (s *robotsServiceImpl) GenerateRobotsTxt(
	directive *models.RobotsDirective,
	sitemapURLs []string,
) string {
	var builder strings.Builder

	// If custom content is provided, use it
	if directive != nil && directive.Content != "" {
		return directive.Content
	}

	// Generate default robots.txt
	userAgent := "*"
	if directive != nil && directive.UserAgent != "" {
		userAgent = directive.UserAgent
	}

	builder.WriteString(fmt.Sprintf("User-agent: %s\n", userAgent))

	// Add disallow rules
	if directive != nil && len(directive.Disallow) > 0 {
		for _, path := range directive.Disallow {
			builder.WriteString(fmt.Sprintf("Disallow: %s\n", path))
		}
	} else {
		// Default: allow all
		builder.WriteString("Disallow:\n")
	}

	// Add allow rules
	if directive != nil && len(directive.Allow) > 0 {
		for _, path := range directive.Allow {
			builder.WriteString(fmt.Sprintf("Allow: %s\n", path))
		}
	}

	// Add sitemap URLs
	if directive != nil && len(directive.Sitemap) > 0 {
		builder.WriteString("\n")
		for _, sitemapURL := range directive.Sitemap {
			builder.WriteString(fmt.Sprintf("Sitemap: %s\n", sitemapURL))
		}
	} else if len(sitemapURLs) > 0 {
		builder.WriteString("\n")
		for _, sitemapURL := range sitemapURLs {
			builder.WriteString(fmt.Sprintf("Sitemap: %s\n", sitemapURL))
		}
	}

	return builder.String()
}

// getRobotsQuery builds the GraphQL query for robots
func (s *robotsServiceImpl) getRobotsQuery(siteName string) string {
	return fmt.Sprintf(`
		query RobotsQuery {
			site {
				siteInfo(site: "%s") {
					robots {
						content
						userAgent
						allow
						disallow
						sitemap
					}
				}
			}
		}
	`, siteName)
}

// parseRobotsResponse parses the robots response
func (s *robotsServiceImpl) parseRobotsResponse(data map[string]any) (*models.RobotsDirective, error) {
	site, ok := data["site"].(map[string]any)
	if !ok {
		return &models.RobotsDirective{}, nil
	}

	siteInfo, ok := site["siteInfo"].(map[string]any)
	if !ok {
		return &models.RobotsDirective{}, nil
	}

	robots, ok := siteInfo["robots"].(map[string]any)
	if !ok {
		return &models.RobotsDirective{}, nil
	}

	directive := &models.RobotsDirective{}

	if content, ok := robots["content"].(string); ok {
		directive.Content = content
	}

	if userAgent, ok := robots["userAgent"].(string); ok {
		directive.UserAgent = userAgent
	}

	if allow, ok := robots["allow"].([]any); ok {
		directive.Allow = make([]string, 0, len(allow))
		for _, item := range allow {
			if str, ok := item.(string); ok {
				directive.Allow = append(directive.Allow, str)
			}
		}
	}

	if disallow, ok := robots["disallow"].([]any); ok {
		directive.Disallow = make([]string, 0, len(disallow))
		for _, item := range disallow {
			if str, ok := item.(string); ok {
				directive.Disallow = append(directive.Disallow, str)
			}
		}
	}

	if sitemap, ok := robots["sitemap"].([]any); ok {
		directive.Sitemap = make([]string, 0, len(sitemap))
		for _, item := range sitemap {
			if str, ok := item.(string); ok {
				directive.Sitemap = append(directive.Sitemap, str)
			}
		}
	}

	return directive, nil
}

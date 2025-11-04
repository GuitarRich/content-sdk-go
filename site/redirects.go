package site

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/guitarrich/content-sdk-go/debug"
	"github.com/guitarrich/content-sdk-go/graphql"
	"github.com/guitarrich/content-sdk-go/models"
)

// RedirectsService fetches and matches URL redirects
type RedirectsService interface {
	FetchRedirects(ctx context.Context, siteName string) ([]models.RedirectInfo, error)
	GetRedirect(path string, redirects []models.RedirectInfo) (*models.RedirectInfo, error)
}

// RedirectsServiceConfig contains configuration for the redirects service
type RedirectsServiceConfig struct {
	GraphQLClient graphql.Client
}

// redirectsServiceImpl is the default implementation
type redirectsServiceImpl struct {
	graphQLClient graphql.Client
}

// NewRedirectsService creates a new redirects service
func NewRedirectsService(config RedirectsServiceConfig) RedirectsService {
	return &redirectsServiceImpl{
		graphQLClient: config.GraphQLClient,
	}
}

// FetchRedirects fetches all redirects for a site
func (s *redirectsServiceImpl) FetchRedirects(ctx context.Context, siteName string) ([]models.RedirectInfo, error) {
	debug.Redirects("fetching redirects for site %s", siteName)

	query := s.getRedirectsQuery(siteName)

	result, err := s.graphQLClient.Request(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch redirects: %w", err)
	}

	redirects, err := s.parseRedirectsResponse(result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redirects response: %w", err)
	}

	debug.Redirects("fetched %d redirects", len(redirects))
	return redirects, nil
}

// GetRedirect finds a matching redirect for a path
func (s *redirectsServiceImpl) GetRedirect(path string, redirects []models.RedirectInfo) (*models.RedirectInfo, error) {
	// Normalize path
	normalizedPath := strings.TrimSpace(path)
	if !strings.HasPrefix(normalizedPath, "/") {
		normalizedPath = "/" + normalizedPath
	}

	// Try exact matches first
	for _, redirect := range redirects {
		if !redirect.IsRegex && redirect.Pattern == normalizedPath {
			return &redirect, nil
		}
	}

	// Try regex matches
	for _, redirect := range redirects {
		if redirect.IsRegex {
			matched, err := regexp.MatchString(redirect.Pattern, normalizedPath)
			if err != nil {
				debug.Redirects("invalid regex pattern: %s", redirect.Pattern)
				continue
			}
			if matched {
				return &redirect, nil
			}
		}
	}

	return nil, nil // No redirect found
}

// getRedirectsQuery builds the GraphQL query for redirects
func (s *redirectsServiceImpl) getRedirectsQuery(siteName string) string {
	return fmt.Sprintf(`
		query RedirectsQuery {
			site {
				siteInfo(site: "%s") {
					redirects {
						pattern
						target
						redirectType
						locale
						isRegex
					}
				}
			}
		}
	`, siteName)
}

// parseRedirectsResponse parses the redirects response
func (s *redirectsServiceImpl) parseRedirectsResponse(data map[string]any) ([]models.RedirectInfo, error) {
	redirects := []models.RedirectInfo{}

	site, ok := data["site"].(map[string]any)
	if !ok {
		return redirects, nil
	}

	siteInfo, ok := site["siteInfo"].(map[string]any)
	if !ok {
		return redirects, nil
	}

	redirectsList, ok := siteInfo["redirects"].([]any)
	if !ok {
		return redirects, nil
	}

	for _, item := range redirectsList {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}

		redirect := models.RedirectInfo{}

		if pattern, ok := itemMap["pattern"].(string); ok {
			redirect.Pattern = pattern
		}

		if target, ok := itemMap["target"].(string); ok {
			redirect.Target = target
		}

		if redirectType, ok := itemMap["redirectType"].(string); ok {
			redirect.RedirectType = models.RedirectType(redirectType)
		}

		if locale, ok := itemMap["locale"].(string); ok {
			redirect.Locale = locale
		}

		if isRegex, ok := itemMap["isRegex"].(bool); ok {
			redirect.IsRegex = isRegex
		}

		redirects = append(redirects, redirect)
	}

	return redirects, nil
}

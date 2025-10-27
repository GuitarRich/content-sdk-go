package seo

import (
	"context"
	"fmt"

	"github.com/content-sdk-go/debug"
	"github.com/content-sdk-go/graphql"
	"github.com/content-sdk-go/models"
)

// ErrorPagesService fetches custom error page definitions
type ErrorPagesService interface {
	FetchErrorPages(ctx context.Context, siteName string) (*models.ErrorPages, error)
}

// ErrorPagesServiceConfig contains configuration for error pages service
type ErrorPagesServiceConfig struct {
	GraphQLClient graphql.Client
}

// errorPagesServiceImpl is the default implementation
type errorPagesServiceImpl struct {
	graphQLClient graphql.Client
}

// NewErrorPagesService creates a new error pages service
func NewErrorPagesService(config ErrorPagesServiceConfig) ErrorPagesService {
	return &errorPagesServiceImpl{
		graphQLClient: config.GraphQLClient,
	}
}

// FetchErrorPages fetches custom error page definitions
func (s *errorPagesServiceImpl) FetchErrorPages(
	ctx context.Context,
	siteName string,
) (*models.ErrorPages, error) {
	debug.ErrorPages("fetching error pages for site %s", siteName)

	query := s.getErrorPagesQuery(siteName)

	result, err := s.graphQLClient.Request(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch error pages: %w", err)
	}

	errorPages, err := s.parseErrorPagesResponse(result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse error pages response: %w", err)
	}

	return errorPages, nil
}

// getErrorPagesQuery builds the GraphQL query for error pages
func (s *errorPagesServiceImpl) getErrorPagesQuery(siteName string) string {
	return fmt.Sprintf(`
		query ErrorPagesQuery {
			site {
				siteInfo(site: "%s") {
					errorHandling {
						notFoundPage {
							rendered
						}
						serverErrorPage {
							rendered
						}
					}
				}
			}
		}
	`, siteName)
}

// parseErrorPagesResponse parses the error pages response
func (s *errorPagesServiceImpl) parseErrorPagesResponse(data map[string]interface{}) (*models.ErrorPages, error) {
	errorPages := &models.ErrorPages{}

	site, ok := data["site"].(map[string]interface{})
	if !ok {
		return errorPages, nil
	}

	siteInfo, ok := site["siteInfo"].(map[string]interface{})
	if !ok {
		return errorPages, nil
	}

	errorHandling, ok := siteInfo["errorHandling"].(map[string]interface{})
	if !ok {
		return errorPages, nil
	}

	// Parse not found page
	if notFoundPage, ok := errorHandling["notFoundPage"].(map[string]interface{}); ok {
		if rendered, ok := notFoundPage["rendered"].(map[string]interface{}); ok {
			errorPages.NotFoundPage = rendered
		}
	}

	// Parse server error page
	if serverErrorPage, ok := errorHandling["serverErrorPage"].(map[string]interface{}); ok {
		if rendered, ok := serverErrorPage["rendered"].(map[string]interface{}); ok {
			errorPages.ServerErrorPage = rendered
		}
	}

	return errorPages, nil
}

package i18n

import (
	"context"
	"fmt"

	"github.com/content-sdk-go/debug"
	"github.com/content-sdk-go/graphql"
	"github.com/content-sdk-go/models"
)

// DictionaryService fetches dictionary phrases for internationalization
type DictionaryService interface {
	FetchDictionaryData(ctx context.Context, locale, siteName string) (models.DictionaryPhrases, error)
}

// DictionaryServiceConfig contains configuration for the dictionary service
type DictionaryServiceConfig struct {
	GraphQLClient graphql.Client
	SiteName      string
}

// dictionaryServiceImpl is the default implementation
type dictionaryServiceImpl struct {
	graphQLClient graphql.Client
	siteName      string
}

// NewDictionaryService creates a new dictionary service
func NewDictionaryService(config DictionaryServiceConfig) DictionaryService {
	return &dictionaryServiceImpl{
		graphQLClient: config.GraphQLClient,
		siteName:      config.SiteName,
	}
}

// FetchDictionaryData fetches all dictionary phrases for a given locale
func (s *dictionaryServiceImpl) FetchDictionaryData(
	ctx context.Context,
	locale string,
	siteName string,
) (models.DictionaryPhrases, error) {
	debug.Dictionary("fetching dictionary for locale=%s, site=%s", locale, siteName)

	// Use provided site name or fall back to config
	site := siteName
	if site == "" {
		site = s.siteName
	}

	// Build GraphQL query
	query := s.getDictionaryQuery(site, locale)

	// Execute query
	result, err := s.graphQLClient.Request(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dictionary data: %w", err)
	}

	// Parse response
	phrases, err := s.parseDictionaryResponse(result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dictionary response: %w", err)
	}

	debug.Dictionary("fetched %d dictionary phrases", len(phrases))
	return phrases, nil
}

// getDictionaryQuery builds the GraphQL query for fetching dictionary
func (s *dictionaryServiceImpl) getDictionaryQuery(siteName, locale string) string {
	return fmt.Sprintf(`
		query DictionaryQuery {
			site {
				siteInfo(site: "%s") {
					dictionary(language: "%s") {
						key
						value
					}
				}
			}
		}
	`, siteName, locale)
}

// parseDictionaryResponse parses the GraphQL response into DictionaryPhrases
func (s *dictionaryServiceImpl) parseDictionaryResponse(data map[string]any) (models.DictionaryPhrases, error) {
	phrases := make(models.DictionaryPhrases)

	// Navigate through response structure
	site, ok := data["site"].(map[string]any)
	if !ok {
		return phrases, nil // Return empty if no site data
	}

	siteInfo, ok := site["siteInfo"].(map[string]any)
	if !ok {
		return phrases, nil // Return empty if no site info
	}

	dictionary, ok := siteInfo["dictionary"].([]any)
	if !ok {
		return phrases, nil // Return empty if no dictionary
	}

	// Parse each dictionary entry
	for _, entry := range dictionary {
		entryMap, ok := entry.(map[string]any)
		if !ok {
			continue
		}

		key, hasKey := entryMap["key"].(string)
		value, hasValue := entryMap["value"].(string)

		if hasKey && hasValue {
			phrases[key] = value
		}
	}

	return phrases, nil
}

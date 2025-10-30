package layoutservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/content-sdk-go/debug"
	"github.com/content-sdk-go/graphql"
)

// GraphQLLayoutQueryName is the name of the GraphQL query for layout data
const GraphQLLayoutQueryName = "ContentSdkLayoutQuery"

// FetchOptions represents options for fetch operations
type FetchOptions struct {
	Retries    *int
	HTTPClient *http.Client
	Headers    map[string]string
	Timeout    *time.Duration
}

// Note: GraphQLClient is now defined in graphql package
// Keeping this as alias for backward compatibility
type GraphQLClient = graphql.Client

// GraphQLServiceConfig contains configuration for GraphQL service
type GraphQLServiceConfig struct {
	Endpoint   string
	APIKey     string
	HTTPClient *http.Client
}

// LayoutServiceConfig contains configuration for the Layout Service
type LayoutServiceConfig struct {
	GraphQLServiceConfig
	// FormatLayoutQuery is an optional function to customize the layout query
	FormatLayoutQuery func(site, itemPath string, language *string) string
}

// LayoutService fetches layout data using Sitecore's GraphQL API
type LayoutService struct {
	serviceConfig LayoutServiceConfig
	graphQLClient GraphQLClient
}

// Note: defaultGraphQLClient removed - now using graphql.Client

// NewLayoutService creates a new LayoutService instance
func NewLayoutService(serviceConfig LayoutServiceConfig) *LayoutService {
	// Create GraphQL client using factory
	graphQLClient := graphql.NewClient(
		serviceConfig.Endpoint,
		serviceConfig.APIKey,
		serviceConfig.HTTPClient,
		nil, // Use default config
	)

	return &LayoutService{
		serviceConfig: serviceConfig,
		graphQLClient: graphQLClient,
	}
}

// NewLayoutServiceWithClient creates a new LayoutService with a custom GraphQL client
func NewLayoutServiceWithClient(serviceConfig LayoutServiceConfig, graphQLClient GraphQLClient) *LayoutService {
	return &LayoutService{
		serviceConfig: serviceConfig,
		graphQLClient: graphQLClient,
	}
}

// FetchLayoutData fetches layout data for an item
// Parameters:
//   - itemPath: item path to fetch layout data for
//   - routeOptions: Request options like language and site to retrieve data for
//   - fetchOptions: Options to override graphQL client details like retries and fetch implementation
//
// Returns: layout service data
func (ls *LayoutService) FetchLayoutData(
	itemPath string,
	routeOptions RouteOptions,
	fetchOptions *FetchOptions,
) (*LayoutServiceData, error) {
	site := routeOptions.Site
	query := ls.getLayoutQuery(itemPath, site, routeOptions.Locale)

	localeStr := ""
	if routeOptions.Locale != nil {
		localeStr = *routeOptions.Locale
	}
	debug.Layout("fetching layout data for %s %s %s", itemPath, localeStr, site)

	// Create context with timeout if specified in fetchOptions
	ctx := context.Background()
	if fetchOptions != nil && fetchOptions.Timeout != nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *fetchOptions.Timeout)
		defer cancel()
	}

	data, err := ls.graphQLClient.Request(ctx, query, map[string]any{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch layout data: %w", err)
	}

	// Parse the response
	layoutData, err := ls.parseLayoutResponse(data, routeOptions.Locale)
	if err != nil {
		return nil, fmt.Errorf("failed to parse layout response: %w", err)
	}

	return layoutData, nil
}

// parseLayoutResponse parses the GraphQL response into LayoutServiceData
func (ls *LayoutService) parseLayoutResponse(data map[string]any, locale *string) (*LayoutServiceData, error) {
	// Navigate through the response structure: data.layout.item.rendered
	layout, ok := data["layout"].(map[string]any)
	if !ok {
		// If `rendered` is empty -> not found, return default structure
		return ls.createDefaultLayoutData(locale), nil
	}

	item, ok := layout["item"].(map[string]any)
	if !ok {
		return ls.createDefaultLayoutData(locale), nil
	}

	rendered, ok := item["rendered"].(map[string]any)
	if !ok || rendered == nil {
		return ls.createDefaultLayoutData(locale), nil
	}

	// Convert map to LayoutServiceData struct
	jsonBytes, err := json.Marshal(rendered)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rendered data: %w", err)
	}

	var layoutData LayoutServiceData
	if err := json.Unmarshal(jsonBytes, &layoutData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal layout data: %w", err)
	}

	return &layoutData, nil
}

// createDefaultLayoutData creates a default LayoutServiceData when item is not found
func (ls *LayoutService) createDefaultLayoutData(locale *string) *LayoutServiceData {
	pageEditing := false
	return &LayoutServiceData{
		Sitecore: struct {
			LayoutServiceContextData
			Route *RouteData `json:"route"`
		}{
			LayoutServiceContextData: LayoutServiceContextData{
				Context: LayoutServiceContext{
					PageEditing: &pageEditing,
					Language:    locale,
				},
			},
			Route: nil,
		},
	}
}

// getLayoutQuery returns GraphQL Layout query
// Parameters:
//   - itemPath: page route
//   - site: site name
//   - language: language (optional)
//
// Returns: GraphQL query
func (ls *LayoutService) getLayoutQuery(itemPath string, site string, language *string) string {
	languageVariable := ""
	if language != nil && *language != "" {
		languageVariable = fmt.Sprintf(`, language:"%s"`, *language)
	}

	var layoutQuery string
	if ls.serviceConfig.FormatLayoutQuery != nil {
		layoutQuery = ls.serviceConfig.FormatLayoutQuery(site, itemPath, language)
	} else {
		layoutQuery = fmt.Sprintf(`layout(site:"%s", routePath:"%s"%s)`, site, itemPath, languageVariable)
	}

	return fmt.Sprintf(`query %s {
      %s{
        item {
          rendered
        }
      }
    }`, GraphQLLayoutQueryName, layoutQuery)
}

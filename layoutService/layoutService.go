package layoutservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/content-sdk-go/debug"
)

// GraphQLLayoutQueryName is the name of the GraphQL query for layout data
const GraphQLLayoutQueryName = "ContentSdkLayoutQuery"

// FetchOptions represents options for fetch operations
type FetchOptions struct {
	Retries    *int
	HTTPClient *http.Client
	Headers    map[string]string
	Timeout    *int
}

// GraphQLClient is an interface for making GraphQL requests
type GraphQLClient interface {
	Request(query string, variables map[string]interface{}, options *FetchOptions) (map[string]interface{}, error)
}

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

// defaultGraphQLClient is a simple implementation of GraphQLClient
type defaultGraphQLClient struct {
	endpoint   string
	apiKey     string
	httpClient *http.Client
}

// NewLayoutService creates a new LayoutService instance
func NewLayoutService(serviceConfig LayoutServiceConfig) *LayoutService {
	// Use provided HTTP client or default
	httpClient := serviceConfig.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	// Create default GraphQL client
	graphQLClient := &defaultGraphQLClient{
		endpoint:   serviceConfig.Endpoint,
		apiKey:     serviceConfig.APIKey,
		httpClient: httpClient,
	}

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

	data, err := ls.graphQLClient.Request(query, map[string]interface{}{}, fetchOptions)
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
func (ls *LayoutService) parseLayoutResponse(data map[string]interface{}, locale *string) (*LayoutServiceData, error) {
	// Navigate through the response structure: data.layout.item.rendered
	layout, ok := data["layout"].(map[string]interface{})
	if !ok {
		// If `rendered` is empty -> not found, return default structure
		return ls.createDefaultLayoutData(locale), nil
	}

	item, ok := layout["item"].(map[string]interface{})
	if !ok {
		return ls.createDefaultLayoutData(locale), nil
	}

	rendered, ok := item["rendered"].(map[string]interface{})
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

// Request implements GraphQLClient interface for defaultGraphQLClient
func (c *defaultGraphQLClient) Request(
	query string,
	variables map[string]interface{},
	options *FetchOptions,
) (map[string]interface{}, error) {
	// Prepare the GraphQL request body
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GraphQL request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("sc_apikey", c.apiKey)
	}

	// Apply custom headers from options
	if options != nil && options.Headers != nil {
		for key, value := range options.Headers {
			req.Header.Set(key, value)
		}
	}

	// Use custom HTTP client from options if provided
	httpClient := c.httpClient
	if options != nil && options.HTTPClient != nil {
		httpClient = options.HTTPClient
	}

	// Execute request
	debug.Http("GraphQL request to %s", c.endpoint)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GraphQL request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GraphQL request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result struct {
		Data   map[string]interface{} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors,omitempty"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GraphQL response: %w", err)
	}

	// Check for GraphQL errors
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", result.Errors)
	}

	return result.Data, nil
}

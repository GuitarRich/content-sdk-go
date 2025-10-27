package config

import (
	"fmt"
	"time"

	"github.com/content-sdk-go/models"
	"github.com/content-sdk-go/utils"
)

// Config contains all configuration for the Sitecore Content SDK
type Config struct {
	// API configuration
	API APIConfig `json:"api"`

	// DefaultSite is the default site name
	DefaultSite string `json:"defaultSite"`

	// DefaultLanguage is the default language
	DefaultLanguage string `json:"defaultLanguage"`

	// Multisite configuration
	Multisite MultisiteConfig `json:"multisite"`

	// Personalization configuration
	Personalize PersonalizeConfig `json:"personalize"`

	// Editing configuration
	Editing EditingConfig `json:"editing"`

	// Timeouts
	EdgeTimeout time.Duration `json:"edgeTimeout"`
	CDPTimeout  time.Duration `json:"cdpTimeout"`

	// EnableDebug enables debug logging
	EnableDebug bool `json:"enableDebug"`
}

// APIConfig contains API configuration
type APIConfig struct {
	// Edge API configuration (recommended for production)
	Edge EdgeAPIConfig `json:"edge"`

	// Local API configuration (for development)
	Local LocalAPIConfig `json:"local"`

	// UseEdge determines whether to use Edge or Local API
	UseEdge bool `json:"useEdge"`
}

// EdgeAPIConfig contains Sitecore Edge API configuration
type EdgeAPIConfig struct {
	// ContextID is the Edge context ID
	ContextID string `json:"contextId"`

	// ClientContextID is the client-side Edge context ID
	ClientContextID string `json:"clientContextId"`

	// EdgeURL is the Edge API URL
	EdgeURL string `json:"edgeUrl"`
}

// LocalAPIConfig contains local Sitecore API configuration
type LocalAPIConfig struct {
	// APIKey is the Sitecore API key
	APIKey string `json:"apiKey"`

	// APIHost is the Sitecore API host
	APIHost string `json:"apiHost"`
}

// MultisiteConfig contains multisite configuration
type MultisiteConfig struct {
	// Enabled determines if multisite is enabled
	Enabled bool `json:"enabled"`

	// Sites is the list of available sites
	Sites []models.SiteInfo `json:"sites"`

	// DefaultSite is the default site
	DefaultSite models.SiteInfo `json:"defaultSite"`

	// UseCookieResolution enables cookie-based site resolution
	UseCookieResolution bool `json:"useCookieResolution"`
}

// PersonalizeConfig contains personalization configuration
type PersonalizeConfig struct {
	// Enabled determines if personalization is enabled
	Enabled bool `json:"enabled"`

	// Scope is the CDP scope
	Scope string `json:"scope"`

	// CDPEndpoint is the CDP API endpoint
	CDPEndpoint string `json:"cdpEndpoint"`
}

// EditingConfig contains editing/preview configuration
type EditingConfig struct {
	// Enabled determines if editing support is enabled
	Enabled bool `json:"enabled"`

	// Secret is the editing secret for security
	Secret string `json:"secret"`

	// InternalHostURL is the internal host URL for server-side requests
	InternalHostURL string `json:"internalHostUrl"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		API: APIConfig{
			Edge: EdgeAPIConfig{
				ContextID:       utils.GetEnvVar("SITECORE_EDGE_CONTEXT_ID"),
				ClientContextID: utils.GetEnvVar("SITECORE_EDGE_CLIENT_CONTEXT_ID"),
				EdgeURL:         utils.GetEnvVarOrDefault("SITECORE_EDGE_URL", "https://edge.sitecorecloud.io"),
			},
			Local: LocalAPIConfig{
				APIKey:  utils.GetEnvVar("SITECORE_API_KEY"),
				APIHost: utils.GetEnvVar("SITECORE_API_HOST"),
			},
			UseEdge: utils.GetEnvVarOrDefault("USE_EDGE_API", "false") == "true",
		},
		DefaultSite:     utils.GetEnvVarOrDefault("DEFAULT_SITE_NAME", "default"),
		DefaultLanguage: utils.GetEnvVarOrDefault("DEFAULT_LANGUAGE", "en"),
		Multisite: MultisiteConfig{
			Enabled:             utils.GetEnvVarOrDefault("MULTISITE_ENABLED", "true") == "true",
			UseCookieResolution: utils.GetEnvVarOrDefault("MULTISITE_USE_COOKIE", "true") == "true",
		},
		Personalize: PersonalizeConfig{
			Enabled:     utils.GetEnvVarOrDefault("PERSONALIZE_ENABLED", "false") == "true",
			Scope:       utils.GetEnvVar("PERSONALIZE_SCOPE"),
			CDPEndpoint: utils.GetEnvVarOrDefault("CDP_ENDPOINT", "https://api.boxever.com"),
		},
		Editing: EditingConfig{
			Enabled:         utils.GetEnvVarOrDefault("EDITING_ENABLED", "false") == "true",
			Secret:          utils.GetEnvVar("EDITING_SECRET"),
			InternalHostURL: utils.GetEnvVar("SITECORE_INTERNAL_EDITING_HOST_URL"),
		},
		EdgeTimeout: parseDuration(utils.GetEnvVarOrDefault("EDGE_TIMEOUT", "10s")),
		CDPTimeout:  parseDuration(utils.GetEnvVarOrDefault("CDP_TIMEOUT", "400ms")),
		EnableDebug: utils.GetEnvVar("DEBUG") != "",
	}

	return config
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate API configuration
	if c.API.UseEdge {
		if c.API.Edge.ContextID == "" {
			return fmt.Errorf("SITECORE_EDGE_CONTEXT_ID is required when using Edge API")
		}
		if c.API.Edge.EdgeURL == "" {
			return fmt.Errorf("SITECORE_EDGE_URL is required when using Edge API")
		}
	} else {
		if c.API.Local.APIKey == "" {
			return fmt.Errorf("SITECORE_API_KEY is required when using Local API")
		}
		if c.API.Local.APIHost == "" {
			return fmt.Errorf("SITECORE_API_HOST is required when using Local API")
		}
	}

	// Validate site configuration
	if c.DefaultSite == "" {
		return fmt.Errorf("DEFAULT_SITE_NAME is required")
	}

	// Validate personalization if enabled
	if c.Personalize.Enabled && c.Personalize.Scope == "" {
		return fmt.Errorf("PERSONALIZE_SCOPE is required when personalization is enabled")
	}

	return nil
}

// GetGraphQLEndpoint returns the appropriate GraphQL endpoint
func (c *Config) GetGraphQLEndpoint() string {
	if c.API.UseEdge {
		return c.API.Edge.EdgeURL + "/api/graphql/v1"
	}
	return c.API.Local.APIHost + "/sitecore/api/graph/edge"
}

// GetAPIKey returns the appropriate API key
func (c *Config) GetAPIKey() string {
	if c.API.UseEdge {
		return c.API.Edge.ContextID
	}
	return c.API.Local.APIKey
}

// parseDuration parses a duration string
func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 10 * time.Second // Default
	}
	return d
}


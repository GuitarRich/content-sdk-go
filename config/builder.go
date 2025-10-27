package config

import (
	"time"

	"github.com/content-sdk-go/models"
)

// ConfigBuilder provides a fluent API for building Config
type ConfigBuilder struct {
	config *Config
}

// NewConfigBuilder creates a new ConfigBuilder with defaults
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: &Config{
			DefaultLanguage: "en",
			EdgeTimeout:     10 * time.Second,
			CDPTimeout:      400 * time.Millisecond,
			API: APIConfig{
				Edge: EdgeAPIConfig{
					EdgeURL: "https://edge.sitecorecloud.io",
				},
			},
		},
	}
}

// WithEdgeAPI configures the Edge API
func (b *ConfigBuilder) WithEdgeAPI(contextID, clientContextID, edgeURL string) *ConfigBuilder {
	b.config.API.Edge.ContextID = contextID
	b.config.API.Edge.ClientContextID = clientContextID
	if edgeURL != "" {
		b.config.API.Edge.EdgeURL = edgeURL
	}
	b.config.API.UseEdge = true
	return b
}

// WithLocalAPI configures the Local API
func (b *ConfigBuilder) WithLocalAPI(apiKey, apiHost string) *ConfigBuilder {
	b.config.API.Local.APIKey = apiKey
	b.config.API.Local.APIHost = apiHost
	b.config.API.UseEdge = false
	return b
}

// WithDefaultSite sets the default site
func (b *ConfigBuilder) WithDefaultSite(siteName string) *ConfigBuilder {
	b.config.DefaultSite = siteName
	return b
}

// WithDefaultLanguage sets the default language
func (b *ConfigBuilder) WithDefaultLanguage(language string) *ConfigBuilder {
	b.config.DefaultLanguage = language
	return b
}

// WithMultisite configures multisite
func (b *ConfigBuilder) WithMultisite(enabled bool, sites []models.SiteInfo, useCookie bool) *ConfigBuilder {
	b.config.Multisite.Enabled = enabled
	b.config.Multisite.Sites = sites
	b.config.Multisite.UseCookieResolution = useCookie
	if len(sites) > 0 {
		b.config.Multisite.DefaultSite = sites[0]
	}
	return b
}

// WithPersonalization configures personalization
func (b *ConfigBuilder) WithPersonalization(enabled bool, scope, cdpEndpoint string) *ConfigBuilder {
	b.config.Personalize.Enabled = enabled
	b.config.Personalize.Scope = scope
	if cdpEndpoint != "" {
		b.config.Personalize.CDPEndpoint = cdpEndpoint
	} else {
		b.config.Personalize.CDPEndpoint = "https://api.boxever.com"
	}
	return b
}

// WithEditing configures editing support
func (b *ConfigBuilder) WithEditing(enabled bool, secret, internalHostURL string) *ConfigBuilder {
	b.config.Editing.Enabled = enabled
	b.config.Editing.Secret = secret
	b.config.Editing.InternalHostURL = internalHostURL
	return b
}

// WithTimeouts sets the API timeouts
func (b *ConfigBuilder) WithTimeouts(edgeTimeout, cdpTimeout time.Duration) *ConfigBuilder {
	b.config.EdgeTimeout = edgeTimeout
	b.config.CDPTimeout = cdpTimeout
	return b
}

// WithDebug enables debug logging
func (b *ConfigBuilder) WithDebug(enabled bool) *ConfigBuilder {
	b.config.EnableDebug = enabled
	return b
}

// Build builds and validates the configuration
func (b *ConfigBuilder) Build() (*Config, error) {
	if err := b.config.Validate(); err != nil {
		return nil, err
	}
	return b.config, nil
}

// BuildOrPanic builds the configuration and panics on error
func (b *ConfigBuilder) BuildOrPanic() *Config {
	config, err := b.Build()
	if err != nil {
		panic(err)
	}
	return config
}


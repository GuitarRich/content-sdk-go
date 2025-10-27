package config

import (
	"time"

	"github.com/content-sdk-go/models"
)

// ExampleEdgeConfig returns an example configuration using Edge API
func ExampleEdgeConfig() (*Config, error) {
	return NewConfigBuilder().
		WithEdgeAPI("your-context-id", "your-client-context-id", "").
		WithDefaultSite("mysite").
		WithDefaultLanguage("en").
		WithMultisite(false, nil, false).
		Build()
}

// ExampleLocalConfig returns an example configuration using local API
func ExampleLocalConfig() (*Config, error) {
	return NewConfigBuilder().
		WithLocalAPI("your-api-key", "https://cm.localhost").
		WithDefaultSite("mysite").
		WithDefaultLanguage("en").
		Build()
}

// ExampleMultisiteConfig returns an example multisite configuration
func ExampleMultisiteConfig() (*Config, error) {
	sites := []models.SiteInfo{
		{
			Name:     "site1",
			Language: "en",
			HostName: "www.site1.com",
		},
		{
			Name:     "site2",
			Language: "fr",
			HostName: "www.site2.com",
		},
	}

	return NewConfigBuilder().
		WithEdgeAPI("your-context-id", "", "").
		WithDefaultSite("site1").
		WithMultisite(true, sites, true).
		Build()
}

// ExamplePersonalizationConfig returns an example configuration with personalization
func ExamplePersonalizationConfig() (*Config, error) {
	return NewConfigBuilder().
		WithEdgeAPI("your-context-id", "", "").
		WithDefaultSite("mysite").
		WithPersonalization(true, "your-personalize-scope", "").
		WithTimeouts(10*time.Second, 400*time.Millisecond).
		Build()
}

// ExampleEditingConfig returns an example configuration with editing support
func ExampleEditingConfig() (*Config, error) {
	return NewConfigBuilder().
		WithEdgeAPI("your-context-id", "", "").
		WithDefaultSite("mysite").
		WithEditing(true, "your-editing-secret", "http://localhost:3000").
		Build()
}

// ExampleFullConfig returns a fully-configured example
func ExampleFullConfig() (*Config, error) {
	sites := []models.SiteInfo{
		{
			Name:     "site1",
			Language: "en",
			HostName: "www.site1.com",
		},
		{
			Name:     "site2",
			Language: "fr",
			HostName: "www.site2.com",
		},
	}

	return NewConfigBuilder().
		WithEdgeAPI("your-context-id", "your-client-context-id", "").
		WithDefaultSite("site1").
		WithDefaultLanguage("en").
		WithMultisite(true, sites, true).
		WithPersonalization(true, "your-personalize-scope", "").
		WithEditing(true, "your-editing-secret", "http://localhost:3000").
		WithTimeouts(10*time.Second, 400*time.Millisecond).
		WithDebug(true).
		Build()
}


package config

import (
	"os"
	"testing"
	"time"

	"github.com/content-sdk-go/models"
)

func TestLoadConfig(t *testing.T) {
	// Set up test environment variables
	os.Setenv("SITECORE_EDGE_CONTEXT_ID", "test-context")
	os.Setenv("USE_EDGE_API", "true")
	os.Setenv("DEFAULT_SITE_NAME", "testsite")
	os.Setenv("DEFAULT_LANGUAGE", "fr")
	defer func() {
		os.Unsetenv("SITECORE_EDGE_CONTEXT_ID")
		os.Unsetenv("USE_EDGE_API")
		os.Unsetenv("DEFAULT_SITE_NAME")
		os.Unsetenv("DEFAULT_LANGUAGE")
	}()

	config := LoadConfig()

	if config.API.Edge.ContextID != "test-context" {
		t.Errorf("expected context ID 'test-context', got '%s'", config.API.Edge.ContextID)
	}

	if !config.API.UseEdge {
		t.Error("expected UseEdge to be true")
	}

	if config.DefaultSite != "testsite" {
		t.Errorf("expected default site 'testsite', got '%s'", config.DefaultSite)
	}

	if config.DefaultLanguage != "fr" {
		t.Errorf("expected default language 'fr', got '%s'", config.DefaultLanguage)
	}
}

func TestConfigValidate_EdgeAPI(t *testing.T) {
	config := &Config{
		API: APIConfig{
			UseEdge: true,
			Edge: EdgeAPIConfig{
				ContextID: "test-context",
				EdgeURL:   "https://edge.sitecorecloud.io",
			},
		},
		DefaultSite: "testsite",
	}

	if err := config.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestConfigValidate_MissingEdgeContext(t *testing.T) {
	config := &Config{
		API: APIConfig{
			UseEdge: true,
			Edge: EdgeAPIConfig{
				EdgeURL: "https://edge.sitecorecloud.io",
			},
		},
		DefaultSite: "testsite",
	}

	if err := config.Validate(); err == nil {
		t.Error("expected validation error for missing context ID")
	}
}

func TestConfigValidate_LocalAPI(t *testing.T) {
	config := &Config{
		API: APIConfig{
			UseEdge: false,
			Local: LocalAPIConfig{
				APIKey:  "test-key",
				APIHost: "https://cm.localhost",
			},
		},
		DefaultSite: "testsite",
	}

	if err := config.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestConfigBuilder_EdgeAPI(t *testing.T) {
	config, err := NewConfigBuilder().
		WithEdgeAPI("test-context", "client-context", "").
		WithDefaultSite("mysite").
		Build()

	if err != nil {
		t.Errorf("unexpected build error: %v", err)
	}

	if config.API.Edge.ContextID != "test-context" {
		t.Errorf("expected context ID 'test-context', got '%s'", config.API.Edge.ContextID)
	}

	if config.DefaultSite != "mysite" {
		t.Errorf("expected default site 'mysite', got '%s'", config.DefaultSite)
	}

	if !config.API.UseEdge {
		t.Error("expected UseEdge to be true")
	}
}

func TestConfigBuilder_LocalAPI(t *testing.T) {
	config, err := NewConfigBuilder().
		WithLocalAPI("test-key", "https://cm.localhost").
		WithDefaultSite("mysite").
		Build()

	if err != nil {
		t.Errorf("unexpected build error: %v", err)
	}

	if config.API.Local.APIKey != "test-key" {
		t.Errorf("expected API key 'test-key', got '%s'", config.API.Local.APIKey)
	}

	if config.API.UseEdge {
		t.Error("expected UseEdge to be false")
	}
}

func TestConfigBuilder_Multisite(t *testing.T) {
	sites := []models.SiteInfo{
		{Name: "site1", Language: "en", HostName: "site1.com"},
		{Name: "site2", Language: "fr", HostName: "site2.com"},
	}

	config, err := NewConfigBuilder().
		WithEdgeAPI("test-context", "", "").
		WithDefaultSite("site1").
		WithMultisite(true, sites, true).
		Build()

	if err != nil {
		t.Errorf("unexpected build error: %v", err)
	}

	if !config.Multisite.Enabled {
		t.Error("expected multisite to be enabled")
	}

	if len(config.Multisite.Sites) != 2 {
		t.Errorf("expected 2 sites, got %d", len(config.Multisite.Sites))
	}

	if config.Multisite.DefaultSite.Name != "site1" {
		t.Errorf("expected default site 'site1', got '%s'", config.Multisite.DefaultSite.Name)
	}
}

func TestConfigBuilder_Personalization(t *testing.T) {
	config, err := NewConfigBuilder().
		WithEdgeAPI("test-context", "", "").
		WithDefaultSite("mysite").
		WithPersonalization(true, "test-scope", "").
		Build()

	if err != nil {
		t.Errorf("unexpected build error: %v", err)
	}

	if !config.Personalize.Enabled {
		t.Error("expected personalization to be enabled")
	}

	if config.Personalize.Scope != "test-scope" {
		t.Errorf("expected scope 'test-scope', got '%s'", config.Personalize.Scope)
	}

	if config.Personalize.CDPEndpoint != "https://api.boxever.com" {
		t.Errorf("expected default CDP endpoint, got '%s'", config.Personalize.CDPEndpoint)
	}
}

func TestConfigBuilder_Timeouts(t *testing.T) {
	config, err := NewConfigBuilder().
		WithEdgeAPI("test-context", "", "").
		WithDefaultSite("mysite").
		WithTimeouts(5*time.Second, 200*time.Millisecond).
		Build()

	if err != nil {
		t.Errorf("unexpected build error: %v", err)
	}

	if config.EdgeTimeout != 5*time.Second {
		t.Errorf("expected edge timeout 5s, got %v", config.EdgeTimeout)
	}

	if config.CDPTimeout != 200*time.Millisecond {
		t.Errorf("expected CDP timeout 200ms, got %v", config.CDPTimeout)
	}
}

func TestConfigGetGraphQLEndpoint_Edge(t *testing.T) {
	config := &Config{
		API: APIConfig{
			UseEdge: true,
			Edge: EdgeAPIConfig{
				EdgeURL: "https://edge.sitecorecloud.io",
			},
		},
	}

	endpoint := config.GetGraphQLEndpoint()
	expected := "https://edge.sitecorecloud.io/api/graphql/v1"

	if endpoint != expected {
		t.Errorf("expected endpoint '%s', got '%s'", expected, endpoint)
	}
}

func TestConfigGetGraphQLEndpoint_Local(t *testing.T) {
	config := &Config{
		API: APIConfig{
			UseEdge: false,
			Local: LocalAPIConfig{
				APIHost: "https://cm.localhost",
			},
		},
	}

	endpoint := config.GetGraphQLEndpoint()
	expected := "https://cm.localhost/sitecore/api/graph/edge"

	if endpoint != expected {
		t.Errorf("expected endpoint '%s', got '%s'", expected, endpoint)
	}
}

func TestConfigGetAPIKey_Edge(t *testing.T) {
	config := &Config{
		API: APIConfig{
			UseEdge: true,
			Edge: EdgeAPIConfig{
				ContextID: "edge-context",
			},
		},
	}

	apiKey := config.GetAPIKey()

	if apiKey != "edge-context" {
		t.Errorf("expected API key 'edge-context', got '%s'", apiKey)
	}
}

func TestConfigGetAPIKey_Local(t *testing.T) {
	config := &Config{
		API: APIConfig{
			UseEdge: false,
			Local: LocalAPIConfig{
				APIKey: "local-key",
			},
		},
	}

	apiKey := config.GetAPIKey()

	if apiKey != "local-key" {
		t.Errorf("expected API key 'local-key', got '%s'", apiKey)
	}
}


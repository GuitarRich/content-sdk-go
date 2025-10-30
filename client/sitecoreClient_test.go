package client

import (
	"testing"

	"github.com/content-sdk-go/models"
)

func TestParsePath_String(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "string with leading slash",
			input:    "/home",
			expected: "/home",
		},
		{
			name:     "string without leading slash",
			input:    "home",
			expected: "/home",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "/",
		},
		{
			name:     "root path",
			input:    "/",
			expected: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePath(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestParsePath_Slice(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "multiple segments",
			input:    []string{"home", "about"},
			expected: "/home/about",
		},
		{
			name:     "segments with slashes",
			input:    []string{"/home/", "/about/"},
			expected: "/home/about",
		},
		{
			name:     "empty segments",
			input:    []string{"", "home", "", "about", ""},
			expected: "/home/about",
		},
		{
			name:     "root segment",
			input:    []string{"/"},
			expected: "/",
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePath(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetSiteRewrite(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		siteName string
		expected string
	}{
		{
			name:     "simple path",
			path:     "/home",
			siteName: "mysite",
			expected: "/_site_mysite/home",
		},
		{
			name:     "path without leading slash",
			path:     "home",
			siteName: "mysite",
			expected: "/_site_mysite/home",
		},
		{
			name:     "nested path",
			path:     "/about/team",
			siteName: "site1",
			expected: "/_site_site1/about/team",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetSiteRewrite(tt.path, tt.siteName)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetSiteRewriteData(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		defaultSite string
		expected    models.SiteRewriteData
	}{
		{
			name:        "path with site prefix",
			path:        "/_site_mysite/home/about",
			defaultSite: "default",
			expected: models.SiteRewriteData{
				SiteName:       "mysite",
				NormalizedPath: "/home/about",
			},
		},
		{
			name:        "path without site prefix",
			path:        "/home/about",
			defaultSite: "default",
			expected: models.SiteRewriteData{
				SiteName:       "default",
				NormalizedPath: "/home/about",
			},
		},
		{
			name:        "root path with site",
			path:        "/_site_site1/",
			defaultSite: "default",
			expected: models.SiteRewriteData{
				SiteName:       "site1",
				NormalizedPath: "/",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetSiteRewriteData(tt.path, tt.defaultSite)
			if result.SiteName != tt.expected.SiteName {
				t.Errorf("expected siteName %s, got %s", tt.expected.SiteName, result.SiteName)
			}
			if result.NormalizedPath != tt.expected.NormalizedPath {
				t.Errorf("expected normalizedPath %s, got %s", tt.expected.NormalizedPath, result.NormalizedPath)
			}
		})
	}
}

func TestNormalizeSiteRewrite(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "path with site prefix",
			path:     "/_site_mysite/home",
			expected: "/home",
		},
		{
			name:     "path without site prefix",
			path:     "/home",
			expected: "/home",
		},
		{
			name:     "nested path with site",
			path:     "/_site_site1/about/team",
			expected: "/about/team",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeSiteRewrite(tt.path)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetPersonalizedRewrite(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		variantId string
		expected  string
	}{
		{
			name:      "simple path",
			path:      "/home",
			variantId: "abc123",
			expected:  "/_variantId_abc123/home",
		},
		{
			name:      "nested path",
			path:      "/products/item",
			variantId: "xyz789",
			expected:  "/_variantId_xyz789/products/item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPersonalizedRewrite(tt.path, tt.variantId)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetPersonalizedRewriteData(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected models.PersonalizeRewriteData
	}{
		{
			name: "path with variant",
			path: "/_variantId_abc123/home",
			expected: models.PersonalizeRewriteData{
				VariantId:      "abc123",
				NormalizedPath: "/home",
			},
		},
		{
			name: "path without variant",
			path: "/home",
			expected: models.PersonalizeRewriteData{
				VariantId:      "",
				NormalizedPath: "/home",
			},
		},
		{
			name: "nested path with variant",
			path: "/_variantId_xyz789/products/item",
			expected: models.PersonalizeRewriteData{
				VariantId:      "xyz789",
				NormalizedPath: "/products/item",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPersonalizedRewriteData(tt.path)
			if result.VariantId != tt.expected.VariantId {
				t.Errorf("expected variantId %s, got %s", tt.expected.VariantId, result.VariantId)
			}
			if result.NormalizedPath != tt.expected.NormalizedPath {
				t.Errorf("expected normalizedPath %s, got %s", tt.expected.NormalizedPath, result.NormalizedPath)
			}
		})
	}
}

func TestNormalizePersonalizedRewrite(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "path with variant",
			path:     "/_variantId_abc123/home",
			expected: "/home",
		},
		{
			name:     "path without variant",
			path:     "/home",
			expected: "/home",
		},
		{
			name:     "nested path with variant",
			path:     "/_variantId_xyz789/products/item",
			expected: "/products/item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizePersonalizedRewrite(tt.path)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestSitecoreClient_ParsePath(t *testing.T) {
	client := &SitecoreClient{
		defaultSite: "mysite",
		defaultLang: "en",
	}

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "path with both site and variant",
			path:     "/_site_site1/_variantId_abc123/home",
			expected: "/home",
		},
		{
			name:     "path with only site",
			path:     "/_site_site1/home",
			expected: "/home",
		},
		{
			name:     "path with only variant",
			path:     "/_variantId_abc123/home",
			expected: "/home",
		},
		{
			name:     "plain path",
			path:     "/home",
			expected: "/home",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.ParsePath(tt.path)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

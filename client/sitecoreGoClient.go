package client

import (
	"net/http"
	"strings"
)

const (
	// SITE_PREFIX is the prefix used to identify site names in paths
	SITE_PREFIX = "/"
)

// SiteRewriteData contains site information extracted from a path
type SiteRewriteData struct {
	SiteName string
}

type SitecoreGoClient struct {
	Client *http.Client
}

func NewSitecoreGoClient(httpClient *http.Client) *SitecoreGoClient {
	return &SitecoreGoClient{Client: httpClient}
}

func (c *SitecoreGoClient) GetPage(path string, lang string) (*http.Response, error) {
	normalizedPath := parsePath(path)

	return c.Client.Get(path)
}

func (c *SitecoreGoClient) GetPreview(path string) (*http.Response, error) {
	return c.Client.Get(path)
}

// parsePath normalizes path regardless of type
// Accepts either string or []string and returns a normalized string path
func parsePath(path interface{}) string {
	switch p := path.(type) {
	case string:
		// If string starts with '/', return as-is, otherwise prepend '/'
		if strings.HasPrefix(p, "/") {
			return p
		}
		return "/" + p
	case []string:
		// Filter out '/' parts and trim slashes from each part
		var parts []string
		for _, part := range p {
			if part == "/" {
				continue
			}
			// Remove leading and trailing slashes
			trimmed := strings.Trim(part, "/")
			if trimmed != "" {
				parts = append(parts, trimmed)
			}
		}
		// Join with '/' and prepend '/'
		return "/" + strings.Join(parts, "/")
	default:
		return ""
	}
}

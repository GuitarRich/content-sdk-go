package models

// SiteInfo contains information about a Sitecore site
type SiteInfo struct {
	// Name is the site name in Sitecore
	Name string `json:"name"`

	// HostName is the primary hostname for the site
	HostName string `json:"hostName"`

	// Language is the default language for the site
	Language string `json:"language"`

	// RootPath is the content root path in Sitecore
	RootPath string `json:"rootPath,omitempty"`

	// Database is the Sitecore database (master, web)
	Database string `json:"database,omitempty"`
}

// SiteRewriteData contains site information extracted from a path
type SiteRewriteData struct {
	// SiteName is the name of the site
	SiteName string `json:"siteName"`

	// NormalizedPath is the path with site prefix removed
	NormalizedPath string `json:"normalizedPath"`
}

// RedirectInfo contains information about a URL redirect
type RedirectInfo struct {
	// Pattern is the URL pattern to match (can be regex)
	Pattern string `json:"pattern"`

	// Target is the redirect destination URL
	Target string `json:"target"`

	// RedirectType is the type of redirect (301, 302, SERVER_TRANSFER)
	RedirectType RedirectType `json:"redirectType"`

	// Locale is the language this redirect applies to (optional)
	Locale string `json:"locale,omitempty"`

	// IsRegex indicates if the pattern is a regular expression
	IsRegex bool `json:"isRegex,omitempty"`
}

// RedirectType represents the type of HTTP redirect
type RedirectType string

const (
	// Redirect301 is a permanent redirect
	Redirect301 RedirectType = "301"

	// Redirect302 is a temporary redirect
	Redirect302 RedirectType = "302"

	// RedirectServerTransfer is a server-side transfer (no HTTP redirect)
	RedirectServerTransfer RedirectType = "SERVER_TRANSFER"
)

// SitemapEntry represents a single entry in a sitemap
type SitemapEntry struct {
	// Loc is the URL of the page
	Loc string `json:"loc"`

	// LastMod is the last modification date (ISO 8601 format)
	LastMod string `json:"lastmod,omitempty"`

	// ChangeFreq is how frequently the page changes
	ChangeFreq string `json:"changefreq,omitempty"`

	// Priority is the priority of this URL (0.0 to 1.0)
	Priority string `json:"priority,omitempty"`
}

// RobotsDirective represents robots.txt directives from Sitecore
type RobotsDirective struct {
	// Content is the robots.txt content from Sitecore
	Content string `json:"content"`

	// UserAgent is the user agent this applies to (optional)
	UserAgent string `json:"userAgent,omitempty"`

	// Allow are the allowed paths
	Allow []string `json:"allow,omitempty"`

	// Disallow are the disallowed paths
	Disallow []string `json:"disallow,omitempty"`

	// Sitemap URLs
	Sitemap []string `json:"sitemap,omitempty"`
}

package site

import (
	"fmt"
	"strings"

	"github.com/guitarrich/content-sdk-go/models"
)

// SiteResolver resolves sites by hostname or name
type SiteResolver interface {
	GetByHost(hostname string) (*models.SiteInfo, error)
	GetByName(name string) (*models.SiteInfo, error)
}

// siteResolverImpl is the default implementation
type siteResolverImpl struct {
	sites       []models.SiteInfo
	defaultSite models.SiteInfo
}

// NewSiteResolver creates a new site resolver
func NewSiteResolver(sites []models.SiteInfo, defaultSite models.SiteInfo) SiteResolver {
	return &siteResolverImpl{
		sites:       sites,
		defaultSite: defaultSite,
	}
}

// GetByHost resolves a site by hostname
func (r *siteResolverImpl) GetByHost(hostname string) (*models.SiteInfo, error) {
	// Normalize hostname (remove port if present)
	host := hostname
	if idx := strings.Index(hostname, ":"); idx > 0 {
		host = hostname[:idx]
	}

	// Convert to lowercase for case-insensitive matching
	host = strings.ToLower(host)

	// Try exact match first
	for _, site := range r.sites {
		if strings.ToLower(site.HostName) == host {
			return &site, nil
		}
	}

	// Try wildcard matching (e.g., *.example.com)
	for _, site := range r.sites {
		if matchesWildcard(site.HostName, host) {
			return &site, nil
		}
	}

	// Return default site if no match found
	return &r.defaultSite, nil
}

// GetByName resolves a site by name
func (r *siteResolverImpl) GetByName(name string) (*models.SiteInfo, error) {
	// Convert to lowercase for case-insensitive matching
	searchName := strings.ToLower(name)

	for _, site := range r.sites {
		if strings.ToLower(site.Name) == searchName {
			return &site, nil
		}
	}

	// Return default site if not found
	if strings.ToLower(r.defaultSite.Name) == searchName {
		return &r.defaultSite, nil
	}

	return nil, fmt.Errorf("site not found: %s", name)
}

// matchesWildcard checks if a hostname matches a wildcard pattern
// Supports patterns like *.example.com
func matchesWildcard(pattern, hostname string) bool {
	pattern = strings.ToLower(pattern)
	hostname = strings.ToLower(hostname)

	if !strings.Contains(pattern, "*") {
		return false
	}

	// Simple wildcard matching - replace * with any characters
	if strings.HasPrefix(pattern, "*.") {
		suffix := pattern[2:] // Remove "*."
		return strings.HasSuffix(hostname, "."+suffix) || hostname == suffix
	}

	return false
}

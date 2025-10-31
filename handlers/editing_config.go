package handlers

import (
	"net/http"
	"slices"

	"github.com/content-sdk-go/debug"
	"github.com/content-sdk-go/middleware"
)

// ComponentRegistry interface for accessing registered components
type ComponentRegistry interface {
	List() []string
}

// EditingConfigHandler provides configuration for the Sitecore Pages editor
type EditingConfigHandler struct {
	registry ComponentRegistry
}

// EditingConfigResponse contains the editing configuration
type EditingConfigResponse struct {
	// Components is the list of registered component names
	Components []string `json:"components"`

	// Packages contains the package versions
	Packages map[string]string `json:"packages"`

	// EditMode is the editing mode
	EditMode string `json:"editMode"`
}

// NewEditingConfigHandler creates a new editing config handler
func NewEditingConfigHandler(registry ComponentRegistry) *EditingConfigHandler {
	return &EditingConfigHandler{
		registry: registry,
	}
}

// Handle processes editing config requests
func (h *EditingConfigHandler) Handle(ctx middleware.Context) error {
	debug.Editing("handling editing config request")

	// Get registered components, excluding "Unknown"
	allComponents := h.registry.List()
	components := make([]string, 0, len(allComponents))
	for _, name := range allComponents {
		if name != "Unknown" {
			components = append(components, name)
		}
	}

	// Sort components alphabetically
	slices.Sort(components)

	// Build response with hardcoded package versions
	response := EditingConfigResponse{
		Components: components,
		Packages: map[string]string{
			"@sitecore-cloudsdk/core":        "0.5.4",
			"@sitecore-cloudsdk/events":      "0.5.4",
			"@sitecore-cloudsdk/personalize": "0.5.4",
			"@sitecore-cloudsdk/utils":       "0.5.4",
			"@sitecore-content-sdk/cli":      "1.1.0",
			"@sitecore-content-sdk/core":     "1.1.0",
			"@sitecore-content-sdk/nextjs":   "1.1.0",
			"@sitecore-content-sdk/react":    "1.1.0",
			"@sitecore-feaas/clientside":     "0.6.2",
			"@sitecore/byoc":                 "0.3.0",
			"@sitecore/components":           "2.1.0",
		},
		EditMode: "metadata",
	}

	// Return configuration as JSON
	return ctx.JSON(http.StatusOK, response)
}

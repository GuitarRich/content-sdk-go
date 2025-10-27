package handlers

import (
	"net/http"

	"github.com/content-sdk-go/debug"
	"github.com/content-sdk-go/middleware"
)

// EditingConfigHandler provides configuration for the Sitecore Pages editor
type EditingConfigHandler struct {
	config EditingConfigResponse
}

// EditingConfigResponse contains the editing configuration
type EditingConfigResponse struct {
	// SitecoreEdgeURL is the Sitecore Edge URL
	SitecoreEdgeURL string `json:"sitecoreEdgeUrl,omitempty"`

	// SitecoreEdgeContextID is the Edge context ID
	SitecoreEdgeContextID string `json:"sitecoreEdgeContextId,omitempty"`

	// DefaultLanguage is the default language
	DefaultLanguage string `json:"defaultLanguage,omitempty"`

	// DefaultSite is the default site name
	DefaultSite string `json:"defaultSite,omitempty"`
}

// EditingConfigHandlerConfig contains configuration for the editing config handler
type EditingConfigHandlerConfig struct {
	SitecoreEdgeURL       string
	SitecoreEdgeContextID string
	DefaultLanguage       string
	DefaultSite           string
}

// NewEditingConfigHandler creates a new editing config handler
func NewEditingConfigHandler(config EditingConfigHandlerConfig) *EditingConfigHandler {
	return &EditingConfigHandler{
		config: EditingConfigResponse{
			SitecoreEdgeURL:       config.SitecoreEdgeURL,
			SitecoreEdgeContextID: config.SitecoreEdgeContextID,
			DefaultLanguage:       config.DefaultLanguage,
			DefaultSite:           config.DefaultSite,
		},
	}
}

// Handle processes editing config requests
func (h *EditingConfigHandler) Handle(ctx middleware.Context) error {
	debug.Editing("handling editing config request")

	// Return configuration as JSON
	return ctx.JSON(http.StatusOK, h.config)
}

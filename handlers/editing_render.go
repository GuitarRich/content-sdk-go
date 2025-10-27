package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/content-sdk-go/client"
	"github.com/content-sdk-go/debug"
	"github.com/content-sdk-go/middleware"
	"github.com/content-sdk-go/models"
)

// EditingRenderHandler handles rendering for the Sitecore Pages editor
type EditingRenderHandler struct {
	client *client.SitecoreClient
}

// NewEditingRenderHandler creates a new editing render handler
func NewEditingRenderHandler(sitecoreClient *client.SitecoreClient) *EditingRenderHandler {
	return &EditingRenderHandler{
		client: sitecoreClient,
	}
}

// Handle processes editing render requests
func (h *EditingRenderHandler) Handle(ctx middleware.Context) error {
	debug.Editing("handling editing render request")

	// Read request body
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		debug.Editing("error reading request body: %v", err)
		return ctx.String(http.StatusBadRequest, "Invalid request body")
	}

	// Parse preview data
	var previewData models.PreviewData
	if err := json.Unmarshal(body, &previewData); err != nil {
		debug.Editing("error parsing preview data: %v", err)
		return ctx.String(http.StatusBadRequest, "Invalid preview data")
	}

	debug.Editing("preview data: itemId=%s, language=%s, site=%s, mode=%s",
		previewData.ItemID, previewData.Language, previewData.Site, previewData.Mode)

	// Fetch preview page
	page, err := h.client.GetPreview(previewData)
	if err != nil {
		debug.Editing("error fetching preview: %v", err)
		return ctx.String(http.StatusInternalServerError, "Error fetching preview")
	}

	// Return page as JSON
	return ctx.JSON(http.StatusOK, page)
}

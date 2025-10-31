package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/content-sdk-go/client"
	"github.com/content-sdk-go/debug"
	"github.com/content-sdk-go/middleware"
	"github.com/content-sdk-go/models"
)

// PageRenderer is an interface for rendering pages
// This allows the handler to work with any renderer implementation
type PageRenderer interface {
	RenderPage(ctx context.Context, page *models.Page) (templ.Component, error)
}

// EditingRenderHandler handles rendering for the Sitecore Pages editor
type EditingRenderHandler struct {
	client   *client.SitecoreClient
	renderer PageRenderer
}

// NewEditingRenderHandler creates a new editing render handler
func NewEditingRenderHandler(sitecoreClient *client.SitecoreClient, renderer PageRenderer) *EditingRenderHandler {
	return &EditingRenderHandler{
		client:   sitecoreClient,
		renderer: renderer,
	}
}

// Handle processes editing render requests from Sitecore Pages
// Expects query parameters: sc_itemid, sc_lang, sc_site, sc_layoutKind, mode, route, secret
func (h *EditingRenderHandler) Handle(ctx middleware.Context) error {
	debug.Editing("handling editing render request")

	// Extract query parameters
	itemID := ctx.Request().URL.Query().Get("sc_itemid")
	language := ctx.Request().URL.Query().Get("sc_lang")
	site := ctx.Request().URL.Query().Get("sc_site")
	layoutKind := ctx.Request().URL.Query().Get("sc_layoutKind")
	mode := ctx.Request().URL.Query().Get("mode")
	route := ctx.Request().URL.Query().Get("route")
	version := ctx.Request().URL.Query().Get("sc_version")

	debug.Editing("query params: itemId=%s, lang=%s, site=%s, layoutKind=%s, mode=%s, route=%s",
		itemID, language, site, layoutKind, mode, route)

	// Validate required parameters
	if itemID == "" {
		debug.Editing("missing required parameter: sc_itemid")
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing required parameter: sc_itemid",
		})
	}

	if language == "" {
		debug.Editing("missing required parameter: sc_lang")
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing required parameter: sc_lang",
		})
	}

	if site == "" {
		debug.Editing("missing required parameter: sc_site")
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing required parameter: sc_site",
		})
	}

	// Build preview data from query parameters
	previewData := models.PreviewData{
		ItemID:   itemID,
		Language: language,
		Site:     site,
		Version:  version,
		Route:    route,
	}

	// Set layout kind (default to final)
	if layoutKind != "" && strings.ToLower(layoutKind) == "shared" {
		previewData.LayoutKind = models.LayoutKindShared
	} else {
		previewData.LayoutKind = models.LayoutKindFinal
	}

	// Set preview mode
	switch strings.ToLower(mode) {
	case "edit":
		previewData.Mode = models.PreviewModeEdit
	case "preview":
		previewData.Mode = models.PreviewModePreview
	case "metadata":
		previewData.Mode = models.PreviewModeMetadata
	default:
		previewData.Mode = models.PreviewModeEdit // Default to edit mode
	}

	debug.Editing("fetching preview: itemId=%s, language=%s, site=%s, mode=%s, layoutKind=%s",
		previewData.ItemID, previewData.Language, previewData.Site, previewData.Mode, previewData.LayoutKind)

	// Fetch preview page
	page, err := h.client.GetPreview(previewData)
	if err != nil {
		debug.Editing("error fetching preview: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error fetching preview: " + err.Error(),
		})
	}

	debug.Editing("preview page fetched, rendering HTML")

	// If no renderer is configured, return JSON (backwards compatibility)
	if h.renderer == nil {
		debug.Editing("no renderer configured, returning JSON")
		return ctx.JSON(http.StatusOK, page)
	}

	// Render the page to HTML using the renderer
	component, err := h.renderer.RenderPage(ctx.Request().Context(), page)
	if err != nil {
		debug.Editing("error rendering page: %v", err)
		return ctx.String(http.StatusInternalServerError, "Error rendering page: "+err.Error())
	}

	// Set content type for HTML
	ctx.SetHeader("Content-Type", "text/html; charset=utf-8")
	ctx.SetHeader("X-Editing-Mode", string(previewData.Mode))

	// Render component to response writer
	// Note: The status code will default to 200 OK when we start writing
	err = component.Render(ctx.Request().Context(), ctx.Response())
	if err != nil {
		debug.Editing("error writing response: %v", err)
		return err
	}

	debug.Editing("page rendered successfully")
	return nil
}

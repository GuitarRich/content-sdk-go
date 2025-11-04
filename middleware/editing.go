package middleware

import (
	"context"

	"github.com/guitarrich/content-sdk-go/models"
	"github.com/labstack/echo/v4"
)

// EditingModeMiddleware detects and stores editing mode information
// It checks for Sitecore query parameters (sc_mode, sc_lang, sc_itemid)
// and stores them in the request context for use by handlers and renderers
func EditingModeMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Detect editing mode from query parameters
			scMode := c.QueryParam("sc_mode")
			scLang := c.QueryParam("sc_lang")
			scItemID := c.QueryParam("sc_itemid")

			// Determine if we're in editing mode
			isEditingMode := scMode == "edit" || scMode == "preview"
			isPreview := scMode == "preview"
			isEdit := scMode == "edit"

			// Determine the page mode
			var mode models.PageMode
			switch scMode {
			case "edit":
				mode = models.PageModeEdit
			case "preview":
				mode = models.PageModePreview
			default:
				mode = models.PageModeNormal
			}

			// Store in request context
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, "sc_mode", scMode)
			ctx = context.WithValue(ctx, "sc_lang", scLang)
			ctx = context.WithValue(ctx, "sc_itemid", scItemID)
			ctx = context.WithValue(ctx, "isEditingMode", isEditingMode)
			ctx = context.WithValue(ctx, "isPreview", isPreview)
			ctx = context.WithValue(ctx, "isEdit", isEdit)
			ctx = context.WithValue(ctx, "pageMode", mode)

			// Create EditingContext for easy access
			editingContext := &models.EditingContext{
				IsEditing: isEdit,
				IsPreview: isPreview,
				Mode:      mode,
				QueryParams: map[string]string{
					"sc_mode":   scMode,
					"sc_lang":   scLang,
					"sc_itemid": scItemID,
				},
			}
			ctx = context.WithValue(ctx, "editingContext", editingContext)

			// Update request with new context
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

// GetEditingContext retrieves the EditingContext from the request context
func GetEditingContext(ctx context.Context) *models.EditingContext {
	if editingCtx, ok := ctx.Value("editingContext").(*models.EditingContext); ok {
		return editingCtx
	}
	return &models.EditingContext{
		IsEditing: false,
		IsPreview: false,
		Mode:      models.PageModeNormal,
	}
}

// IsEditingMode is a helper to check if we're in any editing mode
func IsEditingMode(ctx context.Context) bool {
	if isEditing, ok := ctx.Value("isEditingMode").(bool); ok {
		return isEditing
	}
	return false
}

// IsPreviewMode is a helper to check if we're in preview mode
func IsPreviewMode(ctx context.Context) bool {
	if isPreview, ok := ctx.Value("isPreview").(bool); ok {
		return isPreview
	}
	return false
}

// IsEditMode is a helper to check if we're in edit mode
func IsEditMode(ctx context.Context) bool {
	if isEdit, ok := ctx.Value("isEdit").(bool); ok {
		return isEdit
	}
	return false
}

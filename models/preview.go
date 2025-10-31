package models

// LayoutKind represents the type of layout to fetch
type LayoutKind string

const (
	// LayoutKindFinal is the final layout
	LayoutKindFinal LayoutKind = "final"

	// LayoutKindShared is the shared layout
	LayoutKindShared LayoutKind = "shared"
)

// PreviewData contains data for preview/editing modes
type PreviewData struct {
	// ItemID is the Sitecore item ID being previewed
	ItemID string `json:"itemId"`

	// Language is the language of the item
	Language string `json:"language"`

	// Site is the site name
	Site string `json:"site"`

	// Version is the item version number
	Version string `json:"version"`

	// Mode is the preview mode (preview or edit)
	Mode PreviewMode `json:"mode"`

	// LayoutKind is the layout variant to fetch (final or shared)
	LayoutKind LayoutKind `json:"layoutKind,omitempty"`

	// VariantIds contains personalization variant IDs
	VariantIds []string `json:"variantIds,omitempty"`

	// Route is the page route being previewed
	Route string `json:"route,omitempty"`

	// ServerURL is the Sitecore server URL for preview requests
	ServerURL string `json:"serverUrl,omitempty"`
}

// PreviewMode represents the type of preview/editing mode
type PreviewMode string

const (
	// PreviewModePreview is standard preview mode
	PreviewModePreview PreviewMode = "preview"

	// PreviewModeEdit is Pages editor mode
	PreviewModeEdit PreviewMode = "edit"

	// PreviewModeMetadata is metadata/chromes mode
	PreviewModeMetadata PreviewMode = "metadata"
)

// EditingPreviewData contains data specific to Pages editing mode
type EditingPreviewData struct {
	PreviewData

	// EditMode specifies the type of editing (metadata, chromes, etc.)
	EditMode string `json:"editMode,omitempty"`
}

// DesignLibraryRenderPreviewData contains data for design library preview
type DesignLibraryRenderPreviewData struct {
	// ComponentName is the name of the component to render
	ComponentName string `json:"componentName"`

	// ComponentProps are the props to pass to the component
	ComponentProps map[string]any `json:"componentProps,omitempty"`

	// LayoutData is optional layout data for the component
	LayoutData map[string]any `json:"layoutData,omitempty"`
}

// IsPreviewMode checks if preview data indicates preview mode
func (pd *PreviewData) IsPreviewMode() bool {
	return pd.Mode == PreviewModePreview || pd.Mode == PreviewModeEdit || pd.Mode == PreviewModeMetadata
}

// IsEditMode checks if preview data indicates edit mode
func (pd *PreviewData) IsEditMode() bool {
	return pd.Mode == PreviewModeEdit || pd.Mode == PreviewModeMetadata
}

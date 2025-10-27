package models

// Page represents a complete page from Sitecore with all associated data
// Note: LayoutData uses interface{} to avoid import cycles
// Cast to *layoutService.LayoutServiceData when using
type Page struct {
	// LayoutData contains the page structure and content from Layout Service
	LayoutData interface{} `json:"layoutData"`

	// Dictionary contains i18n phrases for the page's language
	Dictionary DictionaryPhrases `json:"dictionary,omitempty"`

	// ErrorPages contains custom error page definitions
	ErrorPages *ErrorPages `json:"errorPages,omitempty"`

	// HeadLinks contains HTML link tags for the page (stylesheets, etc.)
	HeadLinks []HTMLLink `json:"headLinks,omitempty"`
}

// DictionaryPhrases maps dictionary keys to their translated values
type DictionaryPhrases map[string]string

// ErrorPages contains custom error page items from Sitecore
// Note: Item uses interface{} to avoid import cycles
type ErrorPages struct {
	// NotFoundPage is the custom 404 page
	NotFoundPage interface{} `json:"notFoundPage,omitempty"`

	// ServerErrorPage is the custom 500 page
	ServerErrorPage interface{} `json:"serverErrorPage,omitempty"`
}

// HTMLLink represents an HTML link element (stylesheet, icon, etc.)
type HTMLLink struct {
	Rel         string `json:"rel"`
	Href        string `json:"href"`
	Type        string `json:"type,omitempty"`
	As          string `json:"as,omitempty"`
	Sizes       string `json:"sizes,omitempty"`
	Media       string `json:"media,omitempty"`
	CrossOrigin string `json:"crossOrigin,omitempty"`
}

// PageOptions contains options for fetching a page
type PageOptions struct {
	// Site name to fetch the page for
	Site string `json:"site"`

	// Locale (language) for the page (e.g., "en", "fr-CA")
	Locale *string `json:"locale,omitempty"`

	// Personalize contains personalization variant information
	Personalize *PersonalizeInfo `json:"personalize,omitempty"`
}

// PageMode represents the mode the page is being rendered in
type PageMode string

const (
	// PageModeNormal is the standard public page mode
	PageModeNormal PageMode = "normal"

	// PageModePreview is preview mode (viewing unpublished content)
	PageModePreview PageMode = "preview"

	// PageModeEdit is editing mode (Sitecore Pages editor)
	PageModeEdit PageMode = "edit"
)

// StaticPath represents a path for static site generation
type StaticPath struct {
	// Site name
	Site string `json:"site"`

	// Locale (language)
	Locale string `json:"locale"`

	// Path segments (e.g., ["about", "team"])
	Path []string `json:"path"`
}

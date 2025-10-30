package layoutservice

// LayoutServicePageState represents the page state enum from Layout Service
type LayoutServicePageState string

const (
	PageStatePreview LayoutServicePageState = "preview"
	PageStateEdit    LayoutServicePageState = "edit"
	PageStateNormal  LayoutServicePageState = "normal"
)

// EditMode represents the edit mode for rendering content in Sitecore Editors
type EditMode string

const (
	EditModeMetadata EditMode = "metadata"
)

// RenderingType represents the editing rendering type
type RenderingType string

const (
	RenderingTypeComponent RenderingType = "component"
)

// Constants for component rendering
const (
	EditingComponentPlaceholder = "editing-componentmode-placeholder"
	EditingComponentID          = "editing-component"
)

// GenericFieldValue represents a field value which can be various types
type GenericFieldValue any

// FieldMetadata contains field metadata in editing mode
type FieldMetadata struct {
	Metadata map[string]any `json:"metadata,omitempty"`
}

// Field represents field value data on a component
type Field struct {
	Value    GenericFieldValue `json:"value"`
	Metadata map[string]any    `json:"metadata,omitempty"`
}

// Item represents content data returned from Layout Service
type Item struct {
	Name        string         `json:"name"`
	DisplayName *string        `json:"displayName,omitempty"`
	ID          *string        `json:"id,omitempty"`
	URL         *string        `json:"url,omitempty"`
	Fields      map[string]any `json:"fields"`
}

// ComponentParams represents component parameters
type ComponentParams map[string]string

// ComponentFields represents content field data passed to a component
type ComponentFields map[string]any

// PlaceholdersData represents placeholder contents data
type PlaceholdersData map[string][]ComponentRendering

// ComponentRendering represents a component instance within a placeholder on a route
type ComponentRendering struct {
	ComponentName string           `json:"componentName"`
	DataSource    *string          `json:"dataSource,omitempty"`
	UID           *string          `json:"uid,omitempty"`
	Placeholders  PlaceholdersData `json:"placeholders,omitempty"`
	Fields        ComponentFields  `json:"fields,omitempty"`
	Params        *ComponentParams `json:"params,omitempty"`
}

// LayoutServiceContext represents the shape of context data from the Sitecore Layout Service
type LayoutServiceContext struct {
	PageEditing                    *bool                   `json:"pageEditing,omitempty"`
	Language                       *string                 `json:"language,omitempty"`
	ItemPath                       *string                 `json:"itemPath,omitempty"`
	PageState                      *LayoutServicePageState `json:"pageState,omitempty"`
	VisitorIdentificationTimestamp *int64                  `json:"visitorIdentificationTimestamp,omitempty"`
	Site                           *struct {
		Name *string `json:"name,omitempty"`
	} `json:"site,omitempty"`
	RenderingType *RenderingType            `json:"renderingType,omitempty"`
	ClientScripts []string                  `json:"clientScripts,omitempty"`
	ClientData    map[string]map[string]any `json:"clientData,omitempty"`
	// Additional dynamic properties
	AdditionalProperties map[string]any `json:"-"`
}

// LayoutServiceContextData contains context information from the Sitecore Layout Service
type LayoutServiceContextData struct {
	Context LayoutServiceContext `json:"context"`
}

// RouteData represents the shape of route data returned from Sitecore Layout Service
type RouteData struct {
	Name         string           `json:"name"`
	DisplayName  *string          `json:"displayName,omitempty"`
	Fields       map[string]any   `json:"fields,omitempty"`
	DatabaseName *string          `json:"databaseName,omitempty"`
	DeviceID     *string          `json:"deviceId,omitempty"`
	ItemLanguage *string          `json:"itemLanguage,omitempty"`
	ItemVersion  *int             `json:"itemVersion,omitempty"`
	LayoutID     *string          `json:"layoutId,omitempty"`
	TemplateID   *string          `json:"templateId,omitempty"`
	TemplateName *string          `json:"templateName,omitempty"`
	Placeholders PlaceholdersData `json:"placeholders"`
	ItemID       *string          `json:"itemId,omitempty"`
}

// LayoutServiceData represents a reply from the Sitecore Layout Service
type LayoutServiceData struct {
	Sitecore struct {
		LayoutServiceContextData
		Route *RouteData `json:"route"`
	} `json:"sitecore"`
}

// PlaceholderData represents the contents of a single placeholder returned from placeholder service
type PlaceholderData struct {
	Name     string               `json:"name"`
	Path     string               `json:"path"`
	Elements []ComponentRendering `json:"elements"`
}

// RouteOptions represents additional route options when requesting layout data
type RouteOptions struct {
	Site   string  `json:"site"`
	Locale *string `json:"locale,omitempty"`
}

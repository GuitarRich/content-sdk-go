# Claude Code Agent Guide for Sitecore Content SDK Go

## Project Overview

This is the **Sitecore Content SDK for Go** - a high-performance, strongly-typed SDK for building server-rendered web applications with Sitecore XM Cloud. The project emphasizes type safety, clean architecture, and proper integration with Sitecore's Experience Editor.

### Mission Critical Rules

When working with this codebase, there are three non-negotiable rules that must ALWAYS be followed:

1. **ALWAYS use SDK field renderer components** - Never manually render Sitecore field values
2. **ALWAYS extract datasource data into strongly-typed structs** - Never use deeply nested type assertions
3. **ALWAYS store both typed fields and raw field data** - For validation and rendering respectively

These patterns are not optional - they are fundamental to how the SDK works and ensures compatibility with Sitecore's Experience Editor.

## Tech Stack

- **Language**: Go 1.21+
- **Template Engine**: Templ (https://templ.guide) - Type-safe HTML templates
- **Testing**: Go testing package + testify assertions
- **Build**: `go build`, `templ generate` for template compilation
- **Linting**: golangci-lint with strict settings
- **Web Framework**: Echo v4 (primary), adaptable to Gin, Chi, etc.

## Architecture Overview

```
content-sdk-go/                  # Core SDK package
├── models/                      # Field types, data structures
│   ├── fields.go               # Field type definitions (TextField, LinkField, etc.)
│   ├── field_extractors.go    # Convert raw JSON to typed fields
│   └── field_helpers.go        # Helper functions for field access
├── components/                  # Templ renderer components
│   ├── fields.templ            # Field renderers (Text, Link, Image, RichText)
│   ├── editing.templ           # Experience Editor chrome markers
│   └── shared.templ            # Shared utilities
├── client/                      # API clients for Sitecore
├── layoutService/               # Layout Service integration
├── handlers/                    # HTTP handlers (catch-all routing, etc.)
└── middleware/                  # HTTP middleware (editing, locale, redirects)

xmcloud-starter-go/             # Example applications
└── examples/
    ├── basic-go/               # Minimal example
    └── kit-go-product-listing/ # Full-featured reference app
        ├── components/         # Application-specific components
        ├── models/             # Application data models
        └── handlers/           # Application handlers
```

## Rule #1: Always Use SDK Field Renderers

### The Problem

Sitecore's Experience Editor needs special HTML markup (chrome markers) around editable fields. If you manually render field values, content authors won't be able to edit content in-page.

### The Solution

Always use SDK component renderers from `github.com/content-sdk-go/components`.

### Correct Pattern

```templ
package components

import sdkcomponents "github.com/content-sdk-go/components"

templ MyComponent(fields interface{}, params map[string]interface{}) {
    if ds := models.ExtractMyDatasource(fields); ds != nil {
        <div>
            // ✅ CORRECT - Uses SDK component
            <h1>
                @sdkcomponents.PlainText(ds.TitleRaw, "title", isEditingMode)
            </h1>
            
            // ✅ CORRECT - Rich text with SDK component
            <div>
                @sdkcomponents.RichText(ds.ContentRaw, "content", isEditingMode, "prose")
            </div>
            
            // ✅ CORRECT - Link with SDK component
            @sdkcomponents.Link(ds.LinkRaw, "cta", isEditingMode, "btn btn-primary")
            
            // ✅ CORRECT - Image with SDK component
            @sdkcomponents.Image(ds.ImageRaw, "hero", isEditingMode, "w-full", "", "")
        </div>
    }
}
```

### Incorrect Pattern

```templ
// ❌ WRONG - Manual rendering breaks Experience Editor
templ MyComponent(fields interface{}, params map[string]interface{}) {
    <h1>{ ds.Title.Value }</h1>  // Missing Experience Editor chrome
    
    <a href={ templ.SafeURL(ds.Link.GetHref()) }>
        { ds.Link.GetText() }  // Missing Experience Editor chrome
    </a>
    
    <img src={ ds.Image.GetSrc() } alt={ ds.Image.GetAlt() } />  // Missing Experience Editor chrome
}
```

### Available SDK Components

```go
// Text field rendering
sdkcomponents.PlainText(fieldData, fieldName, isEditingMode)
sdkcomponents.Text(fieldData, fieldName, isEditingMode, htmlTag, cssClass)

// Rich text field rendering
sdkcomponents.RichText(fieldData, fieldName, isEditingMode, cssClass)

// Link field rendering
sdkcomponents.Link(fieldData, fieldName, isEditingMode, cssClass, children...)

// Image field rendering
sdkcomponents.Image(fieldData, fieldName, isEditingMode, cssClass, width, height)
```

### When SDK Components are Used

- **Experience Editor (sc_mode=edit)**: Renders chrome markers for in-page editing
- **Normal Mode**: Renders clean HTML without editing markup
- **Preview Mode**: May render with minimal chrome depending on configuration

### Why This Matters

Experience Editor chrome markers look like this in the HTML:

```html
<code type="text/sitecore" chrometype="field" kind="open" id="r_FA123" ... ></code>
<h1>Page Title</h1>
<code type="text/sitecore" chrometype="field" kind="close" id="r_FA123"></code>
```

SDK components automatically generate this markup when `isEditingMode` is true. Manual rendering bypasses this, breaking the editing experience.

## Rule #2: Extract to Strongly-Typed Structs

### The Problem

Sitecore Layout Service returns deeply nested `map[string]interface{}` structures. Accessing these requires multiple levels of type assertions, making code complex and error-prone.

### The Solution

Create strongly-typed structs and extraction functions. Extract once, use strongly-typed data throughout your component.

### Pattern: Strongly-Typed Data Structure

```go
package models

import sdkmodels "github.com/content-sdk-go/models"

// LinkListDatasource represents a strongly-typed link list component
type LinkListDatasource struct {
    // Typed fields for validation and logic
    Title *sdkmodels.TextField
    Items []*LinkListItem
    
    // Raw field data for SDK component rendering
    TitleRaw interface{}
}

// LinkListItem represents a single link in the list
type LinkListItem struct {
    Link    *sdkmodels.LinkField  // Typed for validation
    LinkRaw interface{}            // Raw for SDK rendering
}

// HasItems checks if there are any items to display
func (d *LinkListDatasource) HasItems() bool {
    return d != nil && len(d.Items) > 0
}

// GetTitle safely returns the title value
func (d *LinkListDatasource) GetTitle() string {
    if d != nil && d.Title != nil {
        return d.Title.Value
    }
    return ""
}
```

### Pattern: Extraction Function

```go
// ExtractLinkListDatasource converts raw datasource to strongly-typed struct
func ExtractLinkListDatasource(datasource map[string]interface{}) *LinkListDatasource {
    if datasource == nil {
        debug.Common("ExtractLinkListDatasource: datasource is nil")
        return &LinkListDatasource{}
    }
    
    result := &LinkListDatasource{
        Items: make([]*LinkListItem, 0),
    }
    
    // Extract title field
    if field, ok := datasource["field"].(map[string]interface{}); ok {
        titleField := GetFieldByName(field, "title")
        result.Title = sdkmodels.ExtractTextFieldFromMap(titleField)
        result.TitleRaw = titleField  // Store raw for SDK components
    }
    
    // Extract children items
    if children, ok := datasource["children"].(map[string]interface{}); ok {
        if results, ok := children["results"].([]interface{}); ok {
            for _, item := range results {
                if itemMap, ok := item.(map[string]interface{}); ok {
                    linkItem := extractLinkListItem(itemMap)
                    if linkItem != nil {
                        result.Items = append(result.Items, linkItem)
                    }
                }
            }
        }
    }
    
    return result
}

// extractLinkListItem extracts a single item (unexported helper)
func extractLinkListItem(itemMap map[string]interface{}) *LinkListItem {
    if fieldData, ok := itemMap["field"].(map[string]interface{}); ok {
        if linkField := fieldData["link"]; linkField != nil {
            return &LinkListItem{
                Link:    sdkmodels.ExtractLinkFieldFromMap(linkField),
                LinkRaw: linkField,
            }
        }
    }
    return nil
}
```

### Pattern: Template Usage

```templ
templ LinkList(fields interface{}, params map[string]interface{}) {
    if datasource := models.ExtractDatasourceField(fields); datasource != nil {
        if isEditingMode := GetEditingMode(params); true {
            // Extract ONCE to strongly-typed data
            if ds := models.ExtractLinkListDatasource(datasource); ds.HasItems() {
                <div class="link-list">
                    <h3>
                        // Use raw data with SDK component
                        @sdkcomponents.PlainText(ds.TitleRaw, "title", isEditingMode)
                    </h3>
                    <ul>
                        // Iterate over strongly-typed items
                        for _, item := range ds.Items {
                            // Use typed field for validation
                            if item.Link != nil && !item.Link.IsEmpty() {
                                <li>
                                    // Use raw data with SDK component
                                    @sdkcomponents.Link(item.LinkRaw, "link", isEditingMode, "")
                                </li>
                            }
                        }
                    </ul>
                </div>
            }
        }
    }
}
```

### Benefits of This Pattern

**Before (Nested Type Assertions)**:
```go
// ❌ BAD - 5+ levels of nesting, repeated everywhere
if data, ok := datasource["data"].(map[string]interface{}); ok {
    if ds, ok := data["datasource"].(map[string]interface{}); ok {
        if children, ok := ds["children"].(map[string]interface{}); ok {
            if results, ok := children["results"].([]interface{}); ok {
                for _, item := range results {
                    if itemMap, ok := item.(map[string]interface{}); ok {
                        if fieldData, ok := itemMap["field"].(map[string]interface{}); ok {
                            if link := fieldData["link"]; link != nil {
                                // Finally can use the data
                            }
                        }
                    }
                }
            }
        }
    }
}
```

**After (Strongly-Typed Extraction)**:
```go
// ✅ GOOD - Extract once, use throughout
ds := models.ExtractLinkListDatasource(datasource)
for _, item := range ds.Items {
    if item.Link != nil && !item.Link.IsEmpty() {
        // Use strongly-typed item.Link
    }
}
```

**Improvements**:
- ✅ **Reduced complexity**: From 5-6 levels to 1 extraction call
- ✅ **Type safety**: Compile-time checking vs runtime assertions
- ✅ **Readability**: Clear data structure vs nested conditionals
- ✅ **Reusability**: Extraction logic in one place
- ✅ **Maintainability**: Changes isolated to extraction function
- ✅ **IDE support**: Better autocomplete and refactoring
- ✅ **Testability**: Easy to unit test extraction logic

## Rule #3: Hybrid Approach - Typed + Raw Data

### The Pattern

Store BOTH strongly-typed fields AND raw field data in your structs.

```go
type ComponentDatasource struct {
    // Strongly-typed fields for validation, checks, logic
    Title       *sdkmodels.TextField
    Description *sdkmodels.RichTextField
    Image       *sdkmodels.ImageField
    Link        *sdkmodels.LinkField
    
    // Raw field data for SDK component rendering
    TitleRaw       interface{}
    DescriptionRaw interface{}
    ImageRaw       interface{}
    LinkRaw        interface{}
}
```

### Why Both?

**Typed Fields** - Use for:
- Validation (`if field.IsEmpty()`)
- Conditionals (`if field.Value != ""`)
- Iteration and counting
- Business logic
- Type safety in Go code

**Raw Fields** - Use for:
- Passing to SDK component renderers
- Preserving Experience Editor metadata
- Supporting all Sitecore field configurations

### Extraction Pattern

```go
func ExtractComponentDatasource(datasource map[string]interface{}) *ComponentDatasource {
    result := &ComponentDatasource{}
    
    if field, ok := datasource["field"].(map[string]interface{}); ok {
        // Title field
        titleField := GetFieldByName(field, "title")
        result.Title = sdkmodels.ExtractTextFieldFromMap(titleField)
        result.TitleRaw = titleField  // Store raw
        
        // Description field
        descField := GetFieldByName(field, "description")
        result.Description = sdkmodels.ExtractRichTextFieldFromMap(descField)
        result.DescriptionRaw = descField  // Store raw
        
        // Image field
        imageField := GetFieldByName(field, "image")
        result.Image = sdkmodels.ExtractImageFieldFromMap(imageField)
        result.ImageRaw = imageField  // Store raw
        
        // Link field
        linkField := GetFieldByName(field, "link")
        result.Link = sdkmodels.ExtractLinkFieldFromMap(linkField)
        result.LinkRaw = linkField  // Store raw
    }
    
    return result
}
```

### Usage in Templates

```templ
templ Component(fields interface{}, params map[string]interface{}) {
    if ds := models.ExtractComponentDatasource(fields); ds != nil {
        <div>
            // Use typed field for validation
            if ds.Title != nil && !ds.Title.IsEmpty() {
                <h1>
                    // Use raw field for rendering
                    @sdkcomponents.PlainText(ds.TitleRaw, "title", isEditingMode)
                </h1>
            }
            
            // Use typed field for business logic
            if ds.Image != nil && !ds.Image.IsEmpty() {
                <figure>
                    // Use raw field for rendering
                    @sdkcomponents.Image(ds.ImageRaw, "hero", isEditingMode, "responsive", "", "")
                    if ds.Description != nil && !ds.Description.IsEmpty() {
                        <figcaption>
                            @sdkcomponents.RichText(ds.DescriptionRaw, "caption", isEditingMode, "")
                        </figcaption>
                    }
                </figure>
            }
            
            // Use typed field for conditional rendering
            if ds.Link != nil && !ds.Link.IsEmpty() && ds.Link.GetHref() != "" {
                // Use raw field for SDK component
                @sdkcomponents.Link(ds.LinkRaw, "cta", isEditingMode, "btn btn-primary")
            }
        </div>
    }
}
```

## Templ Template Syntax

### Component Definition

```templ
// Package must match directory
package components

// Imports
import (
    sdkcomponents "github.com/content-sdk-go/components"
    "github.com/xmcloud-starter-go/examples/kit-go-product-listing/models"
)

// Component definition with typed parameters
templ ComponentName(fields interface{}, params map[string]interface{}) {
    // Template content
}
```

### Calling Components

```templ
// Call SDK components with @ prefix
@sdkcomponents.PlainText(fieldData, "fieldName", isEditingMode)
@sdkcomponents.Link(linkData, "linkName", isEditingMode, "css-class")

// Call other templ components
@HeaderComponent(data, params)
@FooterComponent(data)
```

### Control Flow

```templ
// If statements (no @ prefix for control flow)
if condition {
    <div>Content when true</div>
}

// If-else
if condition {
    <div>True branch</div>
} else {
    <div>False branch</div>
}

// For loops
for index, item := range items {
    <div>{ item.Name }</div>
}

// For with single variable
for _, item := range items {
    @ItemComponent(item)
}

// Switch statements
switch value {
    case "option1":
        <div>Option 1</div>
    case "option2":
        <div>Option 2</div>
    default:
        <div>Default</div>
}
```

### Expressions and Interpolation

```templ
// Go expressions in curly braces
{ variableName }
{ structField.Value }
{ function(param) }

// In attributes
<div class={ className }>
<a href={ templ.SafeURL(url) }>
<img src={ templ.SafeURL(imageSrc) } alt={ altText } />

// Multiple CSS classes
<div class={ "base-class", conditionalClass, "another-class" }>
<div class={ getClassName(index) }>
```

### Special Functions

```templ
// Raw HTML (use sparingly, ensure content is safe)
@templ.Raw(htmlString)

// URL sanitization (ALWAYS use for href attributes)
href={ templ.SafeURL(url) }

// URL sanitization for other URL attributes
src={ templ.SafeURL(imageSrc) }
```

### Best Practices

```templ
// ✅ GOOD - Clean component structure
templ LinkList(fields interface{}, params map[string]interface{}) {
    if datasource := models.ExtractDatasourceField(fields); datasource != nil {
        if ds := models.ExtractLinkListDatasource(datasource); ds.HasItems() {
            @linkListContent(ds, params, isEditingMode)
        }
    }
}

// Sub-components for complex rendering
templ linkListContent(ds *models.LinkListDatasource, params map[string]interface{}, isEditingMode bool) {
    <div class="link-list">
        <h3>
            @sdkcomponents.PlainText(ds.TitleRaw, "title", isEditingMode)
        </h3>
        <ul>
            for _, item := range ds.Items {
                @linkListItem(item, isEditingMode)
            }
        </ul>
    </div>
}

templ linkListItem(item *models.LinkListItem, isEditingMode bool) {
    if item.Link != nil && !item.Link.IsEmpty() {
        <li>
            @sdkcomponents.Link(item.LinkRaw, "link", isEditingMode, "")
        </li>
    }
}

// ❌ BAD - Complex logic in templates
templ Component(fields interface{}, params map[string]interface{}) {
    if data, ok := fields.(map[string]interface{}); ok {
        // Don't do complex type assertions in templates
        if nested, ok := data["nested"].(map[string]interface{}); ok {
            // Move this to extraction function
        }
    }
}

// ✅ GOOD - Logic in Go, simple rendering in template
func (ds *ComponentDatasource) ShouldDisplay() bool {
    return ds != nil && ds.Title != nil && !ds.Title.IsEmpty()
}

templ Component(ds *models.ComponentDatasource, isEditingMode bool) {
    if ds.ShouldDisplay() {
        @componentContent(ds, isEditingMode)
    }
}
```

## Field Type Reference

### TextField

Single-line and multi-line text fields.

```go
type TextField struct {
    Value    string  // The text content
    Editable string  // Experience Editor metadata
}

// Methods
func (f *TextField) GetValue() interface{}
func (f *TextField) GetEditable() string
func (f *TextField) IsEmpty() bool
```

Usage:
```templ
@sdkcomponents.PlainText(ds.TitleRaw, "title", isEditingMode)
@sdkcomponents.Text(ds.SubtitleRaw, "subtitle", isEditingMode, "h2", "text-lg")
```

### RichTextField

Rich text with HTML content.

```go
type RichTextField struct {
    Value    string  // HTML content
    Editable string  // Experience Editor metadata
}

// Methods
func (f *RichTextField) GetValue() interface{}
func (f *RichTextField) GetEditable() string
func (f *RichTextField) IsEmpty() bool
```

Usage:
```templ
@sdkcomponents.RichText(ds.ContentRaw, "content", isEditingMode, "prose prose-lg")
```

### LinkField

General link field (internal, external, mailto, etc).

```go
type LinkField struct {
    Href     string  // Link URL
    Text     string  // Link text
    Target   string  // Link target (_blank, etc)
    Title    string  // Link title attribute
    Class    string  // CSS classes from Sitecore
    Editable string  // Experience Editor metadata
    Value    *LinkFieldValue  // Nested structure
}

// Methods
func (f *LinkField) GetHref() string
func (f *LinkField) GetText() string
func (f *LinkField) GetTarget() string
func (f *LinkField) GetTitle() string
func (f *LinkField) GetClass() string
func (f *LinkField) GetValue() interface{}
func (f *LinkField) GetEditable() string
func (f *LinkField) IsEmpty() bool
```

Usage:
```templ
@sdkcomponents.Link(ds.CTALinkRaw, "cta", isEditingMode, "btn btn-primary")

// With children
@sdkcomponents.Link(ds.LinkRaw, "link", isEditingMode, "card-link") {
    <span class="icon">→</span>
    <span>Learn More</span>
}
```

### ImageField

Image field with responsive image support.

```go
type ImageField struct {
    Src      string  // Image URL
    Alt      string  // Alt text
    Width    string  // Width attribute
    Height   string  // Height attribute
    Editable string  // Experience Editor metadata
    Value    *ImageFieldValue  // Nested structure
}

// Methods
func (f *ImageField) GetSrc() string
func (f *ImageField) GetAlt() string
func (f *ImageField) GetWidth() string
func (f *ImageField) GetHeight() string
func (f *ImageField) GetValue() interface{}
func (f *ImageField) GetEditable() string
func (f *ImageField) IsEmpty() bool
```

Usage:
```templ
// Basic usage
@sdkcomponents.Image(ds.ImageRaw, "hero", isEditingMode, "w-full", "", "")

// With specific dimensions
@sdkcomponents.Image(ds.ThumbnailRaw, "thumbnail", isEditingMode, "thumbnail", "200", "200")

// Responsive image
@sdkcomponents.Image(ds.BannerRaw, "banner", isEditingMode, "responsive", "", "")
```

## Common Patterns

### Component with Variants

```go
// Model
type ContainerVariant string

const (
    VariantDefault     ContainerVariant = "default"
    VariantFullWidth   ContainerVariant = "full-width"
    VariantCentered    ContainerVariant = "centered"
)

func GetVariant(params map[string]interface{}) ContainerVariant {
    if variant, ok := params["Variant"].(string); ok {
        return ContainerVariant(variant)
    }
    return VariantDefault
}
```

```templ
// Template
templ Container(fields interface{}, params map[string]interface{}) {
    switch models.GetVariant(params) {
        case models.VariantFullWidth:
            @containerFullWidth(fields, params)
        case models.VariantCentered:
            @containerCentered(fields, params)
        default:
            @containerDefault(fields, params)
    }
}
```

### Responsive Images

```go
func (ds *HeroDatasource) GetResponsiveImageSizes() string {
    if ds.Layout == "full-width" {
        return "100vw"
    }
    return "(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
}
```

```templ
templ HeroImage(ds *models.HeroDatasource, isEditingMode bool) {
    <picture>
        @sdkcomponents.Image(
            ds.ImageRaw,
            "hero",
            isEditingMode,
            "responsive",
            "",
            "",
        )
    </picture>
}
```

### Conditional CTA

```go
func (ds *ComponentDatasource) HasCTA() bool {
    return ds.CTALink != nil &&
           !ds.CTALink.IsEmpty() &&
           ds.CTAText != nil &&
           !ds.CTAText.IsEmpty()
}
```

```templ
templ ComponentCTA(ds *models.ComponentDatasource, isEditingMode bool) {
    if ds.HasCTA() {
        <div class="cta">
            @sdkcomponents.Link(ds.CTALinkRaw, "cta", isEditingMode, "btn btn-lg") {
                @sdkcomponents.PlainText(ds.CTATextRaw, "cta-text", isEditingMode)
            }
        </div>
    }
}
```

## Testing Best Practices

### Unit Test Example

```go
func TestExtractComponentDatasource(t *testing.T) {
    tests := []struct {
        name    string
        input   map[string]interface{}
        want    *ComponentDatasource
        wantErr bool
    }{
        {
            name: "valid datasource with all fields",
            input: map[string]interface{}{
                "field": map[string]interface{}{
                    "title": map[string]interface{}{
                        "jsonValue": map[string]interface{}{
                            "value": "Test Title",
                        },
                    },
                    "link": map[string]interface{}{
                        "jsonValue": map[string]interface{}{
                            "value": map[string]interface{}{
                                "href": "/test",
                                "text": "Test Link",
                            },
                        },
                    },
                },
            },
            want: &ComponentDatasource{
                Title: &sdkmodels.TextField{Value: "Test Title"},
                // ... other assertions
            },
            wantErr: false,
        },
        {
            name:    "nil datasource",
            input:   nil,
            want:    &ComponentDatasource{},
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := ExtractComponentDatasource(tt.input)
            
            if got == nil {
                t.Fatal("Expected non-nil result")
            }
            
            if got.Title != nil && tt.want.Title != nil {
                assert.Equal(t, tt.want.Title.Value, got.Title.Value)
            }
            
            // More assertions...
        })
    }
}
```

## Error Handling

### Always Handle Errors

```go
// ✅ GOOD
func GetPage(ctx context.Context, path string) (*Page, error) {
    data, err := layoutService.FetchPage(ctx, path)
    if err != nil {
        return nil, fmt.Errorf("fetch page for path %s: %w", path, err)
    }
    
    page, err := parsePage(data)
    if err != nil {
        return nil, fmt.Errorf("parse page data: %w", err)
    }
    
    return page, nil
}

// ❌ BAD - Never ignore errors
data, _ := layoutService.FetchPage(ctx, path)

// ❌ BAD - Never panic in production code
if err != nil {
    panic(err)
}
```

### Guard Clauses

```go
// ✅ GOOD - Guard clauses at the start
func ProcessField(field interface{}) string {
    if field == nil {
        return ""
    }
    
    fieldMap, ok := field.(map[string]interface{})
    if !ok {
        return ""
    }
    
    value, ok := fieldMap["value"].(string)
    if !ok {
        return ""
    }
    
    return value
}
```

## Documentation Standards

### Godoc Comments

```go
// ComponentDatasource represents the strongly-typed structure for a component's data.
// It contains both typed fields for validation/logic and raw data for SDK component rendering.
//
// The typed fields (Title, Image, Link) should be used for:
//   - Validation and emptiness checks
//   - Conditional logic
//   - Type-safe access to field values
//
// The raw fields (TitleRaw, ImageRaw, LinkRaw) should be used for:
//   - Passing to SDK component renderers
//   - Preserving Experience Editor metadata
type ComponentDatasource struct {
    Title    *sdkmodels.TextField
    TitleRaw interface{}
    // ... other fields
}

// HasTitle checks if the component has a non-empty title field.
// Returns true if Title is not nil and has a non-empty value.
func (d *ComponentDatasource) HasTitle() bool {
    return d != nil && d.Title != nil && !d.Title.IsEmpty()
}

// ExtractComponentDatasource converts a generic datasource map into a strongly-typed ComponentDatasource.
// It handles the standard Sitecore Layout Service data structure and extracts all relevant fields.
//
// Parameters:
//   - datasource: The raw datasource map from Sitecore Layout Service
//
// Returns:
//   - *ComponentDatasource: A strongly-typed datasource with extracted fields
//
// Example:
//   ds := ExtractComponentDatasource(rawDatasource)
//   if ds.HasTitle() {
//       fmt.Println(ds.Title.Value)
//   }
func ExtractComponentDatasource(datasource map[string]interface{}) *ComponentDatasource {
    // ... implementation
}
```

## Summary: The Three Rules

### 1. Always Use SDK Field Renderers

```templ
// ✅ DO THIS
@sdkcomponents.PlainText(ds.TitleRaw, "title", isEditingMode)
@sdkcomponents.Link(item.LinkRaw, "link", isEditingMode, "")

// ❌ NOT THIS
{ ds.Title.Value }
<a href={link.GetHref()}>{ link.GetText() }</a>
```

### 2. Extract to Strongly-Typed Structs

```go
// ✅ DO THIS
type LinkListDatasource struct {
    Title    *sdkmodels.TextField
    TitleRaw interface{}
    Items    []*LinkListItem
}

ds := models.ExtractLinkListDatasource(datasource)
for _, item := range ds.Items { ... }

// ❌ NOT THIS
if children, ok := datasource["children"].(map[string]interface{}); ok {
    if results, ok := children["results"].([]interface{}); ok {
        // Multiple levels of nesting
    }
}
```

### 3. Store Both Typed and Raw Data

```go
// ✅ DO THIS
type ComponentDatasource struct {
    Title    *sdkmodels.TextField  // For validation
    TitleRaw interface{}            // For SDK rendering
}

// Use typed for logic
if ds.Title != nil && !ds.Title.IsEmpty() {
    // Use raw for rendering
    @sdkcomponents.PlainText(ds.TitleRaw, "title", isEditingMode)
}
```

These three rules work together to create maintainable, type-safe code that properly integrates with Sitecore's Experience Editor. Follow them consistently and your code will be clean, testable, and work correctly in both editing and normal modes.


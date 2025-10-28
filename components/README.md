# Content SDK Go - Shared Components

## Overview

This package provides reusable templ components that are shared across all Sitecore applications using the Content SDK Go.

These components handle common rendering scenarios and provide visual feedback during development.

---

## Components

### UnknownComponent

Renders a fallback component when a Sitecore component is not registered in the component registry.

**Purpose:**

- Provides visual feedback in development
- Helps identify unimplemented components
- Prevents application crashes from missing components

**Usage:**

```go
import sdkcomponents "github.com/content-sdk-go/components"

// In your renderer when a component is not found
return sdkcomponents.UnknownComponent(componentName, fields)
```

**Rendered Output:**

- Yellow bordered warning box
- Warning icon
- Component name display
- Helpful message for developers

**Example:**

```html
<div
  class="unknown-component border-2 border-dashed border-yellow-500 bg-yellow-50 p-4 rounded my-4"
>
  <strong>Component Not Implemented</strong>
  <p>Component: <code>MyCustomComponent</code></p>
  <p>
    This component needs to be implemented. Register it in your component
    registry.
  </p>
</div>
```

---

### RenderPlaceholder

Renders a Sitecore placeholder with all its nested components.

**Purpose:**

- Renders placeholder containers
- Iterates through and renders all nested components
- Adds data attributes for debugging

**Usage:**

```go
import sdkcomponents "github.com/content-sdk-go/components"

// In your layout template
@sdkcomponents.RenderPlaceholder("main-content", placeholders["main-content"])
```

**Parameters:**

- `name string` - The placeholder name (e.g., "main-content", "sidebar")
- `components []templ.Component` - Array of components to render

**Rendered Output:**

```html
<div class="placeholder-main-content" data-placeholder="main-content">
  <!-- Component 1 -->
  <!-- Component 2 -->
  <!-- Component 3 -->
</div>
```

**Benefits:**

- Semantic HTML with placeholder identification
- Data attributes for debugging and Experience Editor
- Consistent placeholder rendering across applications

---

### RenderComponentWithMetadata

Renders a component with HTML comments containing metadata for debugging.

**Purpose:**

- Adds HTML comments around components for debugging
- Shows component name, UID, and DataSource in rendered HTML
- Useful in development and troubleshooting

**Usage:**

```go
import sdkcomponents "github.com/content-sdk-go/components"

// In your renderer
@sdkcomponents.RenderComponentWithMetadata(component, componentRendering)
```

**Parameters:**

- `component templ.Component` - The component to render
- `componentRendering *layoutservice.ComponentRendering` - The component metadata

**Rendered Output:**

```html
<!-- Component: RichText -->
<!-- UID: abc-123-def-456 -->
<!-- DataSource: {ABC-DEF-GHI} -->
<div class="rich-text">
  <!-- Component content -->
</div>
<!-- End Component: RichText -->
```

**Benefits:**

- Easy component identification in browser DevTools
- Helps debug layout issues
- Shows Sitecore item IDs for reference

---

### RenderEmptyPlaceholder

Renders a visual indicator for empty placeholders in editing mode.

**Purpose:**

- Shows placeholder boundaries in Sitecore Pages editor
- Helps content authors identify where to add components
- Only visible in editing mode

**Usage:**

```go
import sdkcomponents "github.com/content-sdk-go/components"

// When placeholder has no components
@sdkcomponents.RenderEmptyPlaceholder("sidebar", isEditingMode)
```

**Parameters:**

- `name string` - The placeholder name
- `editingMode bool` - Whether the page is in editing mode

**Rendered Output (Editing Mode):**

```html
<div class="placeholder-sidebar" data-placeholder="sidebar">
  <div
    class="sc-jss-empty-placeholder border border-dashed border-gray-300 p-4 text-center"
  >
    <p>Empty placeholder: <strong>sidebar</strong></p>
    <p class="text-xs">Add components in Sitecore Pages</p>
  </div>
</div>
```

**Rendered Output (Normal Mode):**

```html
<div class="placeholder-sidebar" data-placeholder="sidebar">
  <!-- Empty, no visual indicator -->
</div>
```

---

## Integration

### Installation

The components are part of the Content SDK Go module. If you're using the SDK, they're already available.

```bash
# Ensure you have the latest SDK
go get -u github.com/content-sdk-go@latest

# Install templ (required for compilation)
go get github.com/a-h/templ@latest
go install github.com/a-h/templ/cmd/templ@latest
```

### Importing

```go
import sdkcomponents "github.com/content-sdk-go/components"
```

### Using in Application

**1. In Component Registry:**

```go
import sdkcomponents "github.com/content-sdk-go/components"

// Register fallback for unknown components
r.Register("Unknown", func(fields interface{}, params map[string]interface{}) templ.Component {
    componentName := params["componentName"].(string)
    return sdkcomponents.UnknownComponent(componentName, fields)
})
```

**2. In Renderer:**

```go
import sdkcomponents "github.com/content-sdk-go/components"

func (r *Renderer) renderComponent(ctx context.Context, rendering *layoutservice.ComponentRendering) templ.Component {
    fn, err := r.registry.Get(rendering.ComponentName)
    if err != nil {
        // Return fallback component
        return sdkcomponents.UnknownComponent(rendering.ComponentName, rendering)
    }
    return fn(rendering.Fields, rendering.Params)
}
```

**3. In Layout Template:**

```go
// components/layout.templ
package components

import sdkcomponents "github.com/content-sdk-go/components"

templ Layout(page *models.Page, placeholders map[string][]templ.Component) {
    <html>
        <body>
            <header>
                @sdkcomponents.RenderPlaceholder("header", placeholders["header"])
            </header>
            <main>
                @sdkcomponents.RenderPlaceholder("main", placeholders["main"])
            </main>
        </body>
    </html>
}
```

---

## CSS Classes

The components use Tailwind CSS classes. Ensure these are available in your application:

### UnknownComponent Classes

- `border-2`, `border-dashed`, `border-yellow-500`
- `bg-yellow-50`, `bg-yellow-100`
- `text-yellow-600`, `text-yellow-700`, `text-yellow-800`
- `p-4`, `px-2`, `py-1`
- `rounded`, `my-4`
- `flex`, `items-center`, `gap-2`
- `text-sm`, `text-xs`

### RenderEmptyPlaceholder Classes

- `border`, `border-dashed`, `border-gray-300`
- `p-4`
- `text-gray-500`
- `text-sm`, `text-xs`
- `text-center`
- `mt-1`

---

## Development

### Generating Templ Files

After modifying `.templ` files, regenerate the Go code:

```bash
cd content-sdk-go/components
templ generate
```

### Testing

Components can be tested by importing them in a test application:

```go
package main

import (
    sdkcomponents "github.com/content-sdk-go/components"
)

func main() {
    // Test UnknownComponent
    component := sdkcomponents.UnknownComponent("TestComponent", nil)

    // Render to test
    html, _ := component.Render(context.Background(), os.Stdout)
}
```

---

## Best Practices

### 1. Always Use SDK Components for Common Scenarios

✅ **Good:**

```go
import sdkcomponents "github.com/content-sdk-go/components"
return sdkcomponents.UnknownComponent(name, data)
```

❌ **Bad:**

```go
// Don't duplicate UnknownComponent in your application
return components.MyCustomUnknownComponent(name, data)
```

### 2. Keep Application-Specific Components Separate

SDK components are for **reusable, common** scenarios. Application-specific components belong in your application's `components/` directory.

**SDK Components (content-sdk-go/components):**

- UnknownComponent
- RenderPlaceholder
- RenderComponentWithMetadata
- RenderEmptyPlaceholder

**Application Components (your-app/components):**

- Button
- Card
- Hero
- Custom business components

### 3. Use Consistent Naming

When using SDK components, always alias the import:

```go
import sdkcomponents "github.com/content-sdk-go/components"
```

This makes it clear which components come from the SDK vs. your application.

---

## Migration from Application Components

If you have existing applications with `UnknownComponent` or similar utilities in your application's `components/` directory:

### Step 1: Update Imports

```go
// Before
import "yourapp/components"
components.UnknownComponent(name, data)

// After
import sdkcomponents "github.com/content-sdk-go/components"
sdkcomponents.UnknownComponent(name, data)
```

### Step 2: Remove Duplicate Code

Delete your application's version of:

- `UnknownComponent`
- `RenderPlaceholder`
- Any other components now in the SDK

### Step 3: Update Templates

```go
// Before (layout.templ)
@RenderPlaceholder("main", placeholders["main"])

// After
@sdkcomponents.RenderPlaceholder("main", placeholders["main"])
```

### Step 4: Regenerate and Build

```bash
templ generate
go build
```

---

## Troubleshooting

### Issue: "Cannot find package github.com/content-sdk-go/components"

**Solution:** Update your SDK dependency:

```bash
go get -u github.com/content-sdk-go@latest
go mod tidy
```

### Issue: "undefined: RenderPlaceholder"

**Solution:** Import the SDK components:

```go
import sdkcomponents "github.com/content-sdk-go/components"
// Then use: sdkcomponents.RenderPlaceholder(...)
```

### Issue: "templ generate fails"

**Solution:** Ensure templ is installed:

```bash
go install github.com/a-h/templ/cmd/templ@latest
```

---

## Version Compatibility

- **Go:** 1.21+
- **Templ:** 0.3.960+
- **Content SDK Go:** Latest version

---

## Contributing

When adding new shared components:

1. **Ensure Reusability:** Component must be useful across multiple applications
2. **Document Thoroughly:** Add usage examples and parameter descriptions
3. **Test:** Verify component works in at least 2 different applications
4. **Generate:** Run `templ generate` after changes
5. **Update README:** Document the new component in this file

---

## License

Part of the Sitecore Content SDK Go - see main SDK license.

---

**For SDK Documentation:** See [Content SDK Go README](../README.md)

**For Application Integration:** See your application's documentation

# AI Code Assistant Instructions

This repository includes comprehensive instruction files for various AI coding assistants. These files ensure that AI-generated code follows best practices and correctly integrates with the Sitecore Content SDK for Go.

## Available Instruction Files

### `.cursorrules`
**For**: Cursor IDE
**Location**: Repository root
**Auto-detected**: Yes

Cursor automatically loads this file when you open the repository. Contains comprehensive rules for:
- SDK field renderer usage
- Strongly-typed data patterns  
- Go naming conventions
- Templ template syntax
- Error handling patterns

### `copilot-instructions.md`
**For**: GitHub Copilot
**Location**: Repository root
**Auto-detected**: Partial (depends on Copilot version)

Comprehensive guide for GitHub Copilot including:
- Project structure and tech stack
- Detailed code patterns with examples
- Common anti-patterns to avoid
- Testing best practices
- Complete API reference

### `CLAUDE.md`
**For**: Claude (Anthropic) / Cursor with Claude
**Location**: Repository root  
**Auto-detected**: Manual reference

In-depth guide focused on the three critical rules:
1. Always use SDK field renderers
2. Extract to strongly-typed structs
3. Store both typed and raw field data

Includes extensive examples, comparisons, and explanations of why these patterns matter.

### `LLMs.txt`
**For**: Any LLM/AI assistant
**Location**: Repository root
**Auto-detected**: No (manual reference)

Concise reference guide with:
- Quick syntax references
- Essential patterns
- Common anti-patterns
- Code templates
- Key takeaways

Use this when you need to quickly reference the correct patterns or copy examples into prompts.

## The Three Critical Rules

All instruction files emphasize these three non-negotiable rules:

### 1. Always Use SDK Field Renderer Components

```templ
// ✅ CORRECT
@sdkcomponents.PlainText(ds.TitleRaw, "title", isEditingMode)
@sdkcomponents.Link(item.LinkRaw, "link", isEditingMode, "")

// ❌ WRONG
{ ds.Title.Value }
<a href={link.GetHref()}>{ link.GetText() }</a>
```

**Why**: SDK components handle Sitecore Experience Editor chrome markers for in-page editing.

### 2. Extract Datasource to Strongly-Typed Structs

```go
// ✅ CORRECT
type LinkListDatasource struct {
    Title    *sdkmodels.TextField
    TitleRaw interface{}
    Items    []*LinkListItem
}

ds := models.ExtractLinkListDatasource(datasource)
for _, item := range ds.Items { ... }

// ❌ WRONG
if children, ok := datasource["children"].(map[string]interface{}); ok {
    if results, ok := children["results"].([]interface{}); ok {
        // Multiple levels of nesting
    }
}
```

**Why**: Reduces complexity from 5+ nested type assertions to a single extraction call.

### 3. Store Both Typed and Raw Field Data

```go
// ✅ CORRECT
type ComponentDatasource struct {
    Title    *sdkmodels.TextField  // For validation and logic
    TitleRaw interface{}            // For SDK rendering
}

// Use typed for logic, raw for rendering
if ds.Title != nil && !ds.Title.IsEmpty() {
    @sdkcomponents.PlainText(ds.TitleRaw, "title", isEditingMode)
}
```

**Why**: Typed fields for validation/business logic, raw fields for proper SDK rendering.

## Using the Instructions

### With Cursor

1. Open the repository in Cursor
2. `.cursorrules` is automatically loaded
3. Reference `CLAUDE.md` for detailed explanations
4. Use `LLMs.txt` for quick code templates

### With GitHub Copilot

1. `copilot-instructions.md` may be auto-detected
2. If not, reference it in your workspace
3. Copy patterns from the file into your code comments for better suggestions

### With Claude (API/Web)

1. Reference `CLAUDE.md` when starting a conversation
2. Copy specific sections as needed
3. Use examples to show desired patterns

### With Other AI Assistants

1. Start with `LLMs.txt` for a concise overview
2. Reference specific sections from other files as needed
3. Copy code templates and adapt to your needs

## File Comparison

| File | Length | Detail Level | Best For |
|------|--------|--------------|----------|
| `.cursorrules` | Medium | Detailed rules | Cursor IDE |
| `copilot-instructions.md` | Long | Very detailed | GitHub Copilot |
| `CLAUDE.md` | Very Long | Extremely detailed | Claude conversations |
| `LLMs.txt` | Medium | Concise reference | Quick reference, any LLM |

## Key Topics Covered

All instruction files cover:

- ✅ SDK field renderer components
- ✅ Strongly-typed data structures
- ✅ Hybrid typed + raw data pattern
- ✅ Templ template syntax
- ✅ Go naming conventions
- ✅ Error handling patterns
- ✅ Testing best practices
- ✅ Common anti-patterns
- ✅ Documentation standards

## Quick Start

When starting with AI assistance on this project:

1. **Read** one of the instruction files appropriate for your AI tool
2. **Understand** the three critical rules (above)
3. **Reference** existing code in `xmcloud-starter-go/examples/kit-go-product-listing/` for patterns
4. **Follow** the strongly-typed extraction pattern consistently
5. **Always** use SDK component renderers for Sitecore fields

## Example Workflow

### Creating a New Component

1. **Define the model** (strongly-typed struct):
```go
type MyComponentDatasource struct {
    Title    *sdkmodels.TextField
    TitleRaw interface{}
    // ... other fields
}
```

2. **Create extraction function**:
```go
func ExtractMyComponentDatasource(datasource map[string]interface{}) *MyComponentDatasource {
    // Extract fields, populate both typed and raw
}
```

3. **Add helper methods**:
```go
func (d *MyComponentDatasource) HasTitle() bool {
    return d != nil && d.Title != nil && !d.Title.IsEmpty()
}
```

4. **Create Templ template**:
```templ
templ MyComponent(fields interface{}, params map[string]interface{}) {
    if ds := models.ExtractMyComponentDatasource(datasource); ds != nil {
        // Use SDK components for rendering
    }
}
```

5. **Write tests**:
```go
func TestExtractMyComponentDatasource(t *testing.T) {
    // Test extraction logic
}
```

## Additional Resources

- **Templ Documentation**: https://templ.guide
- **Go Best Practices**: https://golang.org/doc/effective_go
- **Sitecore XM Cloud**: https://doc.sitecore.com/xmc

## Feedback

If you find that AI assistants consistently generate incorrect patterns, please update the relevant instruction file(s) to clarify the correct approach.

The goal is to have AI assistants generate code that:
- ✅ Works correctly in Sitecore Experience Editor
- ✅ Uses strongly-typed data structures  
- ✅ Follows Go best practices
- ✅ Is maintainable and testable
- ✅ Is well-documented

## Summary

These instruction files ensure that AI code generation follows the established patterns for:
1. **Proper Experience Editor integration** via SDK components
2. **Type safety and maintainability** via strongly-typed structs
3. **Clean, readable code** via the hybrid typed+raw pattern

Remember: The three rules are not optional—they are fundamental to how this SDK integrates with Sitecore XM Cloud.


package models

import (
	"encoding/json"

	"github.com/guitarrich/content-sdk-go/debug"
)

// extractMetadata extracts field metadata from the field map
// Metadata is only present in editing mode
func extractMetadata(fieldMap map[string]any) *FieldMetadata {
	metadataRaw, ok := fieldMap["metadata"].(map[string]any)
	if !ok {
		return nil
	}

	// Use JSON marshal/unmarshal for easy conversion
	jsonBytes, err := json.Marshal(metadataRaw)
	if err != nil {
		debug.Common("Failed to marshal metadata: %v", err)
		return nil
	}

	var metadata FieldMetadata
	if err := json.Unmarshal(jsonBytes, &metadata); err != nil {
		debug.Common("Failed to unmarshal metadata: %v", err)
		return nil
	}

	return &metadata
}

// ExtractTextFieldFromMap extracts a TextField from generic field data
// Handles both jsonValue.value and direct value patterns
func ExtractTextFieldFromMap(fieldData any) *TextField {
	if fieldData == nil {
		return &TextField{}
	}

	fieldMap, ok := fieldData.(map[string]any)
	if !ok {
		// If it's a string, use it directly
		if str, ok := fieldData.(string); ok {
			return &TextField{Value: str}
		}
		return &TextField{}
	}

	field := &TextField{}

	// Try jsonValue.value pattern (standard Sitecore format)
	if jsonValue, ok := fieldMap["jsonValue"].(map[string]any); ok {
		if value, ok := jsonValue["value"].(string); ok {
			field.Value = value
		}
		// Extract metadata from jsonValue (standard pattern in editing mode)
		field.Metadata = extractMetadata(jsonValue)
	} else if value, ok := fieldMap["value"].(string); ok {
		// Fallback: direct value
		field.Value = value
		// Extract metadata from root level (fallback)
		field.Metadata = extractMetadata(fieldMap)
	}

	// Extract editable metadata (contains pre-wrapped HTML with chrome)
	if editable, ok := fieldMap["editable"].(string); ok {
		field.Editable = editable
	}

	return field
}

// ExtractRichTextFieldFromMap extracts a RichTextField from generic field data
// Handles both jsonValue.value and direct value patterns
func ExtractRichTextFieldFromMap(fieldData any) *RichTextField {
	if fieldData == nil {
		return &RichTextField{}
	}

	fieldMap, ok := fieldData.(map[string]any)
	if !ok {
		// If it's a string, use it directly
		if str, ok := fieldData.(string); ok {
			return &RichTextField{Value: str}
		}
		return &RichTextField{}
	}

	field := &RichTextField{}

	// Try jsonValue.value pattern (standard Sitecore format)
	if jsonValue, ok := fieldMap["jsonValue"].(map[string]any); ok {
		if value, ok := jsonValue["value"].(string); ok {
			field.Value = value
		}
		// Extract metadata from jsonValue (standard pattern in editing mode)
		field.Metadata = extractMetadata(jsonValue)
	} else if value, ok := fieldMap["value"].(string); ok {
		// Fallback: direct value
		field.Value = value
		// Extract metadata from root level (fallback)
		field.Metadata = extractMetadata(fieldMap)
	}

	// Extract editable metadata (contains pre-wrapped HTML with chrome)
	if editable, ok := fieldMap["editable"].(string); ok {
		field.Editable = editable
	}

	return field
}

// ExtractImageFieldFromMap extracts an ImageField from generic field data
// Handles both jsonValue.value and direct property patterns
func ExtractImageFieldFromMap(fieldData any) *ImageField {
	if fieldData == nil {
		debug.Common("ExtractImageFieldFromMap fieldData is nil type: %T", fieldData)
		return &ImageField{}
	}

	fieldMap, ok := fieldData.(map[string]any)
	if !ok {
		debug.Common("ExtractImageFieldFromMap fieldMap is not a map, fieldData type: %T", fieldData)
		return &ImageField{}
	}

	field := &ImageField{}

	fieldValues, ok := fieldMap["value"].(map[string]any)
	if !ok {
		debug.Common("ExtractImageFieldFromMap fieldValues[\"value\"] is not a map, fieldData type: %T", fieldData)
		fieldValues = nil
	}

	useJsonValue := false
	if fieldValues == nil {
		fieldValues, ok = fieldMap["jsonValue"].(map[string]any)
		if !ok {
			debug.Common("ExtractImageFieldFromMap fieldValues[\"jsonValue\"] is not a map, fieldData type: %T", fieldData)
			return &ImageField{}
		}
		useJsonValue = ok
	}

	// Try jsonValue.value pattern (standard Sitecore format)
	if jsonValue, ok := fieldValues["jsonValue"].(map[string]any); ok || useJsonValue {
		if useJsonValue {
			jsonValue = fieldValues
		}
		if value, ok := jsonValue["value"].(map[string]any); ok {
			// Extract nested value structure
			field.Value = &ImageFieldValue{}
			if src, ok := value["src"].(string); ok {
				field.Value.Src = src
				field.Src = src // Also set direct property for convenience
			}
			if alt, ok := value["alt"].(string); ok {
				field.Value.Alt = alt
				field.Alt = alt
			}
			if width, ok := value["width"].(string); ok {
				field.Value.Width = width
				field.Width = width
			}
			if height, ok := value["height"].(string); ok {
				field.Value.Height = height
				field.Height = height
			}
		}
	} else {
		// Try direct properties (fallback)
		if src, ok := fieldValues["src"].(string); ok {
			field.Src = src
		}
		if alt, ok := fieldValues["alt"].(string); ok {
			field.Alt = alt
		}
		if width, ok := fieldValues["width"].(string); ok {
			field.Width = width
		}
		if height, ok := fieldValues["height"].(string); ok {
			field.Height = height
		}
	}

	// Extract editable metadata (contains pre-wrapped HTML with chrome)
	if editable, ok := fieldMap["editable"].(string); ok {
		field.Editable = editable
	}

	// Extract field metadata (only present in editing mode)
	// When using jsonValue pattern, metadata is inside jsonValue
	if useJsonValue {
		field.Metadata = extractMetadata(fieldValues)
	} else {
		field.Metadata = extractMetadata(fieldMap)
	}

	return field
}

// ExtractLinkFieldFromMap extracts a LinkField from generic field data
// Handles both jsonValue.value and direct property patterns
func ExtractLinkFieldFromMap(fieldData any) *LinkField {
	if fieldData == nil {
		return &LinkField{}
	}

	fieldMap, ok := fieldData.(map[string]any)
	if !ok {
		return &LinkField{}
	}

	field := &LinkField{}

	// Try jsonValue.value pattern first (standard Sitecore format)
	fieldValues, ok := fieldMap["value"].(map[string]any)
	if !ok {
		fieldValues = nil
	}

	useJsonValue := false
	if fieldValues == nil {
		fieldValues, ok = fieldMap["jsonValue"].(map[string]any)
		if !ok {
			return &LinkField{}
		}
		useJsonValue = ok
	}

	// Extract link value structure
	if jsonValue, ok := fieldValues["jsonValue"].(map[string]any); ok || useJsonValue {
		if useJsonValue {
			jsonValue = fieldValues
		}
		if value, ok := jsonValue["value"].(map[string]any); ok {
			// Extract nested value structure
			field.Value = &LinkFieldValue{}
			if href, ok := value["href"].(string); ok {
				field.Value.Href = href
				field.Href = href // Also set direct property for convenience
			}
			if text, ok := value["text"].(string); ok {
				field.Value.Text = text
				field.Text = text
			}
			if target, ok := value["target"].(string); ok {
				field.Value.Target = target
				field.Target = target
			}
			if title, ok := value["title"].(string); ok {
				field.Value.Title = title
				field.Title = title
			}
			if class, ok := value["class"].(string); ok {
				field.Value.Class = class
				field.Class = class
			}
		}
	} else {
		// Try direct properties (fallback)
		if href, ok := fieldValues["href"].(string); ok {
			field.Href = href
		}
		if text, ok := fieldValues["text"].(string); ok {
			field.Text = text
		}
		if target, ok := fieldValues["target"].(string); ok {
			field.Target = target
		}
		if title, ok := fieldValues["title"].(string); ok {
			field.Title = title
		}
		if class, ok := fieldValues["class"].(string); ok {
			field.Class = class
		}
	}

	// Extract editable metadata (contains pre-wrapped HTML with chrome)
	if editable, ok := fieldMap["editable"].(string); ok {
		field.Editable = editable
	}

	// Extract field metadata (only present in editing mode)
	// When using jsonValue pattern, metadata is inside jsonValue
	if useJsonValue {
		field.Metadata = extractMetadata(fieldValues)
	} else {
		field.Metadata = extractMetadata(fieldMap)
	}

	return field
}

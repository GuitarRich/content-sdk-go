package models

// Field helper functions for convenient field access
// These combine GetFieldByName with typed extractors

// GetTextField extracts and returns a TextField from fields by name
func GetTextField(fields any, fieldName string) *TextField {
	fieldData := GetFieldByName(fields, fieldName)
	return ExtractTextFieldFromMap(fieldData)
}

// GetRichTextField extracts and returns a RichTextField from fields by name
func GetRichTextField(fields any, fieldName string) *RichTextField {
	fieldData := GetFieldByName(fields, fieldName)
	return ExtractRichTextFieldFromMap(fieldData)
}

// GetImageField extracts and returns an ImageField from fields by name
func GetImageField(fields any, fieldName string) *ImageField {
	fieldData := GetFieldByName(fields, fieldName)
	return ExtractImageFieldFromMap(fieldData)
}

// GetLinkField extracts and returns a LinkField from fields by name
func GetLinkField(fields any, fieldName string) *LinkField {
	fieldData := GetFieldByName(fields, fieldName)
	return ExtractLinkFieldFromMap(fieldData)
}

// GetFieldByName extracts a field by name from the fields interface
func GetFieldByName(fields any, name string) any {
	if fields == nil {
		return nil
	}
	fieldsMap, ok := fields.(map[string]any)
	if !ok {
		return nil
	}
	return fieldsMap[name]
}

// IsFieldEmpty checks if a field is empty
// Works with both typed Field interface and generic field data
func IsFieldEmpty(fieldData any) bool {
	if fieldData == nil {
		return true
	}

	// Try as Field interface
	if field, ok := fieldData.(Field); ok {
		return field.IsEmpty()
	}

	// Fallback checks for generic map data
	fieldMap, ok := fieldData.(map[string]any)
	if !ok {
		return true
	}

	// Check for value in jsonValue.value pattern
	if jsonValue, ok := fieldMap["jsonValue"].(map[string]any); ok {
		if value, ok := jsonValue["value"]; ok && value != nil && value != "" {
			return false
		}
	}

	// Check for direct value
	if value, ok := fieldMap["value"]; ok && value != nil && value != "" {
		return false
	}

	return true
}

// FieldHasValue checks if a field has a non-empty value
// This is the inverse of IsFieldEmpty for more readable templates
func FieldHasValue(fieldData any) bool {
	return !IsFieldEmpty(fieldData)
}

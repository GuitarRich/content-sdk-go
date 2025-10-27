package models

import "fmt"

// NotFoundError represents a 404 not found error
type NotFoundError struct {
	Path string
	Site string
}

func (e *NotFoundError) Error() string {
	if e.Site != "" {
		return fmt.Sprintf("page not found: %s (site: %s)", e.Path, e.Site)
	}
	return fmt.Sprintf("page not found: %s", e.Path)
}

// PreviewError represents an error during preview/editing
type PreviewError struct {
	Message string
	ItemID  string
}

func (e *PreviewError) Error() string {
	if e.ItemID != "" {
		return fmt.Sprintf("preview error for item %s: %s", e.ItemID, e.Message)
	}
	return fmt.Sprintf("preview error: %s", e.Message)
}

// GraphQLError represents an error from a GraphQL request
type GraphQLError struct {
	Message    string
	Path       []interface{}
	Extensions map[string]interface{}
}

func (e *GraphQLError) Error() string {
	if len(e.Path) > 0 {
		return fmt.Sprintf("GraphQL error at %v: %s", e.Path, e.Message)
	}
	return fmt.Sprintf("GraphQL error: %s", e.Message)
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s: %s", e.Field, e.Message)
}

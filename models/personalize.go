package models

// PersonalizeInfo contains personalization variant information
type PersonalizeInfo struct {
	// VariantIds are the selected variant IDs for personalization
	VariantIds []string `json:"variantIds"`

	// Scope is the CDP scope
	Scope string `json:"scope,omitempty"`
}

// ExperienceParams contains parameters for CDP experience tracking
type ExperienceParams struct {
	// Referrer is the referring URL
	Referrer string `json:"referrer,omitempty"`

	// UTM contains UTM tracking parameters
	UTM UTMParams `json:"utm"`

	// PageVariantID is the page variant ID for A/B testing
	PageVariantID string `json:"pageVariantId,omitempty"`
}

// UTMParams contains UTM tracking parameters
type UTMParams struct {
	// Campaign is the marketing campaign name
	Campaign string `json:"campaign,omitempty"`

	// Source is the traffic source
	Source string `json:"source,omitempty"`

	// Medium is the marketing medium
	Medium string `json:"medium,omitempty"`

	// Content is the content identifier
	Content string `json:"content,omitempty"`

	// Term is the paid keyword
	Term string `json:"term,omitempty"`
}

// PersonalizeRewriteData contains personalization path rewrite information
type PersonalizeRewriteData struct {
	// VariantId is the variant ID from the path
	VariantId string `json:"variantId"`

	// NormalizedPath is the path with variant prefix removed
	NormalizedPath string `json:"normalizedPath"`
}

// PersonalizeExecution represents a CDP personalization execution result
type PersonalizeExecution struct {
	// FriendlyId is the friendly ID of the personalization experience
	FriendlyId string `json:"friendlyId"`

	// VariantIds are the selected variant IDs
	VariantIds []string `json:"variantIds"`

	// ExperienceID is the CDP experience ID
	ExperienceID string `json:"experienceId,omitempty"`
}

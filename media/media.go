package media

import (
	"fmt"
	"net/url"
	"strings"
)

// ImageParams contains parameters for image transformation
type ImageParams struct {
	// Width sets the image width
	Width *int `json:"w,omitempty"`

	// Height sets the image height
	Height *int `json:"h,omitempty"`

	// MaxWidth sets the maximum width
	MaxWidth *int `json:"mw,omitempty"`

	// MaxHeight sets the maximum height
	MaxHeight *int `json:"mh,omitempty"`

	// Quality sets the image quality (0-100)
	Quality *int `json:"q,omitempty"`

	// Scale sets the scale factor
	Scale *float64 `json:"scale,omitempty"`

	// AllowStretch allows image stretching
	AllowStretch *bool `json:"as,omitempty"`

	// IgnoreAspectRatio ignores the aspect ratio
	IgnoreAspectRatio *bool `json:"iar,omitempty"`

	// Thumbnail generates a thumbnail
	Thumbnail *bool `json:"thumbnail,omitempty"`

	// BackgroundColor sets the background color (hex)
	BackgroundColor *string `json:"bc,omitempty"`

	// Database specifies the database
	Database *string `json:"db,omitempty"`

	// Language specifies the language
	Language *string `json:"la,omitempty"`

	// Version specifies the version
	Version *string `json:"vs,omitempty"`
}

// ImageField represents a Sitecore image field
type ImageField struct {
	Value *ImageFieldValue `json:"value,omitempty"`
}

// ImageFieldValue contains image field data
type ImageFieldValue struct {
	Src    string `json:"src"`
	Alt    string `json:"alt,omitempty"`
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

// MediaAPI provides functions for working with Sitecore media
type MediaAPI struct {
	mediaServerURL string
}

// NewMediaAPI creates a new media API instance
func NewMediaAPI(mediaServerURL string) *MediaAPI {
	return &MediaAPI{
		mediaServerURL: strings.TrimSuffix(mediaServerURL, "/"),
	}
}

// GetImageURL generates an image URL with optional parameters
func (m *MediaAPI) GetImageURL(imageField *ImageField, params *ImageParams) string {
	if imageField == nil || imageField.Value == nil || imageField.Value.Src == "" {
		return ""
	}

	src := imageField.Value.Src

	// If src is already a full URL, use it as base
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		return m.buildURL(src, params)
	}

	// If src is a relative path, prepend media server URL
	if strings.HasPrefix(src, "/") {
		src = m.mediaServerURL + src
	} else {
		src = m.mediaServerURL + "/" + src
	}

	return m.buildURL(src, params)
}

// buildURL builds the final URL with query parameters
func (m *MediaAPI) buildURL(baseURL string, params *ImageParams) string {
	if params == nil {
		return baseURL
	}

	// Parse the base URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	// Get existing query parameters
	q := u.Query()

	// Add image transformation parameters
	if params.Width != nil {
		q.Set("w", fmt.Sprintf("%d", *params.Width))
	}

	if params.Height != nil {
		q.Set("h", fmt.Sprintf("%d", *params.Height))
	}

	if params.MaxWidth != nil {
		q.Set("mw", fmt.Sprintf("%d", *params.MaxWidth))
	}

	if params.MaxHeight != nil {
		q.Set("mh", fmt.Sprintf("%d", *params.MaxHeight))
	}

	if params.Quality != nil {
		q.Set("q", fmt.Sprintf("%d", *params.Quality))
	}

	if params.Scale != nil {
		q.Set("scale", fmt.Sprintf("%.2f", *params.Scale))
	}

	if params.AllowStretch != nil && *params.AllowStretch {
		q.Set("as", "1")
	}

	if params.IgnoreAspectRatio != nil && *params.IgnoreAspectRatio {
		q.Set("iar", "1")
	}

	if params.Thumbnail != nil && *params.Thumbnail {
		q.Set("thumbnail", "1")
	}

	if params.BackgroundColor != nil {
		q.Set("bc", *params.BackgroundColor)
	}

	if params.Database != nil {
		q.Set("db", *params.Database)
	}

	if params.Language != nil {
		q.Set("la", *params.Language)
	}

	if params.Version != nil {
		q.Set("vs", *params.Version)
	}

	// Set the query string
	u.RawQuery = q.Encode()

	return u.String()
}

// GetResponsiveImageURL generates responsive image URLs
func (m *MediaAPI) GetResponsiveImageURL(imageField *ImageField, widths []int) map[int]string {
	result := make(map[int]string)

	for _, width := range widths {
		w := width
		params := &ImageParams{
			Width: &w,
		}
		result[width] = m.GetImageURL(imageField, params)
	}

	return result
}

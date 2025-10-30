package models

// Field represents a generic Sitecore field
// All field types implement this interface
type Field interface {
	GetValue() any
	GetEditable() string
	IsEmpty() bool
}

// TextField represents a simple text field value (Single-Line Text, Multi-Line Text)
type TextField struct {
	Value    string `json:"value"`
	Editable string `json:"editable,omitempty"`
}

func (f *TextField) GetValue() any {
	return f.Value
}

func (f *TextField) GetEditable() string {
	return f.Editable
}

func (f *TextField) IsEmpty() bool {
	return f.Value == ""
}

// RichTextField represents a rich text field with HTML content
type RichTextField struct {
	Value    string `json:"value"`
	Editable string `json:"editable,omitempty"`
}

func (f *RichTextField) GetValue() any {
	return f.Value
}

func (f *RichTextField) GetEditable() string {
	return f.Editable
}

func (f *RichTextField) IsEmpty() bool {
	return f.Value == ""
}

// ImageField represents an image field
type ImageField struct {
	Src      string           `json:"src"`
	Alt      string           `json:"alt"`
	Width    string           `json:"width,omitempty"`
	Height   string           `json:"height,omitempty"`
	Editable string           `json:"editable,omitempty"`
	Value    *ImageFieldValue `json:"value,omitempty"`
}

// ImageFieldValue contains the nested image value structure
type ImageFieldValue struct {
	Src    string `json:"src"`
	Alt    string `json:"alt"`
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

func (f *ImageField) GetValue() any {
	return f
}

func (f *ImageField) GetEditable() string {
	return f.Editable
}

func (f *ImageField) IsEmpty() bool {
	return f.Src == "" && (f.Value == nil || f.Value.Src == "")
}

// GetSrc returns the image source URL, checking both direct and nested value
func (f *ImageField) GetSrc() string {
	if f.Src != "" {
		return f.Src
	}
	if f.Value != nil && f.Value.Src != "" {
		return f.Value.Src
	}
	return ""
}

// GetAlt returns the image alt text, checking both direct and nested value
func (f *ImageField) GetAlt() string {
	if f.Alt != "" {
		return f.Alt
	}
	if f.Value != nil && f.Value.Alt != "" {
		return f.Value.Alt
	}
	return ""
}

// GetWidth returns the image width, checking both direct and nested value
func (f *ImageField) GetWidth() string {
	if f.Width != "" {
		return f.Width
	}
	if f.Value != nil && f.Value.Width != "" {
		return f.Value.Width
	}
	return ""
}

// GetHeight returns the image height, checking both direct and nested value
func (f *ImageField) GetHeight() string {
	if f.Height != "" {
		return f.Height
	}
	if f.Value != nil && f.Value.Height != "" {
		return f.Value.Height
	}
	return ""
}

// LinkField represents a link field (General Link, Internal Link, External Link)
type LinkField struct {
	Href     string          `json:"href"`
	Text     string          `json:"text"`
	Target   string          `json:"target,omitempty"`
	Title    string          `json:"title,omitempty"`
	Class    string          `json:"class,omitempty"`
	Editable string          `json:"editable,omitempty"`
	Value    *LinkFieldValue `json:"value,omitempty"`
}

// LinkFieldValue contains the nested link value structure
type LinkFieldValue struct {
	Href   string `json:"href"`
	Text   string `json:"text"`
	Target string `json:"target,omitempty"`
	Title  string `json:"title,omitempty"`
	Class  string `json:"class,omitempty"`
}

func (f *LinkField) GetValue() any {
	return f
}

func (f *LinkField) GetEditable() string {
	return f.Editable
}

func (f *LinkField) IsEmpty() bool {
	return f.Href == "" && (f.Value == nil || f.Value.Href == "")
}

// GetHref returns the link href, checking both direct and nested value
func (f *LinkField) GetHref() string {
	if f.Href != "" {
		return f.Href
	}
	if f.Value != nil && f.Value.Href != "" {
		return f.Value.Href
	}
	return ""
}

// GetText returns the link text, checking both direct and nested value
func (f *LinkField) GetText() string {
	if f.Text != "" {
		return f.Text
	}
	if f.Value != nil && f.Value.Text != "" {
		return f.Value.Text
	}
	return ""
}

// GetTarget returns the link target, checking both direct and nested value
func (f *LinkField) GetTarget() string {
	if f.Target != "" {
		return f.Target
	}
	if f.Value != nil && f.Value.Target != "" {
		return f.Value.Target
	}
	return ""
}

// GetTitle returns the link title, checking both direct and nested value
func (f *LinkField) GetTitle() string {
	if f.Title != "" {
		return f.Title
	}
	if f.Value != nil && f.Value.Title != "" {
		return f.Value.Title
	}
	return ""
}

// GetClass returns the link class, checking both direct and nested value
func (f *LinkField) GetClass() string {
	if f.Class != "" {
		return f.Class
	}
	if f.Value != nil && f.Value.Class != "" {
		return f.Value.Class
	}
	return ""
}

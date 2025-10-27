package media

import (
	"strings"
	"testing"
)

func TestMediaAPI_GetImageURL_Simple(t *testing.T) {
	api := NewMediaAPI("https://media.example.com")

	imageField := &ImageField{
		Value: &ImageFieldValue{
			Src: "/media/image.jpg",
			Alt: "Test Image",
		},
	}

	url := api.GetImageURL(imageField, nil)

	expected := "https://media.example.com/media/image.jpg"
	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}

func TestMediaAPI_GetImageURL_WithParams(t *testing.T) {
	api := NewMediaAPI("https://media.example.com")

	width := 800
	height := 600
	quality := 85

	imageField := &ImageField{
		Value: &ImageFieldValue{
			Src: "/media/image.jpg",
		},
	}

	params := &ImageParams{
		Width:   &width,
		Height:  &height,
		Quality: &quality,
	}

	url := api.GetImageURL(imageField, params)

	if !strings.Contains(url, "w=800") {
		t.Error("URL should contain width parameter")
	}

	if !strings.Contains(url, "h=600") {
		t.Error("URL should contain height parameter")
	}

	if !strings.Contains(url, "q=85") {
		t.Error("URL should contain quality parameter")
	}
}

func TestMediaAPI_GetImageURL_FullURL(t *testing.T) {
	api := NewMediaAPI("https://media.example.com")

	imageField := &ImageField{
		Value: &ImageFieldValue{
			Src: "https://cdn.example.com/image.jpg",
		},
	}

	url := api.GetImageURL(imageField, nil)

	expected := "https://cdn.example.com/image.jpg"
	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}

func TestMediaAPI_GetImageURL_EmptyField(t *testing.T) {
	api := NewMediaAPI("https://media.example.com")

	url := api.GetImageURL(nil, nil)

	if url != "" {
		t.Errorf("expected empty string, got %s", url)
	}
}

func TestMediaAPI_GetImageURL_EmptyValue(t *testing.T) {
	api := NewMediaAPI("https://media.example.com")

	imageField := &ImageField{
		Value: &ImageFieldValue{
			Src: "",
		},
	}

	url := api.GetImageURL(imageField, nil)

	if url != "" {
		t.Errorf("expected empty string, got %s", url)
	}
}

func TestMediaAPI_GetResponsiveImageURL(t *testing.T) {
	api := NewMediaAPI("https://media.example.com")

	imageField := &ImageField{
		Value: &ImageFieldValue{
			Src: "/media/image.jpg",
		},
	}

	widths := []int{320, 640, 1024, 1920}
	urls := api.GetResponsiveImageURL(imageField, widths)

	if len(urls) != 4 {
		t.Errorf("expected 4 URLs, got %d", len(urls))
	}

	for _, width := range widths {
		url, exists := urls[width]
		if !exists {
			t.Errorf("missing URL for width %d", width)
		}

		// Just check that URL contains a width parameter
		if !strings.Contains(url, "w=") {
			t.Errorf("URL for width %d should contain width parameter", width)
		}
	}
}

func TestMediaAPI_BuildURL_AllParams(t *testing.T) {
	api := NewMediaAPI("https://media.example.com")

	width := 800
	height := 600
	quality := 90
	scale := 2.0
	allowStretch := true
	ignoreAspect := true
	thumbnail := true
	bgColor := "FFFFFF"
	db := "web"
	lang := "en"
	version := "1"

	params := &ImageParams{
		Width:             &width,
		Height:            &height,
		Quality:           &quality,
		Scale:             &scale,
		AllowStretch:      &allowStretch,
		IgnoreAspectRatio: &ignoreAspect,
		Thumbnail:         &thumbnail,
		BackgroundColor:   &bgColor,
		Database:          &db,
		Language:          &lang,
		Version:           &version,
	}

	imageField := &ImageField{
		Value: &ImageFieldValue{
			Src: "/media/image.jpg",
		},
	}

	url := api.GetImageURL(imageField, params)

	// Verify all parameters are in the URL
	requiredParams := []string{
		"w=800",
		"h=600",
		"q=90",
		"scale=2.00",
		"as=1",
		"iar=1",
		"thumbnail=1",
		"bc=FFFFFF",
		"db=web",
		"la=en",
		"vs=1",
	}

	for _, param := range requiredParams {
		if !strings.Contains(url, param) {
			t.Errorf("URL should contain parameter: %s", param)
		}
	}
}

package seo

import (
	"encoding/json"
	"html/template"
)

// Meta represents SEO metadata for a page
type Meta struct {
	Title       string
	Description string
	Keywords    string
	OGImage     string
	Canonical   string
	Type        string // website, article, product
	JsonLd      map[string]interface{}
}

// NewMeta creates a new Meta instance
func NewMeta(title, description string) *Meta {
	return &Meta{
		Title:       title,
		Description: description,
		Type:        "website",
	}
}

// WithKeywords sets the keywords
func (m *Meta) WithKeywords(keywords string) *Meta {
	m.Keywords = keywords
	return m
}

// WithOGImage sets the Open Graph image
func (m *Meta) WithOGImage(image string) *Meta {
	m.OGImage = image
	return m
}

// WithCanonical sets the canonical URL
func (m *Meta) WithCanonical(url string) *Meta {
	m.Canonical = url
	return m
}

// WithJsonLd sets the JSON-LD structured data
func (m *Meta) WithJsonLd(data map[string]interface{}) *Meta {
	m.JsonLd = data
	return m
}

// RenderJsonLd generates JSON-LD script tag
func (m *Meta) RenderJsonLd() template.HTML {
	if m.JsonLd == nil {
		return ""
	}

	jsonData, err := json.Marshal(m.JsonLd)
	if err != nil {
		return ""
	}

	return template.HTML(`<script type="application/ld+json">` + string(jsonData) + `</script>`)
}

// GenerateWebsiteSchema creates a basic website schema
func GenerateWebsiteSchema(name, url, description string) map[string]interface{} {
	return map[string]interface{}{
		"@context":    "https://schema.org",
		"@type":       "WebSite",
		"name":        name,
		"url":         url,
		"description": description,
	}
}

// GenerateOrganizationSchema creates an organization schema
func GenerateOrganizationSchema(name, url, logo string) map[string]interface{} {
	return map[string]interface{}{
		"@context": "https://schema.org",
		"@type":    "Organization",
		"name":     name,
		"url":      url,
		"logo":     logo,
	}
}

// GenerateBreadcrumbSchema creates breadcrumb navigation schema
func GenerateBreadcrumbSchema(items []BreadcrumbItem) map[string]interface{} {
	listItems := make([]map[string]interface{}, len(items))
	for i, item := range items {
		listItems[i] = map[string]interface{}{
			"@type":    "ListItem",
			"position": i + 1,
			"name":     item.Name,
			"item":     item.URL,
		}
	}

	return map[string]interface{}{
		"@context":        "https://schema.org",
		"@type":           "BreadcrumbList",
		"itemListElement": listItems,
	}
}

// BreadcrumbItem represents a single breadcrumb item
type BreadcrumbItem struct {
	Name string
	URL  string
}

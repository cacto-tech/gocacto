package sitemap

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"cacto-cms/app/domain/page"
)

// URLSet represents the root element of a sitemap
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

// URL represents a single URL in the sitemap
type URL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq,omitempty"`
	Priority   float64 `xml:"priority,omitempty"`
}

// Generator handles sitemap generation
type Generator struct {
	baseURL    string
	outputPath string
	repo       page.Repository
}

// NewGenerator creates a new sitemap generator
func NewGenerator(baseURL, outputPath string, repo page.Repository) *Generator {
	return &Generator{
		baseURL:    baseURL,
		outputPath: outputPath,
		repo:       repo,
	}
}

// Generate generates the sitemap XML file
func (g *Generator) Generate() error {
	pages, err := g.repo.FindPublished()
	if err != nil {
		return fmt.Errorf("failed to fetch pages: %w", err)
	}

	urlset := URLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]URL, 0, len(pages)+1),
	}

	// Add homepage
	urlset.URLs = append(urlset.URLs, URL{
		Loc:        g.baseURL,
		LastMod:    time.Now().Format("2006-01-02"),
		ChangeFreq: "daily",
		Priority:   1.0,
	})

	// Add pages
	for _, p := range pages {
		urlset.URLs = append(urlset.URLs, URL{
			Loc:        fmt.Sprintf("%s/%s", g.baseURL, p.Slug),
			LastMod:    p.UpdatedAt.Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   0.8,
		})
	}

	// Generate XML
	output, err := xml.MarshalIndent(urlset, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal XML: %w", err)
	}

	// Write to file
	xmlContent := []byte(xml.Header + string(output))
	if err := os.WriteFile(g.outputPath, xmlContent, 0644); err != nil {
		return fmt.Errorf("failed to write sitemap: %w", err)
	}

	return nil
}

// ScheduleDaily runs sitemap generation every 24 hours
func (g *Generator) ScheduleDaily() {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		// Generate immediately on start
		if err := g.Generate(); err != nil {
			fmt.Printf("Sitemap generation failed: %v\n", err)
		}

		// Then generate daily
		for range ticker.C {
			if err := g.Generate(); err != nil {
				fmt.Printf("Sitemap generation failed: %v\n", err)
			} else {
				fmt.Println("✓ Sitemap regenerated")
			}
		}
	}()
}

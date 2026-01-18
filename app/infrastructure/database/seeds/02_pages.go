package seeds

import (
	"cacto-cms/app/domain/page"
)

// Pages returns seed data for pages
func Pages() []*PageSeed {
	return []*PageSeed{
		{
			Page: &page.Page{
				Slug:            "",
				Title:           "Home",
				Content:         "",
				MetaTitle:       "Cacto CMS - Performance-Focused Website",
				MetaDescription: "High-performance enterprise CMS built with Go and Templ",
				MetaKeywords:    "go, cms, performance, web",
				Status:          page.StatusPublished,
			},
			ComponentNames: []string{"home-hero", "home-about"},
		},
		{
			Page: &page.Page{
				Slug:            "about",
				Title:           "About",
				Content:         "About page content",
				MetaTitle:       "About - Cacto CMS",
				MetaDescription: "Learn about Cacto CMS",
				MetaKeywords:    "about, cms, go, performance",
				Status:          page.StatusPublished,
			},
			ComponentNames: []string{"about-hero", "about-intro"},
		},
		{
			Page: &page.Page{
				Slug:            "contact",
				Title:           "Contact",
				Content:         "Contact page content",
				MetaTitle:       "Contact - Cacto CMS",
				MetaDescription: "Contact Cacto CMS",
				MetaKeywords:    "contact, support",
				Status:          page.StatusPublished,
			},
			ComponentNames: []string{"contact-intro"},
		},
	}
}

// PageSeed represents a page with its component associations
type PageSeed struct {
	Page            *page.Page
	ComponentNames  []string
}

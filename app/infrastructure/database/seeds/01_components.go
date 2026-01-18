package seeds

import (
	"cacto-cms/app/domain/component"
)

// Components returns seed data for components
func Components() []*component.Component {
	return []*component.Component{
		{
			Type:     component.TypeHero,
			Name:     "home-hero",
			Title:    "Welcome to Cacto CMS",
			Subtitle: "Performance-focused, modern and scalable CMS",
			LinkText: "Explore",
			LinkURL:  "/about",
		},
		{
			Type:    component.TypeAbout,
			Name:    "home-about",
			Title:   "About Us",
			Content: "High-performance CMS system built with Go + Templ + HTMX. Fast and reliable content management with modern web technologies.",
		},
		{
			Type:    component.TypeText,
			Name:    "about-intro",
			Content: "<p>Cacto CMS is a high-performance content management system built with the Go programming language. It provides a modern web experience with the Templ template engine and HTMX.</p><p>It stands out with enterprise-level scalable architecture, clean code structure, and high performance.</p>",
		},
		{
			Type:     component.TypeHero,
			Name:     "about-hero",
			Title:    "About Us",
			Subtitle: "CMS solution built with modern technologies",
			LinkText: "Contact Us",
			LinkURL:  "/contact",
		},
		{
			Type:    component.TypeText,
			Name:    "contact-intro",
			Content: "<p>You can use the form below to contact us. Don't hesitate to reach out to us with your questions.</p>",
		},
	}
}

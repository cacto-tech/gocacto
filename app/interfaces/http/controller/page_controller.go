package controller

import (
	"net/http"

	componentservice "cacto-cms/app/application/component"
	"cacto-cms/app/application/page"
	"cacto-cms/app/domain/component"
	"cacto-cms/app/interfaces/templates/pages"
	componentrenderer "cacto-cms/app/shared/component"
	"cacto-cms/app/shared/sanitize"
	"cacto-cms/app/shared/seo"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

// PageController handles page-related HTTP requests
type PageController struct {
	*BaseController
	pageService     *page.Service
	componentService *componentservice.Service
	componentRenderer *componentrenderer.Renderer
	seoManager      *seo.Manager
}

// NewPageController creates a new page controller
func NewPageController(
	baseURL string,
	pageService *page.Service,
	componentService *componentservice.Service,
	seoManager *seo.Manager,
) *PageController {
	return &PageController{
		BaseController:   NewBaseController(baseURL),
		pageService:      pageService,
		componentService: componentService,
		componentRenderer: componentrenderer.NewRenderer(),
		seoManager:       seoManager,
	}
}

// ShowHome renders the home page
func (c *PageController) ShowHome(w http.ResponseWriter, r *http.Request) {
	// Get home page from database (slug is empty string)
	p, err := c.pageService.GetPageBySlug("")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Get page components
	var renderedComponents []templ.Component
	if len(p.Components) > 0 {
		// Convert page.Component to domain.Component
		domainComponents := make([]*component.Component, len(p.Components))
		for i, comp := range p.Components {
			domainComponents[i] = &component.Component{
				ID:       comp.ID,
				Type:     component.Type(comp.Type),
				Name:     comp.Name,
				Title:    comp.Title,
				Subtitle: comp.Subtitle,
				Content:  comp.Content,
				ImageURL: comp.ImageURL,
				LinkURL:  comp.LinkURL,
				LinkText: comp.LinkText,
				DataJSON: comp.DataJSON,
			}
			// Merge with defaults
			domainComponents[i].MergeWithDefaults()
		}

		renderedComponents, err = c.componentRenderer.RenderMultiple(domainComponents)
		if err != nil {
			http.Error(w, "Failed to render components", http.StatusInternalServerError)
			return
		}
	} else {
		// Fallback to simple content (sanitize content)
		sanitizedContent := sanitize.HTML(p.Content)
		renderedComponents = []templ.Component{
			pages.SimpleContent(p.Title, sanitizedContent),
		}
	}

	// Create page wrapper component
	pageComponent := pages.PageWithComponents(renderedComponents)

	// SEO - use page meta if available, otherwise use defaults
	meta := c.seoManager.ForPageWithDefaults(
		p.MetaTitle,
		p.MetaDescription,
		p.MetaKeywords,
		p.OGImage,
		"",
	)

	// Render
	c.Render(w, r, meta, pageComponent)
}

// ShowPage renders a page by slug
func (c *PageController) ShowPage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	p, err := c.pageService.GetPageBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Get page components
	var renderedComponents []templ.Component
	if len(p.Components) > 0 {
		// Convert page.Component to domain.Component
		domainComponents := make([]*component.Component, len(p.Components))
		for i, comp := range p.Components {
			domainComponents[i] = &component.Component{
				ID:       comp.ID,
				Type:     component.Type(comp.Type),
				Name:     comp.Name,
				Title:    comp.Title,
				Subtitle: comp.Subtitle,
				Content:  comp.Content,
				ImageURL: comp.ImageURL,
				LinkURL:  comp.LinkURL,
				LinkText: comp.LinkText,
				DataJSON: comp.DataJSON,
			}
			// Merge with defaults
			domainComponents[i].MergeWithDefaults()
		}

		renderedComponents, err = c.componentRenderer.RenderMultiple(domainComponents)
		if err != nil {
			http.Error(w, "Failed to render components", http.StatusInternalServerError)
			return
		}
	} else {
		// Fallback to simple content (sanitize content)
		sanitizedContent := sanitize.HTML(p.Content)
		renderedComponents = []templ.Component{
			pages.SimpleContent(p.Title, sanitizedContent),
		}
	}

	// Create page wrapper component
	pageComponent := pages.PageWithComponents(renderedComponents)

	// SEO
	meta := c.seoManager.ForPageWithDefaults(
		p.MetaTitle,
		p.MetaDescription,
		p.MetaKeywords,
		p.OGImage,
		p.Slug,
	)

	// Render
	c.Render(w, r, meta, pageComponent)
}


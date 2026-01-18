package controller

import (
	"net/http"

	"cacto-cms/app/shared/seo"
	"cacto-cms/app/interfaces/templates/layouts"

	"github.com/a-h/templ"
)

// BaseController provides common functionality for all controllers
type BaseController struct {
	baseURL string
}

// NewBaseController creates a new base controller
func NewBaseController(baseURL string) *BaseController {
	return &BaseController{
		baseURL: baseURL,
	}
}

// Render renders a templ component with layout
func (c *BaseController) Render(w http.ResponseWriter, r *http.Request, meta *seo.Meta, content templ.Component) {
	layouts.Base(*meta, content).Render(r.Context(), w)
}

// RenderHTML renders raw HTML (for simple pages)
func (c *BaseController) RenderHTML(w http.ResponseWriter, r *http.Request, meta *seo.Meta, htmlContent string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	// Simple HTML rendering (templ is used for main pages)
	html := `<!DOCTYPE html>
<html lang="tr">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>` + meta.Title + `</title>
	<meta name="description" content="` + meta.Description + `">
	<link rel="stylesheet" href="/static/css/output.css">
</head>
<body>
	<header class="header">
		<nav class="container">
			<div class="logo"><a href="/">Cacto CMS</a></div>
		</nav>
	</header>
	<main class="container" style="padding: 4rem 0;">
		` + htmlContent + `
	</main>
</body>
</html>`
	
	w.Write([]byte(html))
}

// BaseURL returns the base URL
func (c *BaseController) BaseURL() string {
	return c.baseURL
}

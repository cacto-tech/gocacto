package component

import (
	"fmt"
	"cacto-cms/app/domain/component"
	"cacto-cms/app/interfaces/templates/components"
	"cacto-cms/app/shared/sanitize"
	"github.com/a-h/templ"
)

// Component is an alias to avoid naming conflict
type Component = component.Component

// Renderer handles component rendering
type Renderer struct {
	registry map[component.Type]RenderFunc
}

// RenderFunc is a function that renders a component
type RenderFunc func(c *component.Component) (templ.Component, error)

// NewRenderer creates a new component renderer
func NewRenderer() *Renderer {
	r := &Renderer{
		registry: make(map[component.Type]RenderFunc),
	}
	
	// Register default component renderers
	r.Register(component.TypeHero, r.renderHero)
	r.Register(component.TypeAbout, r.renderAbout)
	r.Register(component.TypeText, r.renderText)
	
	return r
}

// Register registers a component renderer
func (r *Renderer) Register(componentType component.Type, renderFunc RenderFunc) {
	r.registry[componentType] = renderFunc
}

// Render renders a component
func (r *Renderer) Render(c *Component) (templ.Component, error) {
	renderFunc, exists := r.registry[c.Type]
	if !exists {
		return nil, fmt.Errorf("no renderer registered for component type: %s", c.Type)
	}
	
	return renderFunc(c)
}

// RenderMultiple renders multiple components in order
func (r *Renderer) RenderMultiple(components []*Component) ([]templ.Component, error) {
	rendered := make([]templ.Component, 0, len(components))
	
	for _, c := range components {
		// Merge with defaults before rendering
		c.MergeWithDefaults()
		
		comp, err := r.Render(c)
		if err != nil {
			return nil, fmt.Errorf("failed to render component %s: %w", c.Name, err)
		}
		rendered = append(rendered, comp)
	}
	
	return rendered, nil
}

// renderHero renders a hero component
func (r *Renderer) renderHero(c *Component) (templ.Component, error) {
	data := components.HeroData{
		Title:      c.Title,
		Subtitle:   c.Subtitle,
		ButtonText: c.LinkText,
		ButtonLink: c.LinkURL,
	}
	return components.Hero(data), nil
}

// renderAbout renders an about component
func (r *Renderer) renderAbout(c *Component) (templ.Component, error) {
	data := components.AboutData{
		Title:   c.Title,
		Content: sanitize.HTML(c.Content), // Sanitize user content
	}
	return components.About(data), nil
}

// renderText renders a text component
func (r *Renderer) renderText(c *Component) (templ.Component, error) {
	data := components.TextData{
		Content: sanitize.HTML(c.Content), // Sanitize user content
	}
	return components.Text(data), nil
}

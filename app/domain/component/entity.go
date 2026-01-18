package component

// Type represents component types
type Type string

const (
	TypeHero   Type = "hero"
	TypeAbout  Type = "about"
	TypeText   Type = "text"
	TypeImage  Type = "image"
	TypeCTA    Type = "cta"
	TypeGrid   Type = "grid"
	TypeList   Type = "list"
)

// Component represents a reusable component entity
type Component struct {
	ID       int    `json:"id"`
	Type     Type   `json:"type"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
	ImageURL string `json:"image_url"`
	LinkURL  string `json:"link_url"`
	LinkText string `json:"link_text"`
	DataJSON string `json:"data_json"`
}

// Defaults returns default values for a component type
func (c *Component) Defaults() *Component {
	switch c.Type {
	case TypeHero:
		return &Component{
			Type:     TypeHero,
			Title:    "Welcome",
			Subtitle: "Modern and performant solutions",
			LinkText: "Explore",
			LinkURL:  "/",
		}
	case TypeAbout:
		return &Component{
			Type:    TypeAbout,
			Title:   "About Us",
			Content: "About us content will appear here.",
		}
	case TypeText:
		return &Component{
			Type:    TypeText,
			Content: "Text content will appear here.",
		}
	default:
		return c
	}
}

// MergeWithDefaults merges component with its defaults
func (c *Component) MergeWithDefaults() *Component {
	defaults := c.Defaults()
	
	if c.Title == "" {
		c.Title = defaults.Title
	}
	if c.Subtitle == "" {
		c.Subtitle = defaults.Subtitle
	}
	if c.Content == "" {
		c.Content = defaults.Content
	}
	if c.LinkText == "" {
		c.LinkText = defaults.LinkText
	}
	if c.LinkURL == "" {
		c.LinkURL = defaults.LinkURL
	}
	
	return c
}

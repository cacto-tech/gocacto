package seo

// Manager provides SEO management utilities
type Manager struct {
	baseURL string
	siteName string
	siteDescription string
}

// NewManager creates a new SEO manager
func NewManager(baseURL, siteName, siteDescription string) *Manager {
	return &Manager{
		baseURL: baseURL,
		siteName: siteName,
		siteDescription: siteDescription,
	}
}

// ForPage creates SEO meta for a page
func (m *Manager) ForPage(title, description, keywords, ogImage, slug string) *Meta {
	meta := NewMeta(title, description)
	
	if keywords != "" {
		meta = meta.WithKeywords(keywords)
	}
	
	if ogImage != "" {
		meta = meta.WithOGImage(ogImage)
	}
	
	canonical := m.baseURL
	if slug != "" && slug != "/" {
		canonical = m.baseURL + "/" + slug
	}
	meta = meta.WithCanonical(canonical)
	
	return meta
}

// ForHomePage creates SEO meta for home page
func (m *Manager) ForHomePage() *Meta {
	title := m.siteName
	if m.siteName == "" {
		title = "Cacto CMS - Performans Odaklı Web Sitesi"
	}
	
	description := m.siteDescription
	if m.siteDescription == "" {
		description = "Go ve Templ ile geliştirilmiş yüksek performanslı kurumsal CMS"
	}
	
	meta := NewMeta(title, description).
		WithKeywords("go, cms, performans, web").
		WithCanonical(m.baseURL).
		WithJsonLd(GenerateWebsiteSchema(
			m.siteName,
			m.baseURL,
			m.siteDescription,
		))
	
	return meta
}

// ForPageWithDefaults creates SEO meta with defaults if values are empty
func (m *Manager) ForPageWithDefaults(pageTitle, pageDescription, pageKeywords, pageOGImage, slug string) *Meta {
	title := pageTitle
	if title == "" {
		title = m.siteName
	}
	
	description := pageDescription
	if description == "" {
		description = m.siteDescription
	}
	
	return m.ForPage(title, description, pageKeywords, pageOGImage, slug)
}

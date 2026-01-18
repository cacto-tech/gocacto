package sanitize

import (
	"github.com/microcosm-cc/bluemonday"
)

var (
	// UGC is a policy for user-generated content (allows most HTML but sanitizes dangerous content)
	UGC = bluemonday.UGCPolicy()
	
	// Strict is a strict policy that only allows text (no HTML)
	Strict = bluemonday.StrictPolicy()
)

// HTML sanitizes HTML content using UGC policy
func HTML(html string) string {
	return UGC.Sanitize(html)
}

// StrictHTML sanitizes HTML content using strict policy (removes all HTML)
func StrictHTML(html string) string {
	return Strict.Sanitize(html)
}

// ConfigureUGCPolicy allows customization of UGC policy
func ConfigureUGCPolicy() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	
	// Allow additional attributes if needed
	// p.AllowAttrs("class").OnElements("div", "span")
	
	return p
}

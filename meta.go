package seo

import (
	"encoding/json"
	"html"
	"strings"
)

// Meta describes one page's SEO. Render it into the document <head> on the server
// (SSR) or mirror it on the client with the frontend helper.
type Meta struct {
	Title       string
	Description string
	Canonical   string // absolute or root-relative; SiteURL is prepended
	Image       string // OG/Twitter image
	Type        string // og:type, default "website"
	SiteName    string
	JSONLD      any // marshaled into a ld+json script (schema.org)
}

// HTML renders the full <head> SEO block: title, description, canonical, OpenGraph,
// Twitter cards, JSON-LD, the search-engine verification metas, and the Google
// Analytics snippet (when SEO_GA_ID is set). Returns ready-to-inject HTML.
func (m Meta) HTML() string {
	c := Settings()
	b := &strings.Builder{}
	w := func(s string) { b.WriteString(s); b.WriteByte('\n') }

	ogType := m.Type
	if ogType == "" {
		ogType = "website"
	}
	canon := abs(c.SiteURL, m.Canonical)
	img := m.Image
	if img != "" {
		img = abs(c.SiteURL, img)
	}

	if m.Title != "" {
		w("<title>" + html.EscapeString(m.Title) + "</title>")
		w(meta("property", "og:title", m.Title))
		w(meta("name", "twitter:title", m.Title))
	}
	if m.Description != "" {
		w(meta("name", "description", m.Description))
		w(meta("property", "og:description", m.Description))
		w(meta("name", "twitter:description", m.Description))
	}
	if canon != "" {
		w(`<link rel="canonical" href="` + html.EscapeString(canon) + `">`)
		w(meta("property", "og:url", canon))
	}
	w(meta("property", "og:type", ogType))
	if m.SiteName != "" {
		w(meta("property", "og:site_name", m.SiteName))
	}
	if img != "" {
		w(meta("property", "og:image", img))
		w(meta("name", "twitter:image", img))
		w(meta("name", "twitter:card", "summary_large_image"))
	} else {
		w(meta("name", "twitter:card", "summary"))
	}

	// Search-engine ownership verification.
	if c.GSCVerification != "" {
		w(meta("name", "google-site-verification", c.GSCVerification))
	}
	if c.BingVerification != "" {
		w(meta("name", "msvalidate.01", c.BingVerification))
	}

	if m.JSONLD != nil {
		if j, err := json.Marshal(m.JSONLD); err == nil {
			w(`<script type="application/ld+json">` + string(j) + `</script>`)
		}
	}

	// Google Analytics (gtag.js).
	if c.GAID != "" {
		b.WriteString(GASnippet(c.GAID))
		b.WriteByte('\n')
	}
	return b.String()
}

func meta(attr, key, content string) string {
	return `<meta ` + attr + `="` + key + `" content="` + html.EscapeString(content) + `">`
}

// GASnippet returns the Google Analytics gtag.js <head> snippet for the given id.
func GASnippet(id string) string {
	id = html.EscapeString(id)
	return `<!-- Google tag (gtag.js) -->
<script async src="https://www.googletagmanager.com/gtag/js?id=` + id + `"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());
  gtag('config', '` + id + `');
</script>`
}

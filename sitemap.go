package seo

import (
	"encoding/xml"
	"net/http"
	"strings"
)

type sitemapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

type urlSet struct {
	XMLName xml.Name     `xml:"urlset"`
	NS      string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

// abs turns a root-relative loc into an absolute URL using SiteURL.
func abs(base, loc string) string {
	if loc == "" {
		return base
	}
	if strings.HasPrefix(loc, "http://") || strings.HasPrefix(loc, "https://") {
		return loc
	}
	if !strings.HasPrefix(loc, "/") {
		loc = "/" + loc
	}
	return base + loc
}

func handleSitemap(w http.ResponseWriter, _ *http.Request) {
	mu.RLock()
	src := urlSource
	base := cfg.SiteURL
	mu.RUnlock()

	set := urlSet{NS: "http://www.sitemaps.org/schemas/sitemap/0.9"}
	if src != nil {
		for _, u := range src() {
			set.URLs = append(set.URLs, sitemapURL{
				Loc:        abs(base, u.Loc),
				LastMod:    u.LastMod,
				ChangeFreq: u.ChangeFreq,
				Priority:   u.Priority,
			})
		}
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	_, _ = w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	_ = enc.Encode(set)
}

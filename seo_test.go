package seo

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMetaHTML(t *testing.T) {
	os.Setenv("SEO_SITE_URL", "https://to-go.dev")
	os.Setenv("SEO_GA_ID", "G-TEST123")
	os.Setenv("SEO_GSC_VERIFICATION", "gsc-token")
	mu.Lock()
	cfg = loadConfig()
	mu.Unlock()

	out := Meta{
		Title:       "ToGO",
		Description: "Full-stack Go + React",
		Canonical:   "/plugins",
		Image:       "/og.png",
		JSONLD:      map[string]any{"@type": "SoftwareApplication", "name": "ToGO"},
	}.HTML()

	for _, want := range []string{
		"<title>ToGO</title>",
		`property="og:title"`,
		`name="description"`,
		`rel="canonical" href="https://to-go.dev/plugins"`,
		`property="og:image" content="https://to-go.dev/og.png"`,
		`name="google-site-verification" content="gsc-token"`,
		`application/ld+json`,
		"G-TEST123",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Meta.HTML() missing %q\n--- got ---\n%s", want, out)
		}
	}
}

func TestSitemap(t *testing.T) {
	os.Setenv("SEO_SITE_URL", "https://to-go.dev")
	mu.Lock()
	cfg = loadConfig()
	urlSource = func() []URL { return []URL{{Loc: "/"}, {Loc: "/plugins", ChangeFreq: "weekly"}} }
	mu.Unlock()

	rec := httptest.NewRecorder()
	handleSitemap(rec, httptest.NewRequest("GET", "/sitemap.xml", nil))
	body := rec.Body.String()
	for _, want := range []string{
		"<loc>https://to-go.dev/</loc>",
		"<loc>https://to-go.dev/plugins</loc>",
		"<changefreq>weekly</changefreq>",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("sitemap missing %q\n%s", want, body)
		}
	}
}

func TestRobots(t *testing.T) {
	os.Setenv("SEO_SITE_URL", "https://to-go.dev")
	mu.Lock()
	cfg = loadConfig()
	mu.Unlock()
	rec := httptest.NewRecorder()
	handleRobots(rec, httptest.NewRequest("GET", "/robots.txt", nil))
	if !strings.Contains(rec.Body.String(), "Sitemap: https://to-go.dev/sitemap.xml") {
		t.Errorf("robots missing sitemap line: %s", rec.Body.String())
	}
}

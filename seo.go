// Package seo is togo's SEO/AEO provider. Blank-import it (or `togo install
// togo-framework/seo`) and it mounts /sitemap.xml, /robots.txt and /llms.txt on
// the kernel router, serves search-engine verification + IndexNow key files, and
// exposes <head> helpers (meta/OG/Twitter/JSON-LD + Google Analytics) for SSR.
//
// AEO (Answer-Engine Optimization): /llms.txt advertises the site to LLM agents,
// and the app can serve per-page raw Markdown so agents read content directly.
//
// Everything is opt-in via the environment (togo convention — .env + hooks):
//
//	SEO_SITE_URL          canonical base, e.g. https://to-go.dev   (required for absolute URLs)
//	SEO_INDEXNOW_KEY      enables IndexNow + serves /<key>.txt
//	SEO_GA_ID             Google Analytics gtag id (G-XXXX) → GA <head> snippet
//	SEO_GSC_VERIFICATION  google-site-verification meta content
//	SEO_BING_VERIFICATION msvalidate.01 meta content
package seo

import (
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/togo-framework/togo"
)

// Config is resolved from the environment.
type Config struct {
	SiteURL          string
	IndexNowKey      string
	GAID             string
	GSCVerification  string
	BingVerification string
}

func loadConfig() Config {
	return Config{
		SiteURL:          strings.TrimRight(os.Getenv("SEO_SITE_URL"), "/"),
		IndexNowKey:      os.Getenv("SEO_INDEXNOW_KEY"),
		GAID:             os.Getenv("SEO_GA_ID"),
		GSCVerification:  os.Getenv("SEO_GSC_VERIFICATION"),
		BingVerification: os.Getenv("SEO_BING_VERIFICATION"),
	}
}

// URL is one sitemap entry. Loc may be absolute or root-relative (SiteURL is prepended).
type URL struct {
	Loc        string
	LastMod    string // ISO-8601, optional
	ChangeFreq string // e.g. "weekly", optional
	Priority   string // e.g. "0.8", optional
}

var (
	mu        sync.RWMutex
	cfg       Config
	urlSource func() []URL
	llmsFn    func() string
)

// RegisterURLs sets the source of sitemap URLs. Called by the app at boot; the
// func is evaluated on each /sitemap.xml request so the map can stay current.
func RegisterURLs(f func() []URL) { mu.Lock(); urlSource = f; mu.Unlock() }

// RegisterLLMS sets the builder for the /llms.txt body (the AEO index).
func RegisterLLMS(f func() string) { mu.Lock(); llmsFn = f; mu.Unlock() }

// Settings returns the resolved config (after the provider has booted).
func Settings() Config { mu.RLock(); defer mu.RUnlock(); return cfg }

func init() {
	togo.RegisterProviderFunc("seo", togo.PriorityService, func(k *togo.Kernel) error {
		mu.Lock()
		cfg = loadConfig()
		if cfg.SiteURL == "" {
			cfg.SiteURL = strings.TrimRight(os.Getenv("ADDR"), "/") // best-effort fallback
		}
		mu.Unlock()
		mount(k.Router)
		return nil
	})
}

func mount(r chi.Router) {
	r.Get("/sitemap.xml", handleSitemap)
	r.Get("/robots.txt", handleRobots)
	r.Get("/llms.txt", handleLLMS)

	c := Settings()
	if c.IndexNowKey != "" {
		// IndexNow ownership proof: the key, served as plain text at /<key>.txt.
		r.Get("/"+c.IndexNowKey+".txt", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = io.WriteString(w, c.IndexNowKey)
		})
	}
}

func handleRobots(w http.ResponseWriter, _ *http.Request) {
	c := Settings()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	b := &strings.Builder{}
	b.WriteString("User-agent: *\nAllow: /\n")
	if c.SiteURL != "" {
		b.WriteString("\nSitemap: " + c.SiteURL + "/sitemap.xml\n")
	}
	_, _ = io.WriteString(w, b.String())
}

func handleLLMS(w http.ResponseWriter, _ *http.Request) {
	mu.RLock()
	fn := llmsFn
	mu.RUnlock()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if fn != nil {
		_, _ = io.WriteString(w, fn())
		return
	}
	c := Settings()
	_, _ = io.WriteString(w, "# "+c.SiteURL+"\n\n> Powered by togo (https://to-go.dev).\n")
}

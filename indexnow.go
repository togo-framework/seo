package seo

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

// IndexNow is the multi-engine instant-indexing protocol (Bing, Yandex, Seznam,
// Naver — and Google honors the same submissions via its own pipeline). Submit
// changed URLs after a deploy to get them crawled fast.
//
// Requires SEO_INDEXNOW_KEY (the plugin serves the proof file at /<key>.txt).

type indexNowPayload struct {
	Host        string   `json:"host"`
	Key         string   `json:"key"`
	KeyLocation string   `json:"keyLocation"`
	URLList     []string `json:"urlList"`
}

// SubmitIndexNow notifies IndexNow of (created/updated/deleted) URLs. No-op error
// if SEO_INDEXNOW_KEY or SEO_SITE_URL is unset. Safe to call from a deploy hook.
func SubmitIndexNow(urls []string) error {
	c := Settings()
	if c.IndexNowKey == "" {
		return errors.New("seo: SEO_INDEXNOW_KEY not set")
	}
	if c.SiteURL == "" {
		return errors.New("seo: SEO_SITE_URL not set")
	}
	if len(urls) == 0 {
		return nil
	}
	host := ""
	if u, err := url.Parse(c.SiteURL); err == nil {
		host = u.Host
	}
	abs := make([]string, 0, len(urls))
	for _, u := range urls {
		abs = append(abs, absURL(c.SiteURL, u))
	}
	body, _ := json.Marshal(indexNowPayload{
		Host:        host,
		Key:         c.IndexNowKey,
		KeyLocation: c.SiteURL + "/" + c.IndexNowKey + ".txt",
		URLList:     abs,
	})
	req, err := http.NewRequest(http.MethodPost, "https://api.indexnow.org/indexnow", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return errors.New("seo: indexnow returned " + resp.Status)
	}
	return nil
}

// absURL is the package-level helper mirroring abs() for exported callers.
func absURL(base, loc string) string { return abs(base, loc) }

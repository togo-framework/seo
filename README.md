<p align="center"><img src="https://raw.githubusercontent.com/togo-framework/.github/main/profile/assets/togo-mark.svg" width="96" alt="togo"></p>
<h1 align="center">@togo-framework/seo</h1>
<p align="center"><sub>part of the <a href="https://github.com/togo-framework">togo-framework</a> — the full-stack Go + React framework</sub></p>

The togo **SEO / AEO** provider. Blank-import it (or `togo install togo-framework/seo`)
and it mounts `/sitemap.xml`, `/robots.txt` and `/llms.txt` on the kernel router, serves
search-engine verification + IndexNow key files, and gives you `<head>` helpers
(meta / OpenGraph / Twitter / JSON-LD + Google Analytics). Every provider is opt-in via
the environment — nothing turns on until you configure it.

```bash
togo install togo-framework/seo
```

## Configure (`.env`)

| Env | Effect |
|---|---|
| `SEO_SITE_URL` | Canonical base (e.g. `https://to-go.dev`). Required for absolute URLs. |
| `SEO_INDEXNOW_KEY` | Enables **IndexNow**; serves the proof file at `/<key>.txt`. |
| `SEO_GA_ID` | **Google Analytics** gtag id (`G-XXXX`) → injected GA `<head>` snippet. |
| `SEO_GSC_VERIFICATION` | **Google Search Console** `google-site-verification` meta. |
| `SEO_BING_VERIFICATION` | **Bing** `msvalidate.01` meta. |

## Use (Go / backend)

```go
import "github.com/togo-framework/seo"

// Feed the sitemap + AEO index from your routes/resources (evaluated per request).
seo.RegisterURLs(func() []seo.URL {
    return []seo.URL{{Loc: "/"}, {Loc: "/plugins", ChangeFreq: "weekly", Priority: "0.8"}}
})
seo.RegisterLLMS(func() string { return "# My app\n\n> what it does…\n" })

// Per-page <head> (SSR): meta + OG + Twitter + JSON-LD + verification + GA.
head := seo.Meta{
    Title: "ToGO", Description: "Full-stack Go + React", Canonical: "/plugins",
    Image: "/og.png", JSONLD: map[string]any{"@type": "SoftwareApplication", "name": "ToGO"},
}.HTML()

// After a deploy, ping IndexNow with the URLs that changed:
_ = seo.SubmitIndexNow([]string{"/", "/plugins"})
```

The provider auto-serves:
- `GET /sitemap.xml` — from `RegisterURLs`
- `GET /robots.txt` — allow-all + `Sitemap:` line
- `GET /llms.txt` — from `RegisterLLMS` (AEO: advertise your site to LLM agents)
- `GET /<SEO_INDEXNOW_KEY>.txt` — IndexNow ownership proof

## Use (React / frontend)

```tsx
import { Helmet } from "react-helmet-async";
import { seoTags, gaSnippet } from "@togo-framework/seo/frontend/seo";

<Helmet>{seoTags({ title, description, canonical: path, image, siteUrl: "https://to-go.dev", jsonLd })}</Helmet>
// inject GA once:
<script dangerouslySetInnerHTML={{ __html: gaSnippet("G-XXXX") }} />
```

## AEO (Answer-Engine Optimization)

Beyond `/llms.txt`, serve each page's content as **raw Markdown** (e.g. `/docs/x.md` as
`text/plain`) so LLM agents read it directly — pair this plugin with that convention and
agents get clean, structured answers about your app.

---

## 💎 Premium sponsors

togo is proudly sponsored by **ID8 Media** and **One Studio**.

<p align="center">
  <a href="https://id8media.com"><img src="https://raw.githubusercontent.com/togo-framework/.github/main/profile/assets/id8media.svg" height="44" alt="ID8 Media" /></a>
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
  <a href="https://one-studio.co"><img src="https://raw.githubusercontent.com/togo-framework/.github/main/profile/assets/one-studio.jpeg" height="44" alt="One Studio" /></a>
</p>

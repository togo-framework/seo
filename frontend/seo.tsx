/**
 * togo SEO/AEO — frontend helpers (framework-agnostic, React).
 *
 * Pairs with the Go `seo` plugin: the server owns /sitemap.xml, /robots.txt,
 * /llms.txt and verification files; the client owns per-page <head> tags.
 *
 * Works with any head manager. With react-helmet-async:
 *
 *   import { Helmet } from "react-helmet-async";
 *   import { seoTags, gaSnippet } from "@togo-framework/seo/frontend/seo";
 *   <Helmet>{seoTags({ title, description, canonical: path, image, jsonLd })}</Helmet>
 *
 * Inject the GA snippet once (e.g. in index.html or a root Helmet):
 *   <script dangerouslySetInnerHTML={{ __html: gaSnippet("G-XXXX") }} />
 */
import * as React from "react";

export interface SeoInput {
  title?: string;
  description?: string;
  /** absolute or root-relative; combine with `siteUrl` for canonical/OG */
  canonical?: string;
  image?: string;
  type?: string; // og:type, default "website"
  siteName?: string;
  siteUrl?: string; // e.g. https://to-go.dev
  jsonLd?: unknown; // schema.org object → ld+json
}

function abs(base: string | undefined, loc?: string): string | undefined {
  if (!loc) return undefined;
  if (/^https?:\/\//.test(loc)) return loc;
  if (!base) return loc;
  return base.replace(/\/$/, "") + (loc.startsWith("/") ? loc : "/" + loc);
}

/** Returns an array of head elements — drop into <Helmet> or a custom head. */
export function seoTags(input: SeoInput): React.ReactNode[] {
  const { title, description, canonical, image, type = "website", siteName, siteUrl, jsonLd } = input;
  const url = abs(siteUrl, canonical);
  const img = abs(siteUrl, image);
  const tags: React.ReactNode[] = [];
  if (title) {
    tags.push(<title key="t">{title}</title>);
    tags.push(<meta key="ogt" property="og:title" content={title} />);
    tags.push(<meta key="twt" name="twitter:title" content={title} />);
  }
  if (description) {
    tags.push(<meta key="d" name="description" content={description} />);
    tags.push(<meta key="ogd" property="og:description" content={description} />);
    tags.push(<meta key="twd" name="twitter:description" content={description} />);
  }
  if (url) {
    tags.push(<link key="c" rel="canonical" href={url} />);
    tags.push(<meta key="ogu" property="og:url" content={url} />);
  }
  tags.push(<meta key="ogtype" property="og:type" content={type} />);
  if (siteName) tags.push(<meta key="ogs" property="og:site_name" content={siteName} />);
  if (img) {
    tags.push(<meta key="ogi" property="og:image" content={img} />);
    tags.push(<meta key="twi" name="twitter:image" content={img} />);
    tags.push(<meta key="twc" name="twitter:card" content="summary_large_image" />);
  } else {
    tags.push(<meta key="twc" name="twitter:card" content="summary" />);
  }
  if (jsonLd) {
    tags.push(
      <script key="ld" type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }} />,
    );
  }
  return tags;
}

/** The Google Analytics gtag.js snippet (inject once). */
export function gaSnippet(id: string): string {
  return `<!-- Google tag (gtag.js) -->
<script async src="https://www.googletagmanager.com/gtag/js?id=${id}"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());
  gtag('config', '${id}');
</script>`;
}

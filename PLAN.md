# Coffeeroasters — Implementation Plan

A meta aggregator for specialty coffee roasters: discover them, scrape their offerings, normalise the data, and present it in a clean UI.

---

## Phase 1: Discovery & List Building

**Goal:** Compile a seed list of roasters to track.

### Sources
- Curated directories: [Roast Magazine](https://roastmagazine.com), [Coffee Review](https://coffeereview.com), [Perfect Daily Grind](https://perfectdailygrind.com)
- Community lists: Reddit r/coffee wikis, specialty coffee forums, Instagram hashtags
- Barista competition sponsors and entrants
- Existing open datasets (e.g. OpenStreetMap `craft=coffee_roasters`)

### Approach
- Manual seed list in a YAML/CSV file to bootstrap (`roasters.yaml`)
- Automated discovery scraper to follow links from directory pages
- De-duplication by domain name + fuzzy name matching
- Human review step before adding new roasters to DB

### Deliverables
- `roasters.yaml` — initial seed list (~100 roasters)
- `cmd/discover/` — Go CLI to crawl discovery sources
- Roaster table seeded in Postgres

---

## Phase 2: Web Scraping

**Goal:** Extract coffee product listings from each roaster's website.

### Strategy
- Per-roaster scraper configs (CSS selectors / JSON-LD / structured data)
- Prefer structured data (Schema.org `Product`, JSON-LD) where available
- Fall back to CSS selector scraping for bespoke sites
- Headless browser (via `chromedp`) only as last resort (JS-heavy Shopify stores)

### Tech
- `colly` or `net/http` + `goquery` for static HTML
- `chromedp` for JS-rendered pages
- Scraper config per roaster stored in DB or YAML

### Fields to Extract
- Coffee name, origin(s), process, variety/cultivar
- Roast level, tasting notes
- Price, weight, currency
- In-stock status, product URL, image URL

### Deliverables
- `internal/scraper/` — scraper framework + per-roaster implementations
- `cmd/scrape/` — CLI to run scrapers (single roaster or all)
- Scrape run logging in `scrape_runs` table

---

## Phase 3: Data Normalisation

**Goal:** Standardise free-text fields into queryable, comparable data.

### Challenges & Solutions
- **Origins**: map free text → country + region using a reference list
- **Process**: normalise variants ("Natural", "Dry Process", "Naturale" → `natural`)
- **Roast level**: map to enum (`light`, `medium-light`, `medium`, `medium-dark`, `dark`)
- **Tasting notes**: extract from comma-separated strings, map to flavour wheel categories
- **Currency/weight**: normalise to cents + grams

### Approach
- Normalisation functions in `internal/normalise/`
- Reference data files (`countries.yaml`, `processes.yaml`, `flavour-wheel.yaml`)
- LLM-assisted extraction for freeform descriptions (optional, later phase)

### Deliverables
- `internal/normalise/` package
- Normalised fields stored alongside raw scraped data
- Admin UI for reviewing/correcting normalisation

---

## Phase 4: Backend API (Go + PostgreSQL)

**Goal:** Expose the data via a clean JSON API.

### Endpoints
```
GET  /api/roasters              # list roasters (filter, sort, paginate)
GET  /api/roasters/:slug        # single roaster + their coffees
GET  /api/coffees               # search/filter coffees across all roasters
GET  /api/coffees/:id           # single coffee detail
GET  /api/origins               # list of origins for filter UI
GET  /api/search?q=             # full-text search
```

### Tech
- Standard library `net/http` + `chi` router
- `pgx` for Postgres queries (or `sqlc` for type-safe generated queries)
- Full-text search via Postgres `tsvector`
- JSON responses, no GraphQL for now

### Features
- Filtering: by origin, process, roast level, price range, country
- Sorting: by name, price, recently updated
- Pagination: cursor-based
- Response caching with `Cache-Control` headers

### Deliverables
- `internal/api/` — route handlers
- `internal/db/` — query layer
- `backend/Dockerfile`
- API integration tests

---

## Phase 5: SvelteKit Frontend

**Goal:** A fast, clean UI for browsing and discovering coffees.

### Pages
- `/` — landing / featured roasters + coffees
- `/roasters` — browse all roasters (filter by country/region)
- `/roasters/[slug]` — roaster profile + their current offerings
- `/coffees` — search & filter all coffees
- `/origins/[country]` — coffees by origin country
- `/about` — project info, data sources, ethical notes

### Features
- Filter sidebar: origin, process, roast level, price, in-stock only
- Full-text search with instant results
- Coffee cards with tasting notes tags, roast level badge
- Roaster cards with country flag, region
- Responsive, mobile-first design
- Light/dark mode

### Tech
- SvelteKit with server-side rendering
- Tailwind CSS for styling
- `@tanstack/svelte-query` for data fetching / caching
- Deployed to Fly.io or Vercel

### Deliverables
- Full SvelteKit app in `frontend/`
- Storybook component library (optional)
- E2E tests with Playwright

---

## Phase 6: Regular Updates & Scheduling

**Goal:** Keep data fresh automatically.

### Strategy
- Scheduled scrape jobs (cron or pg-backed job queue)
- Scrape frequency per roaster (daily for active, weekly for slow-moving)
- Change detection: only update DB rows when data actually changes
- Alert on scrape failures (email / Slack webhook)

### Tech Options
- Simple: cron job on the server calling `cmd/scrape`
- Better: `pgq` or `riverqueue` (Go-native Postgres job queue)

### Deliverables
- `cmd/scheduler/` — job scheduler
- Scrape health dashboard (admin page)
- Alerting on repeated failures

---

## Phase 7: Ethical Scraping & Compliance

**Goal:** Be a good citizen of the web.

### Practices
- **robots.txt**: always check and respect
- **Rate limiting**: minimum 2–5 seconds between requests per domain
- **User-Agent**: identify as `CoffeeroastersBot/1.0 (+https://coffeeroasters.app)`
- **Caching**: never re-fetch unchanged pages (use ETags / Last-Modified)
- **No login walls**: only scrape publicly accessible pages
- **Terms of Service**: review each roaster's ToS; skip if scraping is prohibited
- **Contact page**: provide an opt-out mechanism for roasters who don't want to be listed

### Data & Privacy
- Only store publicly listed product data — no personal data
- Link back to the original roaster page for each coffee
- Attribution on every listing

### Deliverables
- `internal/scraper/robots.go` — robots.txt checker
- `internal/scraper/ratelimiter.go` — per-domain rate limiter
- Opt-out list in DB (`roasters.opted_out BOOLEAN`)
- Public "About our data" page

---

## Milestones

| # | Milestone | Target |
|---|-----------|--------|
| 1 | Repo setup, schema, seed list | Week 1 |
| 2 | 10 working scrapers, data in DB | Week 3 |
| 3 | Basic API + SvelteKit UI live | Week 5 |
| 4 | 50+ roasters, normalisation | Week 7 |
| 5 | Scheduled updates, alerts | Week 9 |
| 6 | Public launch | Week 12 |

---

## Open Questions

- Monetisation? (affiliate links, sponsorship, none)
- User accounts / wishlist / price alerts?
- Submission form for roasters to add themselves?
- Mobile app later?

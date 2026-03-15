# Coffeeroasters

A meta aggregator for discovering and exploring **Australian indie coffee roasters**. Automatically discovers roasters from directories and community sources, scrapes their current offerings, and presents a unified, normalised view — so you can search every bean from every roaster in one place.

## What it does

- Crawls Australian specialty coffee directories and community sources to build a comprehensive roaster index
- Scrapes each roaster's website for their current coffee offerings (name, origin, process, roast level, tasting notes, price)
- Normalises free-text fields into queryable, comparable data
- Exposes a JSON API for filtering and full-text search
- Presents a clean, fast frontend for browsing and discovering coffees

## Stack

- **Backend**: Go + PostgreSQL
- **Frontend**: Next.js (App Router, TypeScript)
- **Scraping**: Go-based scrapers with ethical rate limiting

## Getting Started

### Prerequisites

- Go 1.22+
- Node.js 20+
- Docker + Docker Compose (for Postgres)

### Development

```bash
# Start Postgres
docker compose up -d db

# Run database migrations
just migrate

# Run backend (with live reload)
just dev-backend

# Run frontend
just dev-frontend
```

Or run everything at once:

```bash
just dev
```

### Environment

Copy `.env.example` to `.env` and fill in:

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/coffeeroasters
PORT=8080
```

## Project Structure

```
coffeeroasters/
├── cmd/
│   ├── server/         # HTTP server entrypoint
│   ├── scrape/         # CLI: run scrapers for one or all roasters
│   └── discover/       # CLI: crawl discovery sources to find new roasters
├── internal/
│   ├── api/            # HTTP handlers
│   ├── scraper/        # Scraper framework (runner, rate limiter, robots.txt)
│   ├── normalise/      # Data normalisation (origin, process, roast level)
│   ├── db/             # Postgres query layer (sqlc-generated)
│   ├── domain/         # Domain types
│   └── config/         # App configuration
├── roasters/
│   └── configs/        # Per-roaster scraper config (YAML)
├── data/               # Reference data for normalisation
├── migrations/         # Goose SQL migrations
├── queries/            # sqlc SQL query files
└── src/                # Next.js frontend (App Router)
```

## Adding a New Roaster

1. Copy `roasters/configs/_template.yaml` to `roasters/configs/<roaster-slug>.yaml`
2. Fill in the selectors for their product listing page
3. Run `go run ./cmd/scrape -roaster <slug> -dry-run` to verify output
4. Open a PR — all scraper configs are reviewed before going live

See [roasters/configs/_template.yaml](roasters/configs/_template.yaml) for full documentation of every field.

## Ethical Scraping

- Respects `robots.txt` on every domain
- Enforces a minimum 3–5 second delay between requests per domain
- Identifies as `CoffeeroastersBot/1.0 (+https://coffeeroasters.app)`
- Only scrapes publicly accessible, non-login-walled pages
- Roasters can opt out by emailing hello@coffeeroasters.app

## License

MIT

# Coffeeroasters

A meta aggregator for discovering and exploring specialty coffee roasters. Automatically discovers roasters from directories and community sources, scrapes their offerings, and presents a unified, normalised view.

## Stack

- **Backend**: Go + PostgreSQL
- **Frontend**: SvelteKit
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

# Run backend
cd backend && go run .

# Run frontend
cd frontend && npm install && npm run dev
```

## Project Structure

```
coffeeroasters/
├── backend/          # Go API server + scrapers
│   ├── main.go
│   ├── go.mod
│   ├── internal/
│   │   ├── api/      # HTTP handlers
│   │   ├── db/       # Postgres queries (sqlc)
│   │   ├── scraper/  # Roaster scrapers
│   │   └── models/   # Domain models
│   └── schema.sql    # DB schema
├── frontend/         # SvelteKit UI
└── docker-compose.yml
```

## License

MIT

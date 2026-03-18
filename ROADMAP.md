# Coffeeroasters Roadmap

## Quick Wins

- [x] **Normalise tags to lower case** — tasting notes lowercased during normalisation (e36849f)
- [x] **Similar coffees less prominent** — reduced to 3, styled as "You might also like" (e36849f)
- [x] **Similarity reasons** — human-readable reasons replace percentage match (e36849f)

## Features

- [x] **Coffee tracking (likes/have)** — localStorage tracking with like/tried buttons, /my-coffees page (4301fb8)
- [ ] **Historical availability tracking** — track when coffees come in/out of stock so users can see ones they historically liked
- [ ] **Likes-based recommendations** — use liked coffees to inform similarity recommendations
- [x] **Scheduled scraping** — daily cron via Alpine crond, configurable via SCRAPE_CRON env var (a65c782)
- [x] **Small Batch image extraction** — general HTML image extraction via og:image + fallback selectors (a65c782)
- [ ] **Expanded extractors** — extract descriptive text (flavour writeups, producer info) from seller sites; expose in UI and use for similarity matching (embeddings)
- [ ] **Placeholder images** — generate placeholder coffee bag/box images with company logos in studio ghibli/watercolour style (GPT image generation, similar to recipe app)

## Research

- [ ] **Expand roaster list** — research indie roasters across all Australian major cities and state capitals

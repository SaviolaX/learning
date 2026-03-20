# URL Shortener (v0 / v1 / v2)

This folder contains three incremental versions of the same URL shortener idea.

## Key differences
- API path:
  - v0 uses `POST /short-url`
  - v1 and v2 use `POST /shorten`
- Storage backend:
  - v0/v1 store everything in JSON files (`urlShortenerDB.json` / `UrlShortenerDB.json`)
  - v2 uses SQLite (`urls.db`) with a `urls` table
- Architecture:
  - v0: single-file HTTP server (plus `pkg/hasher` and `pkg/storage`)
  - v1: internal packages (server, URL validation, JSON repository)
  - v2: separated `domain`/`application`/`infrastructure`/`transport` with a repository interface
- Collision handling:
  - v0/v1: dedup primarily by hash code (v1 doesn’t explicitly check for “same code, different URL”)
  - v2: checks collision by verifying the original URL for an existing code
- UI behavior:
  - v2 shows error messages from failed shorten requests

## Versions
- [`url_shortener_v0/README.md`](./url_shortener_v0/README.md)
- [`url_shortener_v1/README.md`](./url_shortener_v1/README.md)
- [`url_shortener_v2/README.md`](./url_shortener_v2/README.md)


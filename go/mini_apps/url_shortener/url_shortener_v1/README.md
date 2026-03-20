# URL Shortener (v1)

Refactored Go URL shortener with internal modules (server, URL validation, JSON repository) and a slightly cleaner HTTP flow.

## Run
`cd url_shortener_v1 && go run ./cmd`
Open `http://localhost:8080/`.

## API
- `GET /` serves `templates/index.html`
- `POST /shorten` (form field: `url`) returns the shortened URL text
- `GET /r/<code>` redirects to the original URL

## Notes
- Short code is `SHA256(normalized(url))` truncated to 10 characters.
- Stores mappings in `UrlShortenerDB.json` (in this folder).
- Deduping is done by hash code: if the code exists, it reuses it (no explicit collision handling).

## Tests
`make test`


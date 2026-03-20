# URL Shortener (v0)

Minimal Go `net/http` URL shortener with JSON file storage.

## Run
`cd url_shortener_v0 && go run ./cmd`
Open `http://localhost:8080/`.

## API
- `GET /` serves `templates/index.html`
- `POST /short-url` (form field: `url`) returns the shortened URL text
- `GET /r/<hash>` redirects to the original URL

## Notes
- Short code is `SHA256(url)` truncated to 10 characters.
- Stores mappings in `urlShortenerDB.json` (in this folder).


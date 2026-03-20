# URL Shortener (v2)

Clean-ish layered refactor (domain/application/infrastructure/transport) and an SQLite backend.

## Run
`cd url_shortener_v2 && go run ./cmd`
Open `http://localhost:8080/`.

## API
- `GET /` serves `templates/index.html`
- `POST /shorten` (form field: `url`) returns the shortened URL text
- `GET /r/<code>` redirects to the original URL

## Storage
- Uses SQLite database file: `urls.db`
- Creates table `urls (code TEXT PRIMARY KEY, original TEXT NOT NULL)`

## Notes
- URL normalization/validation lives in `internal/domain/url`.
- Service layer checks for hash-code collisions:
  - if the code exists with the same original URL => reuse the code
  - if the code exists with a different original URL => returns an error
- UI template shows errors from non-`2xx` responses.

## Tests
`make test`


# URL Shortener Improvement Plan

This document outlines a focused plan to improve the existing simple URL shortener without adding new feature categories such as databases, JWT, or external services.

## Goals
- Keep the app minimal and focused on existing endpoints
- Improve reliability and correctness
- Fix error handling flows
- Add tests and minimal refactor
- Keep no external DB, no JWT, no auth, no extra endpoints

## 1) Robust Error Handling
- `shortenUrl`: handle `ParseForm` failures with `http.StatusBadRequest`.
- `shortenUrl`: validate `url`; return `http.StatusBadRequest` for empty input.
- Check and handle `s.repo.Store()` errors and return `http.StatusInternalServerError`.
- `redirectUrl`: return `http.StatusInternalServerError` if `repo.Load()` fails.

## 2) Input validation
- `strings.TrimSpace` on incoming URL.
- Reject empty URLs after trim.
- Optional simple scheme check for `http://` or `https://`.

## 3) Short URL mapping clarity
- Keep existing URL returned (e.g. `http://localhost:8080/r/abc123`) while storing a token mapping.
- In `redirectUrl`, extract token from path (e.g. `/r/{token}`) and lookup by token.
- Avoid direct host normalization in storage key.

## 4) Concurrency safety for file storage
- Add `sync.Mutex` to `Repository`.
- Lock for `Load` and `Store` to avoid race conditions.

## 5) Test coverage
- `pkg/hasher/hasher_test.go`: test empty URL, invalid hash lengths, valid results.
- `pkg/storage/storage_test.go`: roundtrip load/store with temp file.
- `main_test.go`: tests for `buildMap`, `shortenUrl`, `redirectUrl`, HTTP status behavior.

## 6) Simple config (optional)
- Keep defaults but add `flag` options: `port`, `db path`, `hash length`, `base URL`.
- Do not require external config.

## 7) Code cleanup
- Remove alias imports in `main.go` (use plain packages).
- Rename handlers to Go style (`shortenURL`, `redirectURL`).

## 8) Maintain existing endpoints
- `GET /` serves index.
- `POST /short-url` shortens URL.
- `GET /r/{token}` redirects.

## Next steps
1. Implement and run `go test ./...`.
2. Add and run tests in fast iterations.
3. Keep small commits and keep plan as a roadmap.

# Backend Tests Guide

Scope: `backend/tests/`.

Use this package for integration-style backend tests around Gin handlers, middleware, and database-backed behavior.

- Prefer Go's standard `testing` package.
- Use `httptest` for route tests.
- Use the in-memory SQLite helper for DB-backed controller tests when practical.
- Restore global state such as `models.DB` with `t.Cleanup`.
- Keep tests focused on observable behavior: status codes, headers, response JSON, auth enforcement, and persistence effects.

Run from `backend/`:

```sh
go test ./...
```

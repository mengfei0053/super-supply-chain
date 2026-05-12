# Backend Guide

Scope: everything under `backend/`.

This is the Go API server. Keep changes aligned with the existing Gin, GORM, and helper-package layout.

## Entry Points

- `main.go` wires config, DB initialization, middleware, static file serving, and routes.
- `configs/config.go` loads environment variables and development `.env` data.
- `models/init-db.go` initializes the global `models.DB`.

## Commands

Run from `backend/`:

```sh
go test ./...
go build ./...
```

For local hot reload, run `./dev.sh` from the repository root. It changes into `backend/` and runs `air`.

## Conventions

- Use standard Go formatting (`gofmt`) for touched Go files.
- Prefer the standard `testing` package for tests. Add focused helpers in `backend/tests` when API setup is shared.
- Keep controller request/response behavior compatible with React Admin where existing endpoints already do so.
- Do not introduce new package-level globals unless they match the current DB/config pattern and are necessary.
- Treat `backend/tmp`, `backend/logs`, and generated binaries as runtime artifacts. Do not edit them for feature work.
- Treat files in `backend/uploads` as business templates or fixtures; confirm before replacing or deleting them.

## Environment

Backend code reads:

- `ENVIRONMENT`
- `PORT`
- `MYSQL_USER`
- `MYSQL_PASSWORD`
- `MYSQL_SERVER`
- `UPLOAD_USER`
- `UPLOAD_PASSWORD`
- `UPLOAD_SERVER`

Avoid printing new secrets or expanding existing credential logging.

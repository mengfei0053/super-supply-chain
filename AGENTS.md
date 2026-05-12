# AGENTS.md

## Project Overview

This repository contains Super Supply Chain, a supply-chain admin system with:

- `backend/`: Go 1.23.6 API server using Gin, GORM, MySQL, JWT auth, Excel parsing/export helpers, and optional WebDAV upload/download.
- `frontend/`: React 19 + Vite + React Admin UI, with MUI components and custom pages for settlement forms, dynamic Excel tables, dict management, read rules, and Yifan cost calculation.
- `ssc-sqls/`: SQL initialization or shared database scripts.
- `Dockerfile`: multi-stage build that compiles the frontend, copies it into the Go image, builds `backend`, and serves the production app.
- `ssc-deployment.yaml`: Kubernetes deployment/service manifest for namespace `ssc`.
- `dev.sh` and `run.sh`: local backend start scripts.

`node-backend/` was removed in commit `98dd5c9` and should not be reintroduced unless explicitly requested.

## Repository Notes

`batch-print/` and `py-backend/` were removed in commit `ebb0b43`.

Detailed local guidance lives in nested `AGENTS.md` files under `backend/` and `frontend/`. Read the closest one before editing files in those directories.

## Backend

Entry point: `backend/main.go`

Main behavior:

- Loads config with `configs.LoadConfigFile()`.
- Sets Gin release mode.
- Connects to MySQL via `models.InitDB()`.
- Installs zap-backed Gin logger and recovery middleware.
- Serves frontend static files at `/super-supply-chain` only in production.
- Redirects `/` to `/super-supply-chain`.
- Public auth endpoints:
  - `POST /api/register`
  - `POST /api/login`
- Protected API group under `/api/admin`, guarded by `middleware.AuthMiddleware()`.

Important backend directories:

- `backend/controllers/`: HTTP handlers and route behavior.
- `backend/models/`: GORM models and DB initialization.
- `backend/middleware/`: auth and logging middleware.
- `backend/utils/`: Excel parsing/export, logger, WebDAV, and helper utilities.
- `backend/utils/excel-template-engines/`: business-specific Excel export engines.
- `backend/uploads/`: tracked Excel templates used by export logic.
- `backend/scripts/`: conversion helper assets.

Configuration:

- In development, `backend/configs/config.go` loads `../configs/.env` relative to the backend working directory.
- In production, environment variables are expected from the process/Kubernetes secret.
- Required environment variables used by code:
  - `ENVIRONMENT`
  - `PORT`
  - `MYSQL_USER`
  - `MYSQL_PASSWORD`
  - `MYSQL_SERVER`
  - `UPLOAD_USER`
  - `UPLOAD_PASSWORD`
  - `UPLOAD_SERVER`

Security caveat:

- JWT signing key is hard-coded in `backend/controllers/auth.go`. Treat changes around auth as security-sensitive.
- `LoadConfigFile()` currently prints upload credentials and other config to stdout. Be careful when sharing logs.

Common backend commands:

```sh
cd backend
go test ./...
go build ./...
```

Development server:

```sh
./dev.sh
```

`dev.sh` sets `ENVIRONMENT=development`, changes into `backend`, and runs `air`. This requires `air` to be installed locally.

Production-style local run:

```sh
./run.sh
```

`run.sh` sets `ENVIRONMENT=production`, builds in `backend`, and runs `./super-supply-chain`.

## Frontend

Entry point: `frontend/src/App.tsx`

The frontend is a React Admin application. Resources are registered for:

- `settlement-form-entry`
- `excel-read-rules`
- `dict-manage`
- `yifan/cost-calculation`
- dynamic `excel/:tableName` and `excel/:tableName/:id` routes

Data provider:

- `frontend/src/dataProvider.ts` wraps `ra-data-simple-rest`.
- API base URL comes from `VITE_JSON_SERVER_URL`.
- Auth token is loaded from `localforage` user state and sent as `Authorization: Bearer ...`.
- `401` responses clear the local user and re-run auth checks.

Auth provider:

- `frontend/src/authProvider.ts` uses `VITE_LOGIN_URL` for login.

Common frontend commands:

```sh
cd frontend
yarn
yarn dev
yarn build
yarn type-check
yarn lint
yarn format
```

Notes:

- `frontend/node_modules/` exists locally and is ignored by Git.
- Prefer `yarn` because `frontend/yarn.lock` is present.
- Frontend `.env` is tracked in the repository. Avoid exposing secret values in logs or summaries.

## Deployment

Docker build:

```sh
docker build -t super-supply-chain .
```

The Dockerfile:

1. Builds `frontend` with Node 22.
2. Copies `frontend/dist` into the final Go image.
3. Builds the Go backend binary as `app`.
4. Exposes `8081`.
5. Sets `ENVIRONMENT=production`.

Kubernetes:

- `ssc-deployment.yaml` deploys image `registry.cn-hangzhou.aliyuncs.com/mengfei0053/ssc:1.0.5`.
- Runtime env comes from inline `PORT=80`, `ENVIRONMENT=production`, plus secret `ssc-secret`.
- Container port and service target are `80`.

GitHub Actions:

- `.github/workflows/docker-build-push.yml` runs only on tags matching `v*`.
- It builds and pushes the Docker image to Aliyun registry using `DOCKER_REPOSITORY` secret.

## Coding Guidelines

- Preserve the existing Go package layout and controller/model/util split.
- For backend list endpoints, follow the existing React Admin compatibility pattern: parse query params and set `Content-Range` via helpers in `backend/utils`.
- For frontend pages, follow existing React Admin + MUI conventions instead of introducing another UI framework.
- Do not commit generated binaries or local runtime outputs. Existing tracked binaries under `backend/tmp` and old helper outputs are historical artifacts; avoid touching them unless cleanup is explicitly requested.
- Treat Excel templates and uploaded sample files as business fixtures. Do not delete or replace them without confirming the workflow they support.
- Keep environment-specific secrets out of new source files.

## Verification Checklist

Use the narrowest checks that match the change:

- Backend-only change: `cd backend && go test ./...`
- Backend compile-sensitive change: `cd backend && go build ./...`
- Frontend-only change: `cd frontend && yarn type-check` and usually `yarn build`
- Frontend formatting/linting change: `cd frontend && yarn lint`
- Docker/deployment change: inspect `Dockerfile`, `ssc-deployment.yaml`, and build locally when practical.

If a command cannot run because dependencies or services are missing, report the exact blocker.

## Ongoing Maintenance

Keep this file current. Update it whenever:

- project structure changes,
- build/test/run commands change,
- environment variables change,
- deployment workflow changes,
- a tracked legacy directory is removed or restored,
- new generated artifacts should be ignored or protected.

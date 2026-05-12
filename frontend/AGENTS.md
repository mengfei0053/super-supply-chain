# Frontend Guide

Scope: everything under `frontend/`.

This is a React 19, Vite, React Admin, and MUI application.

## Commands

Run from `frontend/`:

```sh
yarn
yarn dev
yarn build
yarn type-check
yarn lint
yarn format
```

Prefer `yarn` because `yarn.lock` is present.

## Conventions

- Keep UI patterns consistent with React Admin resources and MUI components already in use.
- Use existing `dataProvider` and `authProvider` instead of creating parallel API clients.
- Resource names should match backend routes under `/api/admin`.
- Keep environment variable usage limited to Vite `import.meta.env` access.
- Do not edit `node_modules`, `tmp`, or generated `dist` output.
- Avoid introducing a second component library.

## Environment

Known Vite variables:

- `VITE_JSON_SERVER_URL`
- `VITE_LOGIN_URL`

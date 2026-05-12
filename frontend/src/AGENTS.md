# Frontend Source Guide

Scope: `frontend/src/`.

Main files:

- `App.tsx` registers React Admin resources and custom routes.
- `dataProvider.ts` wraps `ra-data-simple-rest` and handles React Admin list/create behavior.
- `authProvider.ts` handles login persistence and auth checks.
- `index.tsx` mounts the app.

Guidelines:

- Keep API behavior centralized in `dataProvider.ts` unless a component needs a narrow custom request.
- Keep auth storage consistent with `localforage` user data.
- Route and resource additions must line up with backend endpoints.
- Prefer typed props and existing React Admin component patterns.

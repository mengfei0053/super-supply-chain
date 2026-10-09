# Frontend Layout Guide

Scope: `frontend/src/layout/`.

Layout files define the React Admin shell, menu, and submenus.

- Keep menu entries aligned with `App.tsx` resources and backend menu data.
- The header docs link (`AppBar.tsx`) points at same-origin `/docs/`, which is the Docusaurus site, not a React Admin route.
- Preserve navigation paths used by existing pages.
- Avoid page-specific business logic in layout components.
- Test layout changes in the browser when visual or navigation behavior changes.

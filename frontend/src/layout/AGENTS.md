# Frontend Layout Guide

Scope: `frontend/src/layout/`.

Layout files define the React Admin shell, menu, and submenus.

- Keep menu entries aligned with `App.tsx` resources and backend menu data.
- Preserve navigation paths used by existing pages.
- Avoid page-specific business logic in layout components.
- Test layout changes in the browser when visual or navigation behavior changes.

# Frontend Pages Guide

Scope: `frontend/src/pages/`.

Pages are organized by React Admin resource or business workflow.

- Follow the local `Create.tsx`, `Edit.tsx`, `List.tsx`, and `index.ts` pattern where present.
- Keep page-specific API calls close to the page unless they belong in `dataProvider.ts`.
- Preserve React Admin record shape expectations (`id`, list totals, resource names).
- Use MUI layout primitives consistently and keep dense admin workflows scannable.

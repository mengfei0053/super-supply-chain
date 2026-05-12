# Backend Controllers Guide

Scope: `backend/controllers/`.

Controllers are Gin handlers for API and static serving behavior.

- Keep route behavior in sync with `backend/main.go` and frontend resource names.
- For list endpoints used by React Admin, preserve `range`, `sort`, `filter`, and `Content-Range` behavior.
- Use `models.DB` consistently with existing handlers.
- Return JSON shapes expected by the frontend; avoid renaming response fields casually.
- File upload/export handlers should preserve existing Excel template and WebDAV workflows.
- Add or update API tests in `backend/tests` when auth, status codes, query handling, or response shapes change.

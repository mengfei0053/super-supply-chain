# Backend Middleware Guide

Scope: `backend/middleware/`.

Middleware currently handles auth and request logging/recovery.

- Auth changes are security-sensitive. Keep both Authorization header and cookie behavior in mind.
- JWT validation depends on `controllers.JwtKey` and `controllers.Claims`.
- `AuthMiddleware` accepts a session JWT or a personal access token. PAT credentials may arrive as `Authorization: Bearer`, `X-API-Key`, or `?token=`. A PAT in the Authorization header is preferred over a browser session cookie. For JWTs, the cookie still wins over the header.
- `RequireSessionJWT` protects token management routes so a PAT cannot create or revoke tokens.
- Tests for auth behavior belong in `backend/tests`.
- Logging middleware should not emit request bodies or secrets unless explicitly required and scrubbed. Redact `token` query values and `Authorization`, `X-API-Key`, and `Cookie` headers before writing panic dumps.

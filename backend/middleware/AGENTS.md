# Backend Middleware Guide

Scope: `backend/middleware/`.

Middleware currently handles auth and request logging/recovery.

- Auth changes are security-sensitive. Keep both Authorization header and cookie behavior in mind.
- JWT validation depends on `controllers.JwtKey` and `controllers.Claims`.
- Tests for auth behavior belong in `backend/tests`.
- Logging middleware should not emit request bodies or secrets unless explicitly required and scrubbed.

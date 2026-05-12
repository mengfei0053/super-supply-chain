# Backend Uploads Guide

Scope: `backend/uploads/`.

This directory stores tracked Excel templates used by backend export logic.

- Treat these files as business fixtures, not disposable uploads.
- Confirm the related export workflow before replacing, deleting, or renaming templates.
- Avoid adding ad hoc user-uploaded files here; runtime upload outputs should stay outside versioned fixtures.

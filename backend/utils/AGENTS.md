# Backend Utils Guide

Scope: `backend/utils/`.

Utilities cover Excel parsing/export helpers, query parsing, logging, WebDAV upload/download, and small generic helpers.

- Keep generic helpers small and type-safe.
- For Excel parsing/export behavior, prefer structured workbook APIs over ad hoc string manipulation.
- Preserve React Admin query and `Content-Range` compatibility helpers.
- WebDAV helpers depend on config values from `backend/configs`; avoid hard-coded credentials or URLs.
- Add focused tests for pure helpers when logic is nontrivial.

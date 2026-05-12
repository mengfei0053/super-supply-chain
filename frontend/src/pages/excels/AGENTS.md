# Excel Pages Guide

Scope: `frontend/src/pages/excels/`.

These pages manage dynamic Excel data, uploads, export rules, and batch/single exports.

- Keep `tableName` route handling compatible with backend dynamic Excel endpoints.
- Preserve upload request shape and multipart behavior.
- Export actions should continue to use backend endpoints and returned files/paths as currently expected.
- Be careful changing query serialization; it affects backend list and export behavior.
- Verify user-facing Excel workflows when changing upload, template, or export code.

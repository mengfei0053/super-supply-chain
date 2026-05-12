# Backend Models Guide

Scope: `backend/models/`.

Models define GORM database shape and global DB initialization.

- Preserve table/field mappings expected by existing SQL and controllers.
- Be conservative with `AutoMigrate`; it is currently commented out in production initialization.
- Keep database initialization in `init-db.go`.
- Model changes normally need controller review and targeted tests using the in-memory test DB helpers in `backend/tests`.
- Avoid adding business logic to models unless it is tightly tied to data invariants, such as password hashing.

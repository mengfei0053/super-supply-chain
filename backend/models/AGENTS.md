# Backend Models Guide

Scope: `backend/models/`.

Models define GORM database shape and global DB initialization.

- Preserve table/field mappings expected by existing SQL and controllers.
- Be conservative with `AutoMigrate`. Existing business tables stay commented out in `init-db.go`.
- `PersonalAccessToken` is the exception: startup AutoMigrates only that new table. The same schema is in `ssc-sqls/personal_access_tokens.sql` for database users who cannot create tables. A migrate failure is logged and does not stop the process.
- Keep database initialization in `init-db.go`.
- Model changes normally need controller review and targeted tests using the in-memory test DB helpers in `backend/tests`.
- Avoid adding business logic to models unless it is tightly tied to data invariants, such as password hashing.

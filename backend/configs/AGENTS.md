# Backend Config Guide

Scope: `backend/configs/`.

This package owns process configuration.

- Development mode loads `../configs/.env` relative to the backend working directory.
- Production mode expects environment variables from the process or Kubernetes secret.
- Keep config reads centralized here unless a value is genuinely local to another package.
- Do not add default secrets or credentials to source code.
- Be careful with logging. This package currently prints sensitive upload values; do not add more secret output.

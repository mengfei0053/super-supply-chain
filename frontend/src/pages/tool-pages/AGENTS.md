# Tool Pages Guide

Scope: `frontend/src/pages/tool-pages/`.

Tool pages manage configuration-style resources such as dicts and Excel read rules.

- Keep resource modules organized with `Create`, `Edit`, `List`, and `index` files.
- Prefer reusable inputs from `frontend/src/components` when adding repeated controls.
- Changes here often affect backend controller assumptions; verify route names and payload shapes.

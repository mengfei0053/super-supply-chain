# Yifan Pages Guide

Scope: `frontend/src/pages/yifan/`.

This area contains Yifan-specific business workflows.

- Keep business-specific behavior isolated here unless it is truly reusable.
- Coordinate field changes with backend Excel export engines under `backend/utils/excel-template-engines`.
- Prefer explicit naming over generic abstractions because these workflows encode business rules.

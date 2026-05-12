# Excel Template Engines Guide

Scope: `backend/utils/excel-template-engines/`.

This directory contains business-specific Excel export engines.

- Preserve existing template paths and output assumptions unless the business workflow changes.
- Keep calculations explicit and testable; add tests for parsing, totals, and rule-heavy transformations.
- Avoid broad refactors across multiple engines unless shared behavior is clearly duplicated and covered.
- When changing output layout, verify generated workbooks manually or with targeted assertions where practical.

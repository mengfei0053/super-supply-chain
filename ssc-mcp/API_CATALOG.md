# SSC Backend API Catalog

Source of truth: `/opt/super-supply-chain/backend/main.go` on host `101.36.111.17`.

Base origin (production): `https://ssc.mengfei.tech`  
SPA path: `/super-supply-chain/` (static only; not an API prefix).

Auth: session JWT or personal access token via `Authorization: Bearer`, `X-API-Key`, or `?token=` (`middleware.AuthMiddleware`). Public endpoints do not require auth. PAT management is session-JWT only. See `website/docs/personal-access-tokens.md`.

The MCP server exposes only five tools, all 常用: `ssc_upload_excel`, `ssc_list_excel`, `ssc_delete_excel_row`, `ssc_export_excel` (types `shortHaulInvoice`, `invoice_unpacking`, `invoice_clearance_only`, `invoice_freight`), and `ssc_search_companies`. Other routes below are not MCP tools. Auth is `SSC_TOKEN` or `SSC_USERNAME` + `SSC_PASSWORD` (no login tool).

## Public

| Method | Path | Auth | Purpose | Body / params |
| --- | --- | --- | --- | --- |
| POST | `/api/register` | no | Create account | `{username, password}` |
| POST | `/api/login` | no | Login, return JWT + user | `{username, password}` → `{id,username,fullName,email,token,avatar}` |

## Protected (`/api/admin`)

### Settlement / orders

| Method | Path | Purpose | Notes |
| --- | --- | --- | --- |
| GET | `/settlement-form-entry` | List order summaries | `id, orderNumber, arrivalDate, arrivalPort` |
| GET | `/settlement-form-entry/:id` | Order detail (same fields) | |
| POST | `/settlement-form-entries` | Upload settlement file | multipart `file` → NAS |
| PUT | `/settlement-form-entry/:id` | Stub | returns success message only |
| DELETE | `/settlement-form-entry/:id` | Stub | returns success message only |

### Excel read rules

| Method | Path | Purpose | Notes |
| --- | --- | --- | --- |
| GET | `/excel-read-rules` | List rules | requires `range` query (RA style) |
| GET | `/excel-read-rules/:id` | Rule detail | includes mapping `rules` |
| POST | `/excel-read-rules` | Create rule + migrate dynamic table | JSON body |
| PUT | `/excel-read-rules/:id` | Update rule | |
| DELETE | `/excel-read-rules/:id` | Delete rule | |

### Dictionary

| Method | Path | Purpose | Notes |
| --- | --- | --- | --- |
| GET | `/dict-manage` | List dicts | requires `range` |
| GET | `/dict-manage/:id` | Detail | |
| POST | `/dict-manage` | Create | `{key,value,type}` |
| PUT | `/dict-manage/:id` | Update | |
| DELETE | `/dict-manage/:id` | Delete | |
| GET | `/dict-manage/map/:type` | key→value map by type | |

### Dynamic Excel tables

| Method | Path | Purpose | Notes |
| --- | --- | --- | --- |
| GET | `/excel/:tableName` | List rows | `filter` JSON `{start,end}` (or `filter.start` / `filter.end`) on `created_at`; repeated `range=start&range=end` or JSON `range=[start,end]`; optional `sort` object or array (accepted, not applied as ORDER BY). `Content-Range` is the match count. MCP 常用: `ssc_list_excel` |
| GET | `/excel/:tableName/:id` | Row detail | |
| POST | `/excel/:tableName` | Upload Excel → parse → insert | multipart `file` + `name`. MCP 常用: `ssc_upload_excel` |
| PUT | `/excel/:tableName/:id` | Update row | JSON `DynamicExcelTable` |
| DELETE | `/excel/:tableName/:id` | Hard delete | Unscoped. MCP 常用: `ssc_delete_excel_row` |
| GET | `/excel-exports/:tableName` | Export Excel file | repeated `ids` plus `type`. 常用 types: `shortHaulInvoice`, `invoice_unpacking`, `invoice_clearance_only`, `invoice_freight`. MCP 常用: `ssc_export_excel` saves `.xlsx` from `Content-Disposition` (not `.zip` when `Content-Type` is `application/zip`) |

Known dynamic tables (from MySQL):  
`dynamic_Integrity_packaging_invoice`, `dynamic_customs_declaration_form`, `dynamic_settlement_statement_fenchang`, `dynamic_settlement_statement_suqian`, `dynamic_yifan_cost_cal`

### Excel export rules / templates

| Method | Path | Purpose | Notes |
| --- | --- | --- | --- |
| GET | `/excel-export-rule/template/:tableName` | List templates | requires `range` |
| GET | `/excel-export-rule/template/:tableName/:id` | Template detail | |
| POST | `/excel-export-rule/template/:tableName` | Upload template | multipart `file` + `alias` |
| PUT | `/excel-export-rule/:tableName/:id` | Stub | |
| DELETE | `/excel-export-rule/:tableName/:id` | Delete template | |
| POST | `/excel-export-rule/:tableName/export` | Export | |
| GET | `/excel-export-rule/:tableName/export/:id` | Single export | |

### Misc

| Method | Path | Purpose | Notes |
| --- | --- | --- | --- |
| GET | `/options/:key` | Select options | `key=export-templates` + `associated_table` |
| GET | `/menus` | Dynamic menus | `{id, menuName, dynamicTableName}[]` |
| GET | `/companies` | Company keyword search | `keyword` required. `LIKE` on `name` and `alias`. Default `deleted_at IS NULL`; `includeDeleted=true` keeps soft-deleted rows. Returns all columns. MCP 常用: `ssc_search_companies` |

## List query convention

Most list endpoints use React Admin params via `utils.GetListQueryParams`:

- `range` **required** on React Admin lists — JSON `[start, end]` or repeated `range=start&range=end` (end is exclusive; limit = end - start, offset = start). The admin UI sends `range=0&range=50`.
- optional `filter` JSON `{start,end}` or `filter.start` / `filter.end` (dynamic excel `created_at` window; empty dates yield no rows in practice)
- optional `sort` JSON object `{"field":"id","order":"ASC"}` or array `["id","ASC"]`. `GetDynamicExcelTableList` stores the raw value and does not `ORDER BY` it.

---
name: SSC MCP
description: When connecting to or calling Super Supply Chain via remote MCP (upload/list/delete/export Excel, search companies).
---

# SSC MCP

Use this when connecting Cursor or another agent to Super Supply Chain over the public remote MCP, or when calling its Excel and company tools.

## Connect

Use public HTTPS only.

- Preferred (Streamable HTTP): `https://ssc.mengfei.tech/mcp`
- Older clients that only speak HTTP+SSE (protocol 2024-11-05): `https://ssc.mengfei.tech/sse`

A client that takes one URL (Cursor) sets `url` to `https://ssc.mengfei.tech/sse` for that older transport. Clients that post SSE messages themselves use `https://ssc.mengfei.tech/messages?sessionId=...` after the session id is known.

The admin UI is `https://ssc.mengfei.tech/super-supply-chain/`. That path is the web app. JSON APIs sit beside it, under `/api` and `/api/admin`.

## Cursor `mcpServers`

Create a personal access token first. Put it in `headers` together with `url`. `Authorization` includes the `Bearer ` prefix and a space. Replace `ssc_pat_…` with the token. Never commit a real token, and do not paste one into this skill, the repository, or shared logs.

```json
{
  "mcpServers": {
    "ssc": {
      "url": "https://ssc.mengfei.tech/mcp",
      "headers": {
        "Authorization": "Bearer ssc_pat_…"
      }
    }
  }
}
```

Optional: send the same token as `X-API-Key` with no `Bearer ` prefix.

```json
{
  "mcpServers": {
    "ssc": {
      "url": "https://ssc.mengfei.tech/mcp",
      "headers": {
        "X-API-Key": "ssc_pat_…"
      }
    }
  }
}
```

For SSE, keep the same `headers` and set `url` to `https://ssc.mengfei.tech/sse`.

Save this in `~/.cursor/mcp.json`, or in the project file `.cursor/mcp.json`.

## Personal access token

1. Sign in to the admin UI at `https://ssc.mengfei.tech/super-supply-chain/`.
2. Open **工具 → 个人访问令牌**.
3. Create a token. A name such as `cursor-mcp` is enough. Leave the scope empty so the token can call the admin APIs these tools use. Set an expiry only when you want one.
4. Copy the plaintext once. It starts with `ssc_pat_`. The UI does not show it again.
5. Store it only in the local MCP `headers`. Revoke it on the same screen when it is no longer needed.

## Tools

Five tools. Example table name: `dynamic_settlement_statement_suqian`.

### `ssc_upload_excel`

Use when inserting one parsed workbook into a dynamic table (`POST /api/admin/excel/{tableName}`).

- `tableName`
- `filePath` — `.xlsx` or `.xls` path the MCP server process can read. On `https://ssc.mengfei.tech/mcp`, that path is on the server.
- `name` — multipart name field. When it ends in `.xls` or `.xlsx`, it is also the stored file name.

### `ssc_list_excel`

Use when listing rows (`GET /api/admin/excel/{tableName}`).

- `tableName`
- `filterStart` and `filterEnd` — `YYYY-MM-DD`, inclusive, on `created_at`. Always send both. The 宿迁结算 table returns an empty list when this date window is missing or does not cover the rows.
- `sort` — optional. `{"field":"id","order":"ASC"}` or `["id","ASC"]`. The API accepts either form and does not apply it as `ORDER BY`.
- `range` — optional `[start, end)`. Default `[0, 50]`.

### `ssc_delete_excel_row`

Use when permanently removing one row (`DELETE /api/admin/excel/{tableName}/{id}`).

- `tableName`
- `id` — one row id

### `ssc_export_excel`

Use when downloading an invoice workbook (`GET /api/admin/excel-exports/{tableName}`). The tool saves one `.xlsx` and returns that path. Keep the `.xlsx` extension when `Content-Type` is `application/zip` (an xlsx file is a zip container).

- `tableName`
- `ids` — one id, a comma-separated string such as `"896,897"`, or an array such as `["896","897"]`
- `type` — one of:
  - `shortHaulInvoice` (短驳发票)
  - `invoice_unpacking` (拆箱发票)
  - `invoice_clearance_only` (清关发票)
  - `invoice_freight` (运费发票)
- `outputPath` — optional file or directory on the MCP host. A path that ends in `.zip` is still written as `.xlsx` for a single workbook.

After connecting, confirm the credential with one read-only export of an existing row (`type` `shortHaulInvoice`). Success returns a path ending in `.xlsx`.

### `ssc_search_companies`

Use when looking up companies by keyword (`GET /api/admin/companies?keyword=`). Matches `name` or `alias` as a literal substring (`%` and `_` are not wildcards). Returns every column.

- `keyword`
- `includeDeleted` — optional. Default `false` omits soft-deleted rows.

## Auth

Each tool calls the SSC backend under `/api/admin`. The personal access token in the MCP `headers` is the same token an HTTP client sends as `Authorization: Bearer ssc_pat_…` or `X-API-Key: ssc_pat_…`.

A login JWT uses that same header and expires after 24 hours. A personal access token stays valid until it expires or is revoked. There is no MCP login tool. Settlement-file upload (`POST /api/admin/settlement-form-entries`) and account registration are not wrapped.

# ssc-mcp

MCP server for the commonly used **Super Supply Chain (SSC)** Excel list, upload, delete, and export APIs, plus company keyword search.

- Project path: `/workspace/ssc-mcp`
- Live UI: <https://ssc.mengfei.tech/super-supply-chain/>
- Live API origin: `https://ssc.mengfei.tech` (routes under `/api` and `/api/admin`)

> **Important:** Do **not** set `SSC_BASE_URL` to `.../super-supply-chain`. That path is the SPA only. JSON APIs are at `/api/*` on the same host.

## Install

```bash
cd /workspace/ssc-mcp
npm install
npm run build
```

## Environment

| Variable | Required | Description |
| --- | --- | --- |
| `SSC_BASE_URL` | no | API origin, default `https://ssc.mengfei.tech` |
| `SSC_TOKEN` | one of token / user+pass | Login JWT, or a personal access token from the admin UI (preferred for long-lived MCP) |
| `SSC_USERNAME` | with password | Account for `POST /api/login` |
| `SSC_PASSWORD` | with username | Password for login |
| `SSC_MCP_TRANSPORT` | no | `stdio` (default) or `http` |
| `SSC_MCP_HOST` | no | HTTP bind address, default `0.0.0.0` |
| `SSC_MCP_PORT` | no | HTTP port, default `3100` |
| `SSC_MCP_CORS_ORIGIN` | no | Browser CORS allow-list, default `*` |
| `SSC_MCP_ALLOWED_HOSTS` | no | Optional comma-separated `Host` allow-list for DNS-rebinding protection |
| `SSC_EXPORT_DIR` | no | Directory for `ssc_export_excel` downloads |

Copy `.env.example` → `.env` (mode `600`) and fill values. Never commit tokens or passwords.

CLI flags override the environment: `--transport`, `--host`, `--port`, `--cors-origin`, `--allowed-hosts`.

## Run (stdio)

```bash
export SSC_BASE_URL=https://ssc.mengfei.tech
export SSC_TOKEN='<jwt from browser localStorage / ssc login CLI>'
npm run start
# or: npx tsx src/index.ts
```

## Cursor / MCP config snippet

Add to your MCP settings (e.g. Cursor `mcp.json`):

```json
{
  "mcpServers": {
    "ssc": {
      "command": "node",
      "args": ["/workspace/ssc-mcp/dist/index.js"],
      "env": {
        "SSC_BASE_URL": "https://ssc.mengfei.tech",
        "SSC_TOKEN": "<paste JWT here>"
      }
    }
  }
}
```

Or with username/password (the server calls `POST /api/login` on the first tool call):

```json
{
  "mcpServers": {
    "ssc": {
      "command": "npx",
      "args": ["tsx", "/workspace/ssc-mcp/src/index.ts"],
      "env": {
        "SSC_BASE_URL": "https://ssc.mengfei.tech",
        "SSC_USERNAME": "testuser",
        "SSC_PASSWORD": "<password>"
      }
    }
  }
}
```

## Run (HTTP)

HTTP mode is opt-in. It keeps the same tools and the same `SscClient` auth behavior. Streamable HTTP (the SDK's current remote transport) is served at `/mcp`. Legacy HTTP+SSE (protocol 2024-11-05) is also served for older clients.

**Production** publishes MCP behind Nginx Proxy Manager on the same host as the app:

| Transport | Production URL | Local debug | Methods |
| --- | --- | --- | --- |
| Streamable HTTP (use this) | `https://ssc.mengfei.tech/mcp` | `http://127.0.0.1:3100/mcp` | `POST`, `GET`, `DELETE` |
| Legacy SSE | `https://ssc.mengfei.tech/sse` | `http://127.0.0.1:3100/sse` | `GET` |
| Legacy SSE messages | `https://ssc.mengfei.tech/messages?sessionId=...` | `http://127.0.0.1:3100/messages?sessionId=...` | `POST` |

Compose service `ssc-mcp` binds `172.17.0.1:3100` only (not the public NIC). Credentials live in `.env.compose` (`SSC_TOKEN`, or `SSC_USERNAME` + `SSC_PASSWORD`).

Local HTTP process:

```bash
export SSC_BASE_URL=https://ssc.mengfei.tech
export SSC_TOKEN='<jwt>'
export SSC_MCP_TRANSPORT=http
export SSC_MCP_HOST=0.0.0.0
export SSC_MCP_PORT=3100
npm run start
# or: npm run dev -- --transport http --port 3100
```

The process listens on `SSC_MCP_HOST` (default `0.0.0.0`). Point clients at the production HTTPS URL, or at `127.0.0.1` for local debug — never at `0.0.0.0`.

One OS process shares one SSC credential (`SSC_TOKEN`, or `SSC_USERNAME` + `SSC_PASSWORD`) across every MCP session. Run a separate process per tenant. Public access should go through the domain reverse proxy, not raw `:3100`.

Sessions are kept in memory. A session is removed when the client sends `DELETE /mcp`, the legacy SSE connection closes, or the session has had no open request for 30 minutes. At most 200 sessions are accepted.

### CORS (browser clients)

Desktop Cursor talks to the URL directly and does not need CORS. Browser-based MCP clients do. The server answers `OPTIONS` and sets:

- `Access-Control-Allow-Origin`: `*` by default, or the request `Origin` when it is in `SSC_MCP_CORS_ORIGIN` (comma-separated)
- `Access-Control-Allow-Methods`: `GET, POST, DELETE, OPTIONS`
- `Access-Control-Allow-Headers`: `Content-Type`, `Accept`, `Authorization`, `Mcp-Session-Id`, `Mcp-Protocol-Version`, `Last-Event-ID`
- `Access-Control-Expose-Headers`: `Mcp-Session-Id`, `Mcp-Protocol-Version`

`*` cannot be combined with credentialed browser requests. For a specific site, set `SSC_MCP_CORS_ORIGIN=https://your-app.example`.

Binding to `0.0.0.0` or `::` without `SSC_MCP_ALLOWED_HOSTS` disables the SDK's DNS-rebinding check and prints a warning. Localhost binds (`127.0.0.1`, `localhost`, `::1`) enable that check automatically. To listen on all interfaces and still restrict `Host`, set `SSC_MCP_ALLOWED_HOSTS=localhost,127.0.0.1`.

### Cursor HTTP config

Production (recommended) — MCP is already behind `ssc.mengfei.tech`:

```json
{
  "mcpServers": {
    "ssc": {
      "url": "https://ssc.mengfei.tech/mcp"
    }
  }
}
```

Older clients that only speak SSE:

```json
{
  "mcpServers": {
    "ssc": {
      "url": "https://ssc.mengfei.tech/sse"
    }
  }
}
```

For local debug only, use `http://127.0.0.1:3100/mcp` (or `/sse`) after starting the process yourself.

`SSC_TOKEN` / `SSC_USERNAME` / `SSC_PASSWORD` belong in the server process environment (production: compose `.env.compose`), not in the URL snippet.

## Tools (5)

These are the only tools. There is no login or status tool: set `SSC_TOKEN`, or `SSC_USERNAME` plus `SSC_PASSWORD`. With username and password, the process calls `POST /api/login` on the first tool call and keeps the JWT in memory. The token is never printed.

### 常用

Upload, list, delete, these four exports, and company keyword search are the commonly used APIs:

| 常用 | API | Tool |
| --- | --- | --- |
| 常用 | `POST /api/admin/excel/{tableName}` multipart `file` + `name` | `ssc_upload_excel` |
| 常用 | `GET /api/admin/excel/{tableName}?filter={start,end}&range=0&range=50` | `ssc_list_excel` |
| 常用 | `DELETE /api/admin/excel/{tableName}/{id}` | `ssc_delete_excel_row` |
| 常用 | `GET /api/admin/excel-exports/{tableName}?ids={ids}&type=shortHaulInvoice` | `ssc_export_excel` |
| 常用 | `GET /api/admin/excel-exports/{tableName}?ids={ids}&type=invoice_unpacking` | `ssc_export_excel` |
| 常用 | `GET /api/admin/excel-exports/{tableName}?ids={ids}&type=invoice_clearance_only` | `ssc_export_excel` |
| 常用 | `GET /api/admin/excel-exports/{tableName}?ids={ids}&type=invoice_freight` | `ssc_export_excel` |
| 常用 | `GET /api/admin/companies?keyword={keyword}` | `ssc_search_companies` |

`ssc_upload_excel` (常用) reads a local `.xlsx`/`.xls` path on the machine running the MCP server and posts it as multipart `file`, plus form field `name`. Example table: `dynamic_settlement_statement_suqian`. This inserts a row.

`ssc_list_excel` (常用) lists rows from a dynamic table. `filterStart` / `filterEnd` are `YYYY-MM-DD` and are sent as `filter={"start","end"}` (`created_at` window; required for non-empty 宿迁结算 results). Optional `sort` is `{"field":"id","order":"ASC"}` or `["id","ASC"]`. `range` defaults to `[0, 50]` and is sent as repeated `range=0&range=50` (`range=[0,50]` also works on the API). The result includes the JSON rows and the `Content-Range` total when the API sends it. The list handler accepts `sort` but does not `ORDER BY` it.

`ssc_delete_excel_row` (常用) hard-deletes one row. Example: `DELETE /api/admin/excel/dynamic_settlement_statement_suqian/896`.

`ssc_export_excel` (常用) downloads a workbook and writes it under the OS temp directory (`ssc-mcp-exports`, or `SSC_EXPORT_DIR` / `outputPath`). The tool result is metadata plus `path` — read that file; the bytes are not inlined. The saved name comes from `Content-Disposition` (`filename*` / `filename`) and ends in `.xlsx`. These responses are often `Content-Type: application/zip` because Office Open XML is a zip; that header is not used as the file extension. If `outputPath` ends in `.zip` and the body is one workbook, the tool still writes `.xlsx`. A `.zip` path is kept only when the body is a zip of several separate workbooks. `ids` may be `896`, `896,897`, or an array. Commonly used `type` values are `shortHaulInvoice`, `invoice_unpacking`, `invoice_clearance_only`, and `invoice_freight`.

`ssc_search_companies` (常用) searches `base_companies_infos` by a literal substring of `name` or `alias` (`LIKE %keyword%`, with `%` and `_` escaped). The `rows` array includes every column: `id`, `created_at`, `updated_at`, `deleted_at`, `name`, `addr_country`, `addr_province`, `addr_city`, `addr_street`, `unified_social_credit_code`, `bank_code`, `phone_num`, `alias`, `target_addr`. Rows with `deleted_at` set are omitted unless `includeDeleted` is true.

Orders, dictionaries, menus, read rules, and other export types are not exposed.

## Smoke test

Local tests do not need SSC credentials. They mock the upload/export HTTP calls and exercise both MCP transports:

```bash
npm test
npm run build
```

Transport smoke (no SSC credentials). Starts HTTP on an ephemeral port, handshakes Streamable HTTP and legacy SSE, lists the five tools, calls delete without credentials, checks CORS, then handshakes the default stdio server:

```bash
npm run smoke:http
```

Live check against `https://ssc.mengfei.tech` (does not print the JWT):

```bash
export SSC_BASE_URL=https://ssc.mengfei.tech
export SSC_TOKEN='...'   # or SSC_USERNAME + SSC_PASSWORD
npm run smoke
```

Optional live export (read-only). Prints the saved path and byte size only:

```bash
export SSC_SMOKE_EXPORT_TABLE=dynamic_settlement_statement_suqian
export SSC_SMOKE_EXPORT_IDS=896
export SSC_SMOKE_EXPORT_TYPE=shortHaulInvoice
npm run smoke
```

Optional live upload **inserts a row**. Set `SSC_SMOKE_UPLOAD_FILE` to a local `.xlsx` and `SSC_SMOKE_UPLOAD_NAME` (defaults to the file name). Skip this unless you intend to write to that table.

### Manual verification

Use this when `SSC_TOKEN` / `SSC_USERNAME`+`SSC_PASSWORD` are not available in the environment. Do not paste the JWT into logs or chat.

1. Log in at <https://ssc.mengfei.tech/super-supply-chain/> and copy the JWT from the app's stored user (or mint one with the `ssc login` CLI). Export it as `SSC_TOKEN` in a shell that is not recorded.
2. Export an existing row (no database write):

```bash
cd /workspace/ssc-mcp
export SSC_BASE_URL=https://ssc.mengfei.tech
# export SSC_TOKEN=...   # do not echo it
npx tsx -e '
import { SscClient, loadConfigFromEnv } from "./src/client.ts";
const file = await new SscClient(loadConfigFromEnv()).exportExcel({
  tableName: "dynamic_settlement_statement_suqian",
  ids: "896",
  type: "shortHaulInvoice",
});
console.log(JSON.stringify({ path: file.path, bytes: file.bytes, fileName: file.fileName, xlsx: file.xlsx }));
'
```

3. Confirm the printed `path` is a non-empty `.xlsx` (`xlsx: true`). A `Content-Type` of `application/zip` does not change that extension.
4. Upload only against a table you mean to change:

```bash
npx tsx -e '
import { SscClient, loadConfigFromEnv } from "./src/client.ts";
const uploaded = await new SscClient(loadConfigFromEnv()).uploadExcel({
  tableName: "dynamic_settlement_statement_suqian",
  filePath: "/absolute/path/to/file.xlsx",
  name: "file.xlsx",
});
console.log(JSON.stringify({ fileName: uploaded.fileName, bytes: uploaded.bytes, ok: true }));
'
```

## Auth notes

- No MCP login tool. Put a JWT or personal access token in `SSC_TOKEN`, or set `SSC_USERNAME` and `SSC_PASSWORD` so the server can call `POST /api/login`.
- JWT is HS256, 24h expiry. A personal access token stays valid until it expires or is revoked. Both use `Authorization: Bearer <token>`. Do not print it. Create a PAT in the admin UI; see `website/docs/personal-access-tokens.md`.
- Existing DB account: `testuser` (password not stored in this repo).
- Project CLI: `/opt/super-supply-chain/cli` (`ssc login`) can mint a local token to place in `SSC_TOKEN`.

## API catalog

See [API_CATALOG.md](./API_CATALOG.md).

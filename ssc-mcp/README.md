# ssc-mcp

MCP server that wraps the **Super Supply Chain (SSC)** Go/Gin backend HTTP APIs for Cursor / Grok Bot.

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
| `SSC_TOKEN` | one of token / user+pass | JWT from login (preferred) |
| `SSC_USERNAME` | with password | Account for `POST /api/login` |
| `SSC_PASSWORD` | with username | Password for login |
| `SSC_MCP_TRANSPORT` | no | `stdio` (default) or `http` |
| `SSC_MCP_HOST` | no | HTTP bind address, default `0.0.0.0` |
| `SSC_MCP_PORT` | no | HTTP port, default `3100` |
| `SSC_MCP_CORS_ORIGIN` | no | Browser CORS allow-list, default `*` |
| `SSC_MCP_ALLOWED_HOSTS` | no | Optional comma-separated `Host` allow-list for DNS-rebinding protection |

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

Or with username/password (server calls login on first authenticated tool):

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

```bash
export SSC_BASE_URL=https://ssc.mengfei.tech
export SSC_TOKEN='<jwt>'
export SSC_MCP_TRANSPORT=http
export SSC_MCP_HOST=0.0.0.0
export SSC_MCP_PORT=3100
npm run start
# or: npm run dev -- --transport http --port 3100
```

Endpoints after startup:

| Transport | URL | Methods |
| --- | --- | --- |
| Streamable HTTP (use this) | `http://127.0.0.1:3100/mcp` | `POST`, `GET`, `DELETE` |
| Legacy SSE | `http://127.0.0.1:3100/sse` | `GET` |
| Legacy SSE messages | `http://127.0.0.1:3100/messages?sessionId=...` | `POST` |

The process listens on `SSC_MCP_HOST` (default `0.0.0.0`). Point clients at `127.0.0.1` or another reachable address, not at `0.0.0.0`.

One OS process shares one SSC login (`SSC_TOKEN` or `ssc_login`) across every MCP session. Run a separate process per tenant, and do not expose the port on an untrusted network.

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

Start the server yourself, then point Cursor at the Streamable HTTP endpoint:

```json
{
  "mcpServers": {
    "ssc": {
      "url": "http://127.0.0.1:3100/mcp"
    }
  }
}
```

Older clients that only speak SSE can use the legacy endpoint instead:

```json
{
  "mcpServers": {
    "ssc": {
      "url": "http://127.0.0.1:3100/sse"
    }
  }
}
```

`SSC_TOKEN` / `SSC_USERNAME` / `SSC_PASSWORD` belong in the server process environment, not in the URL snippet.

## Tools (20)

| Tool | Backend |
| --- | --- |
| `ssc_login` | `POST /api/login` |
| `ssc_status` | `GET /api/admin/menus` (auth check) |
| `ssc_list_menus` | `GET /api/admin/menus` |
| `ssc_list_orders` | `GET /api/admin/settlement-form-entry` |
| `ssc_get_order` | `GET /api/admin/settlement-form-entry/:id` |
| `ssc_list_dicts` | `GET /api/admin/dict-manage` |
| `ssc_get_dict` | `GET /api/admin/dict-manage/:id` |
| `ssc_get_dict_map` | `GET /api/admin/dict-manage/map/:type` |
| `ssc_create_dict` | `POST /api/admin/dict-manage` |
| `ssc_update_dict` | `PUT /api/admin/dict-manage/:id` |
| `ssc_delete_dict` | `DELETE /api/admin/dict-manage/:id` |
| `ssc_list_excel_read_rules` | `GET /api/admin/excel-read-rules` |
| `ssc_get_excel_read_rule` | `GET /api/admin/excel-read-rules/:id` |
| `ssc_list_excel_rows` | `GET /api/admin/excel/:tableName` |
| `ssc_get_excel_row` | `GET /api/admin/excel/:tableName/:id` |
| `ssc_update_excel_row` | `PUT /api/admin/excel/:tableName/:id` |
| `ssc_delete_excel_row` | `DELETE /api/admin/excel/:tableName/:id` |
| `ssc_list_export_rules` | `GET /api/admin/excel-export-rule/template/:tableName` |
| `ssc_list_export_template_options` | `GET /api/admin/options/export-templates` |
| `ssc_get_export_rule` | `GET /api/admin/excel-export-rule/template/:tableName/:id` |

Not wrapped (multipart / binary / stubs): Excel file upload create, settlement file upload, bulk Excel export download, register user.

## Smoke test

API smoke (live SSC, needs a token or username/password; does not print secrets):

```bash
export SSC_BASE_URL=https://ssc.mengfei.tech
export SSC_TOKEN='...'   # or USERNAME+PASSWORD
npm run smoke
```

Transport smoke (no SSC credentials). Starts HTTP on an ephemeral port, handshakes Streamable HTTP and legacy SSE, lists tools, calls `ssc_status`, checks CORS, then handshakes the default stdio server:

```bash
npm run smoke:http
```

## Auth notes

- JWT is HS256, 24h expiry, header `Authorization: Bearer <token>`.
- Existing DB account: `testuser` (password not stored in this repo).
- Project CLI: `/opt/super-supply-chain/cli` (`ssc login` / `ssc status`) can mint a local token file.

## API catalog

See [API_CATALOG.md](./API_CATALOG.md).

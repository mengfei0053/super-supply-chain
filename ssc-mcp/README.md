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

Copy `.env.example` → `.env` (mode `600`) and fill values. Never commit tokens or passwords.

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

## Tools (19)

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

```bash
export SSC_BASE_URL=https://ssc.mengfei.tech
export SSC_TOKEN='...'   # or USERNAME+PASSWORD
npm run smoke
```

## Auth notes

- JWT is HS256, 24h expiry, header `Authorization: Bearer <token>`.
- Existing DB account: `testuser` (password not stored in this repo).
- Project CLI: `/opt/super-supply-chain/cli` (`ssc login` / `ssc status`) can mint a local token file.

## API catalog

See [API_CATALOG.md](./API_CATALOG.md).

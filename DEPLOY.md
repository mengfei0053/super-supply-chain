# Super Supply Chain on 101

- Path: `/opt/super-supply-chain`
- Secrets: `.env.compose` (mode 600); original Mac env also at `configs/.env` (mode 600)
- URL: https://ssc.mengfei.tech/super-supply-chain/ (docs: https://ssc.mengfei.tech/docs/)
- MCP: https://ssc.mengfei.tech/mcp (legacy SSE: https://ssc.mengfei.tech/sse , POST messages: https://ssc.mengfei.tech/messages)
- App port: **172.17.0.1:8088** -> container 8081 (not on the public NIC; NPM uses docker0). 8317 is CLI Proxy — do not reuse
- MCP port: **172.17.0.1:3100** -> container 3100 (docker0 only; not a public raw port)
- Stack: `ssc-mysql` + `ssc-app` + `ssc-mcp`, network `ssc-net`, restart unless-stopped
- DB name: `super_supply_chain` (imported dump, 18 tables)
- App image: root `Dockerfile` (`registry.cn-hangzhou.aliyuncs.com/mengfei0053/node:22` and `golang:1.23.6-alpine`)
- MCP image: `ssc-mcp/Dockerfile` (same Node 22 base)

## Update

```sh
# from a machine that can reach this VPS
rsync -az --exclude frontend/node_modules --exclude .git ./ ubuntu@101.36.111.17:/opt/super-supply-chain/
ssh ubuntu@101.36.111.17 'cd /opt/super-supply-chain && sudo docker compose --env-file .env.compose build ssc-app ssc-mcp && sudo docker compose --env-file .env.compose up -d'
```

## MCP (`ssc-mcp`)

The compose service runs the Node server in HTTP mode:

- `SSC_MCP_TRANSPORT=http`
- `SSC_MCP_HOST=0.0.0.0`
- `SSC_MCP_PORT=3100`
- `SSC_BASE_URL=https://ssc.mengfei.tech`
- `SSC_MCP_ALLOWED_HOSTS=ssc.mengfei.tech,172.17.0.1,127.0.0.1,localhost`

The process still calls the SSC API with a server credential from `.env.compose` (mode 600). Set one of these. Do not commit the file.

- `SSC_TOKEN` — JWT from login (preferred)
- `SSC_USERNAME` and `SSC_PASSWORD` — the process calls `POST /api/login`

That credential is the existing SSC login, not a new MCP-only secret. Cursor's public HTTP config is only the URL until client auth ships.

Personal access tokens will be created, listed, and revoked in the admin UI (backend and frontend). After that API exists, `https://ssc.mengfei.tech/mcp` will accept those PATs as the client credential. Do not add a second token type, query-string secret, or MCP-only key in the meantime. This deploy does not invent the PAT request shape.

Until that client check exists, anyone who can open the public MCP URLs uses the account configured on the container.

The service keeps the three tools only: upload, delete, and export.

### Nginx Proxy Manager

On the `ssc.mengfei.tech` proxy host, add three custom locations. Scheme `http`, forward hostname `172.17.0.1`, forward port `3100`. Leave the request path intact so each location hits the same path on the MCP container.

| Custom location | Upstream |
| --- | --- |
| `/mcp` | `http://172.17.0.1:3100` |
| `/sse` | `http://172.17.0.1:3100` |
| `/messages` | `http://172.17.0.1:3100` |

`/mcp` is Streamable HTTP. `/sse` is legacy SSE (`GET`). `/messages` is the legacy SSE `POST` (`/messages?sessionId=...`). These paths must not be forwarded to the app on `172.17.0.1:8088`. The Go process does not serve MCP.

Streamable HTTP and legacy SSE keep a response open. In each custom location's advanced config:

```nginx
proxy_buffering off;
proxy_cache off;
proxy_read_timeout 3600s;
proxy_send_timeout 3600s;
```

NPM's default `Host` header is the public name (`ssc.mengfei.tech`). That name is in `SSC_MCP_ALLOWED_HOSTS`. Do not point clients at `http://<public-ip>:3100`.

## Notes

- Public `http://101.36.111.17:8088` and `:3100` are intentionally closed. Use the domain via Nginx Proxy Manager (`172.17.0.1:8088` for the app, `172.17.0.1:3100` for MCP).
- WebDAV/NAS (`UPLOAD_SERVER` points at home LAN 192.168.50.63) is unreachable from 101; uploads will fail until VPN/tunnel or new URL.
- Production logs print UPLOAD_* via `configs.LoadConfigFile()` — avoid sharing raw `docker logs`.
- Root `/` redirects on GET (301) to `/super-supply-chain`; HEAD may 404 (Gin route is GET-only).

## Uploads (production on 101)

`UPLOAD_SERVER=file:///data/ssc-uploads` stores Excel uploads on Docker volume `ssc-uploads`.
Home LAN WebDAV (`192.168.50.63:5000`) is unreachable from this VPS.
Set `UPLOAD_SERVER` to a reachable WebDAV URL to restore NAS mode; `file://` / `local://` keep local disk mode.
